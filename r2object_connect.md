# Panduan Implementasi Upload Gambar ke Cloudflare R2

### Stack: Go + Gin (tanpa GORM) & Flutter Clean Architecture

---

## 1. Latar Belakang & Use Case

Fitur ini digunakan untuk menyimpan berbagai jenis gambar dalam aplikasi:

- Foto profil pengguna
- Foto produk / resep
- Foto harian anak (bergaya feed, mirip Instagram — kemungkinan volume upload tinggi dan berulang)

Gambar disimpan di **Cloudflare R2** (object storage S3-compatible), sedangkan **URL/metadata gambar disimpan di database**, bukan file binary-nya.

---

## 2. Pola Arsitektur yang Direkomendasikan: Presigned URL

**Prinsip utama:** file tidak diupload lewat backend Go (kecuali kasus tertentu yang dijelaskan di bagian 2.1). Backend hanya berperan sebagai "penjaga gerbang" yang menerbitkan izin upload dan mencatat metadata — file transfer terjadi langsung antara Flutter dan R2.

### Kenapa presigned URL, bukan proxy upload?

- Bandwidth file besar tidak membebani compute server Go.
- R2 tidak mengenakan biaya egress, jadi cocok dipakai sebagai storage utama untuk media.
- Lebih scalable untuk fitur dengan volume upload tinggi/berulang (seperti foto harian anak).

### 2.1 Kapan proxy upload (satu endpoint saja) masih masuk akal?

Untuk kasus sederhana dan jarang seperti foto profil, upload langsung lewat satu endpoint backend (multipart form → backend forward ke R2) tetap valid dan lebih simpel. Presigned URL lebih relevan untuk fitur dengan upload sering/banyak sekaligus.

---

## 3. Flow Lengkap (Pola Presigned URL — 3 Tahap)

| #   | Endpoint / Aksi                                              | Penjelasan                                                                                                                                                                                                                      |
| --- | ------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | `POST /api/images/presign-upload` (backend Go)               | Flutter minta izin upload. Backend generate `object_key` unik dan presigned PUT URL dari R2. Balikan: `{upload_url, object_key}`                                                                                                |
| 2   | `PUT <upload_url>` (langsung ke domain R2, bukan ke backend) | Flutter upload file langsung ke Cloudflare R2 menggunakan presigned URL. Header `Content-Type` **harus sama persis** dengan yang dipakai saat generate presign URL, atau signature akan mismatch (403).                         |
| 3   | `POST /api/images/confirm` (backend Go)                      | Setelah upload sukses, Flutter kirim `object_key` + metadata (owner_id, category) ke backend. Backend baru menyimpan record ke database. Opsional: backend bisa `HEAD` object ke R2 dulu untuk verifikasi file benar-benar ada. |

### Kenapa 3 tahap ini dipisah, bukan digabung satu endpoint?

**Tahap 1 vs 2 — wajib terpisah (bukan pilihan desain):**
Tahap 2 adalah request ke server Cloudflare R2 (`https://<account>.r2.cloudflarestorage.com/...`), bukan ke server backend Go. Ini dua server fisik berbeda — tidak mungkin digabung.

**Tahap 1 vs 3 — pilihan desain, dengan alasan berikut:**

1. **Generate presign itu instan, upload itu lama & tidak pasti.** Kalau digabung, backend harus menahan koneksi terbuka selama proses upload berlangsung — menghilangkan keuntungan utama presigned URL.
2. **Retry lebih murah.** Kalau koneksi Flutter putus di tengah upload (umum di mobile), Flutter cukup retry `PUT` ke URL yang sama tanpa perlu presign ulang, selama URL belum expired.
3. **Mencegah orphan record.** Record baru disimpan ke DB _setelah_ file terbukti ada di R2 (tahap 3), sehingga menghindari data kotor berupa baris DB yang menunjuk ke file yang gagal/batal diupload.
4. **Progress tracking akurat di UI.** Karena Flutter upload langsung ke R2, progress bar yang ditampilkan adalah progress transfer yang sesungguhnya — bukan progress upload ke backend saja seperti pada pola proxy.

---

## 4. Setup Client R2 di Go (S3-Compatible SDK)

R2 kompatibel dengan S3 API, sehingga digunakan `aws-sdk-go-v2` dengan endpoint custom mengarah ke R2 (bukan AWS).

```go
package r2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3        *s3.Client
	presigner *s3.PresignClient
	bucket    string
}

func New(ctx context.Context, accountID, accessKeyID, secretAccessKey, bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""),
		),
		config.WithRegion("auto"), // wajib diisi SDK, tapi tidak dipakai R2
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &Client{
		s3:        s3Client,
		presigner: s3.NewPresignClient(s3Client),
		bucket:    bucket,
	}, nil
}

func (c *Client) PresignPutURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error) {
	req, err := c.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (c *Client) PresignGetURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	req, err := c.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err
}
```

---

## 5. Penempatan di Clean Architecture (Go + Gin, tanpa GORM)

Karena tidak memakai GORM, buat abstraksi storage sebagai interface terpisah dari repository database (yang query-nya manual, misalnya pakai `sqlx` atau `pgx`).

```
internal/
  domain/
    entity/image.go
    repository/image_repository.go     // interface, implementasi pakai sqlx/pgx
    storage/object_storage.go          // interface, decouple dari R2 spesifik
  usecase/
    image_usecase.go
  infrastructure/
    persistence/postgres/image_repository.go
    storage/r2/client.go               // implementasi ObjectStorage
  delivery/
    http/handler/image_handler.go
```

```go
// domain/storage/object_storage.go
type ObjectStorage interface {
	PresignPutURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
	PresignGetURL(ctx context.Context, key string, expires time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}
```

Dengan interface ini, layer `usecase` tidak bergantung langsung pada R2 — kalau provider storage berganti di kemudian hari, cukup ganti implementasinya saja.

---

## 6. Konvensi Penamaan Object Key

Gunakan prefix per kategori + owner agar mudah dikelola/dihapus per user:

```
profile/{user_id}/{uuid}.webp
product/{product_id}/{uuid}.webp
recipe/{recipe_id}/{uuid}.webp
child-daily/{child_id}/{yyyy}/{mm}/{uuid}.webp
```

---

## 7. Skema Database

**Prinsip:** kolom bernama `photo_url` / `image_url` **tetap dipertahankan** (tidak diganti jadi `*_key`) agar konsisten dengan konvensi yang sudah dipakai di skema saat ini. Yang berubah hanya _apa yang disimpan di dalamnya_ — bukan URL lengkap, tapi **`object_key`** R2. URL final dihitung saat response, bukan disimpan permanen.

Tidak ada tabel `images` generik terpisah seperti contoh sebelumnya — tiap domain sudah punya kolom foto/gambar sendiri:

| Tabel                       | Kolom                | Isi (`object_key`)                                   | Strategi akses                                                                                     |
| --------------------------- | -------------------- | ---------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `users`                     | `photo_url`          | foto profil pengguna                                 | Private → presigned GET (±15 menit)                                                                |
| `children`                  | `photo_url`          | foto profil anak                                     | Private → presigned GET                                                                            |
| `child_photos`              | `photo_url`          | foto harian anak (feed, volume tinggi)               | Private → presigned GET. Kontrol tambahan lewat kolom `visibility` (enum) dan tabel `photo_shares` |
| `ingredients`               | `image_url`          | foto bahan makanan                                   | Public → `baseURL + object_key`                                                                    |
| `cooking_steps`             | `image_url`          | foto tiap langkah masak                              | Public → `baseURL + object_key`                                                                    |
| `assessment_kpsp_questions` | `image_url` _(baru)_ | ilustrasi/gambar pendukung pertanyaan KPSP, nullable | Public → `baseURL + object_key`                                                                    |

**Alasan menyimpan `object_key`, bukan URL penuh:**

- Kalau custom domain/CDN berganti di kemudian hari, tidak perlu migrasi seluruh baris URL di database.
- Untuk gambar **publik** (bahan makanan, langkah masak, ilustrasi KPSP) → URL dihitung saat response: `baseURL + object_key`.
- Untuk gambar **privat/sensitif** (foto profil, dan terutama foto harian anak) → generate presigned GET URL secara on-demand dengan masa berlaku pendek (misal 15 menit), bukan menyimpan URL permanen.

> **Rekomendasi khusus foto harian anak:** gunakan bucket **privat** + presigned GET URL, bukan public bucket, agar URL foto tidak bisa ditebak/diakses pihak luar. Skema ini sudah didukung tambahan kontrol lewat kolom `visibility` (`all` / `private` / `selected_only`) dan tabel junction `photo_shares` untuk berbagi ke kontak tertentu.

### Perubahan Migration Terkait

Skema di atas mengacu ke tabel yang sudah ada di `000004_create_core_domain_tables_up.sql` dan `000005_simplify_medical_module_up.sql`. Karena migration ini diasumsikan belum pernah dijalankan di environment manapun (masih tahap development), perubahan dilakukan dengan **mengedit langsung** file migration tersebut — bukan menambah migration `ALTER TABLE` baru.

**`000004_create_core_domain_tables_up.sql`** — tambahkan kolom `image_url TEXT` (nullable) pada `assessment_kpsp_questions`:

```diff
 CREATE TABLE assessment_kpsp_questions (
     id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     month_target            INTEGER NOT NULL CHECK (month_target IN (3,6,9,12,15,18,21,24,30,36,42,48,54,60)),
     developmental_domain    developmental_domain NOT NULL,
     question_text           TEXT NOT NULL,
+    image_url               TEXT,
     created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
     updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
 );
```

Tabel lain (`users`, `children`, `child_photos`, `ingredients`, `cooking_steps`) **tidak perlu diubah** — kolom `photo_url`/`image_url` sudah ada, cukup pastikan backend menulis `object_key` (bukan URL) ke kolom tersebut.

**`000005_simplify_medical_module_up.sql`** — tidak ada perubahan. Modul medical tidak menyentuh tabel media apa pun.

> ⚠️ Kalau `000004` sudah pernah dijalankan di environment manapun (staging/production), pendekatan edit-langsung ini **tidak berlaku lagi** — harus pakai migration baru (`ALTER TABLE ... ADD COLUMN`) agar riwayat migration tetap sinkron dengan state database aktual.

Gap yang belum ditangani di luar perubahan ini: tabel `recipes` masih belum punya kolom cover image sendiri (yang ada baru per-langkah lewat `cooking_steps.image_url`).

## 8. Pertimbangan Keamanan

- Validasi `content_type` dan ukuran file **di sisi backend** sebelum generate presigned URL — jangan percaya input dari client begitu saja.
- Set masa berlaku (expiry) presigned PUT URL singkat, misalnya 5–10 menit.
- Simpan R2 access key & secret hanya di backend — jangan pernah expose ke aplikasi Flutter.
- Jika ada dukungan Flutter Web, atur CORS policy di sisi R2 bucket agar browser bisa melakukan `PUT` langsung.

---

## 9. Sisi Flutter (Clean Architecture)

- `data/datasource/remote`: `ImageRemoteDataSource` — memanggil endpoint presign & confirm, lalu melakukan `PUT` langsung ke presigned URL (bukan lewat backend Go).
- `domain/usecase`: `UploadImageUseCase` — mengorkestrasi alur: `getPresignUrl → putFile → confirmUpload`.
- Disarankan memakai `dio` (dibanding `http`) karena mendukung progress callback — berguna untuk UX upload foto harian yang sering dilakukan berkali-kali/multi-file.

---

## 10. Referensi Dokumentasi Resmi Cloudflare R2

- Get Started: https://developers.cloudflare.com/r2/get-started/
- S3 API Compatibility: https://developers.cloudflare.com/r2/api/s3/api/
- Presigned URLs: https://developers.cloudflare.com/r2/api/s3/presigned-urls/
- Contoh aws-sdk-go: https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/
- Public Buckets: https://developers.cloudflare.com/r2/buckets/public-buckets/
- Konfigurasi CORS: https://developers.cloudflare.com/r2/buckets/cors/

---

## 11. Checklist Implementasi

- [ ] Buat bucket R2 di dashboard Cloudflare + catat Account ID
- [ ] Buat R2 API Token (access key & secret) dengan permission sesuai kebutuhan (read/write)
- [ ] Implementasi `ObjectStorage` interface di layer infrastructure (Go)
- [ ] Implementasi endpoint `POST /images/presign-upload`
- [ ] Implementasi endpoint `POST /images/confirm`
- [ ] Buat tabel `images` di database
- [ ] (Opsional) Konfigurasi custom domain untuk gambar publik (produk/resep)
- [ ] Pastikan bucket foto anak tetap privat + gunakan presigned GET saat serve ke Flutter
- [ ] Implementasi `ImageRemoteDataSource` & `UploadImageUseCase` di Flutter
- [ ] Setup CORS di R2 jika ada dukungan Flutter Web
- [ ] Uji retry upload saat koneksi terputus di tengah proses

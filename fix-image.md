# Solusi: Gambar R2 Gagal / Lambat Termuat (Mode `r2.dev`, Tanpa Custom Domain)

> Dokumen ini merangkum **mekanisme** perbaikan untuk kasus gambar dari Cloudflare R2 yang gagal atau lambat termuat di Flutter — termasuk sejak render pertama kali masuk halaman. Tidak berisi source code, fokus pada penyebab dan alur solusinya.

---

## 1. Ringkasan Masalah

Gambar (foto profil user, foto profil anak, dst.) yang sudah tersimpan di R2 kadang gagal termuat atau sangat lambat — bahkan pada render pertama saat halaman baru dibuka, bukan hanya saat render ulang.

---

## 2. Akar Masalah

### 2.1 Penyebab Utama — `object_key` Ditampilkan Langsung sebagai URL

Sesuai keputusan skema database sebelumnya, kolom `photo_url` di database **menyimpan `object_key`** (mis. `profile/{user_id}/{uuid}.webp`), **bukan URL lengkap** — ini prinsip yang memang disengaja agar tidak perlu migrasi data kalau domain/CDN berganti.

Masalahnya: di sisi Flutter, nilai `photoUrl` dari response API (yang sebenarnya `object_key` mentah) langsung dipakai sebagai argumen widget gambar berbasis jaringan, tanpa pernah dikonversi dulu menjadi URL yang bisa diakses (`https://...`). Widget gambar berbasis jaringan di Flutter mensyaratkan URI **absolut** — kalau diberi string relatif seperti `profile/uuid/xxx.webp`, permintaan gagal **sejak awal**, bukan soal lambat. Ini cocok persis dengan gejala: gagal bahkan di render pertama.

Titik penyebab ini ditemukan di lebih dari satu halaman (edit profil user, edit profil caregiver, form tambah/edit profil anak) — masing-masing memanggil widget gambar secara manual dan terpisah, sehingga potensi lupa/tidak konsisten menjadi tinggi.

### 2.2 Faktor Kontribusi — Karakteristik Resmi `r2.dev`

Setelah bug di atas diperbaiki pun, karena kondisi Anda memakai subdomain publik `r2.dev` (bukan custom domain), ada batasan resmi dari Cloudflare yang tetap relevan:

- **Tidak melewati Cloudflare Cache** — fitur caching, WAF, dan kontrol akses hanya tersedia lewat custom domain, tidak tersedia di `r2.dev`.
- **Rate-limited** — request ke `r2.dev` dibatasi (skala ratusan request/detik); saat melebihi batas, permintaan akan mendapat respons "terlalu banyak permintaan" dan throughput bisa ditahan sementara.
- **Ditujukan untuk testing/development**, bukan traffic produksi bervolume tinggi — relevan untuk fitur seperti feed foto harian anak yang memang didesain bervolume tinggi.
- **Tanpa autentikasi sama sekali** — begitu diaktifkan, semua objek di bucket bisa diakses siapa pun yang tahu/menebak object key-nya, tanpa proses login/otorisasi apa pun.

---

## 3. Mekanisme Solusi

### 3.1 Perbaikan Wajib — Resolusi `object_key` menjadi URL

Prinsipnya: **tidak ada satu pun tempat di aplikasi yang boleh menampilkan `object_key` mentah sebagai sumber gambar.** Konversi ke URL harus terjadi di satu titik terpusat, bukan tersebar di tiap halaman.

**Opsi A — Backend yang mengonversi (disarankan, sesuai prinsip skema semula):**
Backend menyimpan konfigurasi base URL publik R2 (`pub-<hash>.r2.dev`). Setiap kali backend mengembalikan data yang mengandung field foto (respons profil user, profil anak, dst.), backend menggabungkan base URL tersebut dengan `object_key` sebelum dikirim ke Flutter — sehingga field yang diterima Flutter sudah berupa URL utuh yang siap dipakai, bukan `object_key` mentah. Ini sejalan dengan prinsip awal: "URL dihitung saat response", dan menghindari duplikasi logika konversi di banyak halaman Flutter sekaligus mengunci base URL storage di satu tempat yang mudah diganti kelak.

**Opsi B — Flutter yang mengonversi (quick-fix sementara tanpa menunggu perubahan backend):**
Buat satu fungsi utilitas terpusat yang menerima nilai `photoUrl` dari model manapun, lalu: jika nilai tersebut sudah berupa URL absolut (dimulai `http`), pakai apa adanya (untuk kompatibilitas ke depan bila backend sudah dikonversi); jika bukan, gabungkan dengan base URL `r2.dev` yang disimpan sebagai konfigurasi aplikasi. Fungsi ini dipakai di **semua** titik yang menampilkan gambar dari R2, menggantikan pemanggilan widget gambar langsung dengan `photoUrl` mentah yang tersebar di berbagai halaman saat ini.

> Kedua opsi bisa dipakai sekaligus sebagai jaring pengaman selama masa transisi — backend melakukan konversi utama, Flutter tetap punya fallback kalau suatu saat ada data lama yang belum terkonversi.

### 3.2 Setup Public Development URL di Cloudflare

Karena tidak memakai custom domain, langkah aktivasi publik dilakukan lewat dashboard bucket R2: aktifkan opsi Public Development URL pada bucket yang relevan, lalu catat URL yang dihasilkan (`pub-<hash>.r2.dev`) sebagai konfigurasi base URL — baik untuk dipakai backend (Opsi A) maupun Flutter (Opsi B). Base URL ini sebaiknya disimpan sebagai environment/config variable, bukan ditulis langsung (hardcode) di kode, agar mudah diganti saat suatu saat pindah ke custom domain.

### 3.3 Mitigasi Karakteristik `r2.dev`

Karena batasan rate-limit dan tidak ada caching di `r2.dev` tetap berlaku meski bug di atas sudah diperbaiki, tambahan mitigasi berikut membantu mengurangi kemungkinan lambat/gagal terutama pada fitur bervolume tinggi (feed foto harian anak):

- **Caching gambar di sisi Flutter** — pakai pendekatan image caching berbasis disk (bukan widget gambar polos bawaan) supaya gambar yang sudah pernah dimuat tidak diminta ulang ke `r2.dev` setiap kali widget dibangun ulang atau halaman dibuka kembali. Ini langsung mengurangi jumlah request ke `r2.dev`, sehingga risiko kena rate-limit menurun.
- **Kompresi/pengecilan gambar sebelum upload** — semakin kecil ukuran file yang tersimpan, semakin cepat diambil kembali dan semakin kecil peluang throughput ditahan saat diminta banyak sekaligus (relevan untuk feed foto harian anak).
- **Penanganan error & retry di widget gambar** — tampilkan placeholder/ikon fallback saat gambar gagal dimuat (alih-alih layar rusak), dan sediakan mekanisme retry manual (mis. tap untuk memuat ulang) untuk mengatasi kegagalan sementara akibat throttling.
- **Pantau frekuensi akses** — kalau fitur foto harian anak benar-benar bervolume tinggi sesuai use case awal, pertimbangkan jadwal migrasi ke custom domain sebagai solusi jangka panjang, karena `r2.dev` secara resmi memang tidak ditujukan untuk traffic produksi.

---

## 4. Catatan Keamanan Khusus Foto Anak

Mengaktifkan `r2.dev` membuat bucket tersebut **publik tanpa autentikasi** — siapa pun dengan `object_key` yang benar bisa mengakses filenya langsung, tanpa perlu login. Ini berbeda dari rekomendasi awal (bucket privat + presigned GET khusus foto harian anak, karena sifatnya sensitif).

Kalau solusi `r2.dev` ini diterapkan untuk **semua** kategori foto termasuk foto harian anak, konsekuensinya keamanan foto anak menurun dibanding rencana awal. Kalau ingin tetap menjaga foto anak privat sambil kategori lain (foto profil, bahan makanan, dll.) tetap pakai `r2.dev` publik, pisahkan strateginya per bucket/kategori: hanya bucket untuk kategori non-sensitif yang diaktifkan Public Development URL, sedangkan bucket foto harian anak tetap privat dan diakses lewat presigned GET dari backend seperti rencana semula.

---

## 5. Checklist Perbaikan

- [ ] Aktifkan Public Development URL pada bucket R2 yang relevan, catat base URL `pub-<hash>.r2.dev`
- [ ] Simpan base URL tersebut sebagai konfigurasi (bukan hardcode), baik di backend maupun Flutter
- [ ] Implementasikan konversi `object_key → URL` di backend saat serialisasi response (Opsi A)
- [ ] (Opsional, jaring pengaman) implementasikan fungsi utilitas konversi terpusat di Flutter (Opsi B)
- [ ] Ganti seluruh pemanggilan widget gambar yang langsung memakai `photoUrl` mentah di semua halaman (edit profil user, edit profil caregiver, form profil anak, dan halaman lain yang menampilkan foto) agar melalui hasil konversi/fungsi utilitas tersebut
- [ ] Tambahkan image caching berbasis disk di sisi Flutter untuk semua gambar dari R2
- [ ] Tambahkan fallback/placeholder + retry pada widget gambar untuk menangani kegagalan sementara
- [ ] Evaluasi ulang strategi privat/publik khusus foto harian anak (pisahkan bucket kalau perlu)
- [ ] Rencanakan migrasi ke custom domain untuk jangka panjang, terutama kalau volume foto harian anak terus bertambah

---

## 6. Referensi Dokumentasi Resmi

- Cloudflare R2 — Public buckets (perbedaan `r2.dev` vs custom domain, batasan caching): https://developers.cloudflare.com/r2/buckets/public-buckets/
- Cloudflare R2 — Platform limits (rate limiting pada `r2.dev`): https://developers.cloudflare.com/r2/platform/limits/

---

## 7. Status Implementasi Aktual (Diselesaikan)

Berdasarkan permintaan Anda pada poin Opsi A yang dieksekusi (serta mengabaikan catatan privasi pada poin 4), berikut adalah langkah riil yang telah saya selesaikan di **Backend**:

1. Menyesuaikan injektor _storage_ dari skema Presigned URL menjadi R2 Public URL secara global untuk model Child dan User.
2. Mengganti Dependency `storage.ObjectStorage` dengan `cfg.R2PublicURL` pada layer `UserService` dan `ChildService`.
3. Mengganti utilitas `ResolvePhotoURL` (yang memproduksi path relatif / presigned URL berdurasi) dengan utilitas `ResolvePublicPhotoURL` yang memastikan **semua URL respon (termasuk untuk foto profil Caregiver, foto Ibu, dan profil Anak) langsung menjadi URL `http` penuh/absolut dengan base domain `r2.dev`**.
4. **Normalisasi Data Seed Database**: Karena backend sekarang proaktif menyematkan _Base Domain Public_ pada URL di setiap _Request_ yang membawa _response_ foto, skrip seed lokal PostgreSQL (seperti `ingredients.sql`, `children.sql`, `child_nutrition_derivative.sql`, dan `asesmen_kpsp_recom_actions.sql`) yang sebelumnya memuat statis URL `https://...` lengkap, secara retrospektif **telah dinormalkan/dibersihkan** agar murni kembali menyimpan _relative object_key_ saja (contoh: `profile/anak/m-razky.webp` atau `ingredients/GR048.jpg`). Hal ini menjaga agar database mutlak tidak mengandung base host S3 manapun sesuai poin `2.1` di atas.

Dengan langkah ini, Backend kita sepenuhnya mereturn URL absolut publik yang sudah jadi untuk seluruh konten gambar, dengan demikian memecahkan masalah widget network image Flutter secara instan layaknya di resep dan pertanyaan modul KPSP.

# NusaGizi — Dokumentasi Endpoint API (v4)

Dokumen ini adalah revisi dari `NusaGizi_Endpoint_Documentation_v3.md`, berdasarkan review langsung per endpoint. Ringkasan perubahan v3 -> v4 ada di tabel paling bawah dokumen. Riwayat sebelumnya:

1. **v3**: standar dokumentasi per endpoint dibuat konsisten penuh (_Path Params_, _Query Params_, _Request Body_, _Response Body_, _Response Error_ di semua endpoint), endpoint baru nomor 47–67 ditambahkan untuk menutup celah CRUD.
2. **Tidak ada pagination** — riwayat/list difilter per-bulan (`month`[+`year`]) atau age-range, bukan `page`/`limit`.
3. **Endpoint autentikasi (login/register/refresh/logout) sengaja tidak dibuat** — sepenuhnya ditangani Auth0 di luar API ini.

Nomor endpoint asli dari v1/v2 dipertahankan agar mudah dibandingkan; endpoint yang isinya diubah ditandai **[FIX]**, endpoint baru ditandai **[BARU]**.

---

## 0. Konvensi Umum

- **Auth**: ditangani provider eksternal (Auth0). Setiap request terautentikasi membawa `user_id` dari token. Role (mother/caregiver) ditentukan dari ada-tidaknya row terkait di `mother_profiles`/`caregiver_profiles` untuk `user_id` tsb.
- **Di luar cakupan dokumen ini**: login, register, refresh token, logout (semua ditangani Auth0). Pembuatan row `users` diasumsikan terjadi otomatis (mis. via Auth0 Action/webhook) saat signup pertama kali — dokumen ini mulai dari endpoint pertama yang dipanggil aplikasi **setelah** row `users` ada, yaitu pemilihan role (lihat endpoint 4/48 di bawah).
- **Format ID**: semua ID adalah UUID (string), termasuk di request body — bukan integer.
- **Format tanggal** (`date`, contoh kolom `birth_date`, `measured_at`, `valid_until`): `DD-MM-YYYY` (contoh: `25-07-2026`).
- **Format timestamp** (`datetime`, contoh kolom `created_at`, `updated_at`): ISO 8601 UTC (contoh: `2026-07-25T09:30:00Z`).
- **Penamaan field**: English snake_case, mengikuti nama kolom database (contoh: `full_name`, `birth_date`).
- **Soft delete**: hanya untuk tabel yang punya kolom `deleted_at` (`users`, `children`, `caregiver_engagements`) — endpoint DELETE terkait melakukan `UPDATE ... SET deleted_at = now()`. Tabel lain (`child_growth_reports`, `child_development_reports`, `child_photos`, `medical_notes`, `contacts`, `notifications`, `ingredient_shopping_items`, dst) **tidak** punya `deleted_at`, sehingga endpoint DELETE-nya melakukan hard delete (baris benar-benar dihapus).
- **Response POST**: kecuali dinyatakan lain, endpoint POST yang membuat resource baru mengembalikan minimal `{ "id": "<uuid>" }` (status `201 Created`), supaya client tidak perlu re-fetch.
- **Response PATCH/DELETE tanpa body**: status `200 OK` dengan body kosong (`-`), kecuali dinyatakan lain.
- **Tidak ada pagination**: endpoint list/riwayat difilter per periode (`month` + `year`) atau age-range, bukan `page`/`limit`. Jika suatu saat volume data jadi masalah, pagination bisa ditambahkan sebagai perubahan terpisah.
- **Base path**: contoh path di dokumen ini pakai konvensi resource-oriented (`/children/{child_id}/...`), sesuaikan dengan prefix API Anda (mis. `/api/v1/...`).
- **HTTP Header**: Seluruh _response_ API secara default sudah disisipkan header `Content-Type: application/json; charset=utf-8` (ditangani secara otomatis oleh framework Gin pada setiap pemanggilan `c.JSON()`).

### Format Error Standar

Semua endpoint yang gagal mengembalikan body dengan bentuk yang sama:

```json
{
    "error": {
        "code": "NOT_FOUND",
        "message": "Resource dengan id tsb tidak ditemukan"
    }
}
```

Status code yang dipakai di seluruh dokumen (di tiap endpoint hanya subset yang relevan yang dicantumkan):

| Status | `code`                 | Kapan terjadi                                                                           |
| ------ | ---------------------- | --------------------------------------------------------------------------------------- |
| 400    | `BAD_REQUEST`          | Body/query param tidak valid (format salah, tipe salah, field wajib kosong)             |
| 401    | `UNAUTHORIZED`         | Token tidak ada/tidak valid                                                             |
| 403    | `FORBIDDEN`            | Resource yang diminta bukan milik `user_id` yang login                                  |
| 404    | `NOT_FOUND`            | Resource dengan id di path tidak ditemukan (atau sudah di-soft-delete)                  |
| 409    | `CONFLICT`             | Duplikasi data (mis. `mother_profiles` sudah ada untuk user ini, `email` sudah dipakai) |
| 422    | `UNPROCESSABLE_ENTITY` | Validasi bisnis gagal (mis. `visibility=only` tapi `list_visibility` kosong)            |

---

## 1. Modul: Auth & User Profile

> Alur singkat: user signup/login via Auth0 → row `users` tersedia (di luar cakupan) → aplikasi tanya "Anda ibu atau caregiver?" → panggil endpoint 4 atau 48 untuk membuat profile sesuai role.

### 0. Konfirmasi Email (OTP) **[BARU]**

- **Method & Path**: `POST /confirm-email`
- **Authorization**: Tanpa token (public endpoint)
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: Wajib menyertakan `id_token` dari Auth0

```json
{
    "id_token": "eyJhbGciOiJ..."
}
```

| Field    | Type   | Wajib | Keterangan                                |
| -------- | ------ | ----- | ----------------------------------------- |
| id_token | string | Ya    | ID Token dari flow Passwordless OTP Auth0 |

- **Response Body (200)**:

```json
{
    "status": "verified"
}
```

- **Response Error**:

| Status | Kasus                                                |
| ------ | ---------------------------------------------------- |
| 400    | `id_token` wajib diisi atau klaim email tidak ada    |
| 401    | `id_token` tidak valid atau email belum diverifikasi |
| 404    | akun database untuk email ini tidak ditemukan        |
| 500    | gagal memperbarui status verifikasi email            |

### 1. Update profil user

- **Method & Path**: `PATCH /users`
- **Authorization**: `user_id` diambil dari token
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: semua field opsional, kirim hanya yang ingin diubah

```json
{
    "full_name": "ibor alwan",
    "email": "robialwan8@gmail.com",
    "phone_number": "081234567899",
    "gender": "male"
}
```

| Field        | Type   | Wajib | Keterangan                  |
| ------------ | ------ | ----- | --------------------------- |
| full_name    | string | Tidak | Maks 150 karakter           |
| email        | string | Tidak | Harus unik di tabel `users` |
| phone_number | string | Tidak | Maks 20 karakter            |
| gender       | enum   | Tidak | `male` / `female`           |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                     |
| ------ | ----------------------------------------- |
| 400    | `email` bukan format email valid          |
| 403    | Tidak relevan (sudah otomatis dari token) |
| 404    | `user_id` tidak ditemukan                 |
| 409    | `email` sudah dipakai user lain           |

### 2. Hapus user (soft delete)

- **Method & Path**: `DELETE /users`
- **Authorization**: `user_id` diambil dari token
- **Path Params**: `-` (tidak ada)

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                     |
| ------ | ----------------------------------------- |
| 403    | Tidak relevan (sudah otomatis dari token) |
| 404    | `user_id` tidak ditemukan / sudah dihapus |

- **Catatan**: soft delete — set `users.deleted_at = now()`. Perlu didiskusikan terpisah (di luar scope dokumen ini) apakah `children`/data terkait ikut di-cascade-soft-delete atau tetap ada.

### 3. [BARU] Get profil user sendiri

- **Method & Path**: `GET /users`
- **Authorization**: `user_id` diambil dari token
- **Path Params**: `-` (tidak ada)

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
{
    "id": "8e3b...",
    "full_name": "ibor alwan",
    "gender": "male",
    "email": "robialwan8@gmail.com",
    "phone_number": "081234567899",
    "photo_url": "https://...",
    "has_mother_profile": true,
    "has_caregiver_profile": false
}
```

| Field                 | Type           | Keterangan                                                                                                                                         |
| --------------------- | -------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| id                    | UUID           | `users.id`                                                                                                                                         |
| full_name             | string         |                                                                                                                                                    |
| gender                | enum \| null   | `male` / `female`                                                                                                                                  |
| email                 | string         |                                                                                                                                                    |
| phone_number          | string \| null |                                                                                                                                                    |
| photo_url             | string \| null |                                                                                                                                                    |
| has_mother_profile    | boolean        | **[FIX ikutan]** true jika ada row di `mother_profiles` untuk user ini — dipakai FE untuk tahu perlu redirect ke role-selection (47/48) atau tidak |
| has_caregiver_profile | boolean        | idem, untuk `caregiver_profiles`                                                                                                                   |

- **Response Error**:

| Status | Kasus                                     |
| ------ | ----------------------------------------- |
| 403    | Tidak relevan (sudah otomatis dari token) |
| 404    | `user_id` tidak ditemukan                 |

- **Catatan**: ditambahkan agar form edit di endpoint 1 punya data awal untuk di-load. Field `has_mother_profile`/`has_caregiver_profile` ditambahkan di v3 supaya app tahu kapan harus mengarahkan user ke endpoint 4/48.

### 4. [BARU] Daftar sebagai mother (pilih role)

- **Method & Path**: `POST /mother-profiles`
- **Authorization**: `user_id` diambil dari token (bukan dari body)
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada — cukup token)
- **Response Body (201)**:

```json
{ "id": "<mother_profile_id>" }
```

- **Response Error**:

| Status | Kasus                                                                    |
| ------ | ------------------------------------------------------------------------ |
| 409    | `mother_profiles` untuk `user_id` ini sudah ada (kolom `user_id` UNIQUE) |

- **Catatan [ASUMSI]**: skema `mother_profiles`/`caregiver_profiles` secara teknis tidak saling eksklusif (satu `user_id` bisa punya baris di keduanya sekaligus). Jika bisnis mewajibkan 1 user = 1 role saja, validasi itu perlu ditambahkan di application layer (bukan dari DB constraint).

### 5. [BARU] Daftar sebagai caregiver (pilih role)

- **Method & Path**: `POST /caregiver-profiles`
- **Authorization**: `user_id` diambil dari token
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada — cukup token)
- **Response Body (201)**:

```json
{ "id": "<caregiver_profile_id>" }
```

- **Response Error**:

| Status | Kasus                                              |
| ------ | -------------------------------------------------- |
| 409    | `caregiver_profiles` untuk `user_id` ini sudah ada |

### 6. [BARU] Get profil mother sendiri

- **Method & Path**: `GET /mother-profiles`
- **Authorization**: `user_id` diambil dari token
- **Path Params**: `-` (tidak ada)

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
{
    "id": "<mother_profile_id>",
    "user_id": "<user_id>",
    "full_name": "...",
    "email": "...",
    "phone_number": "...",
    "photo_url": "..."
}
```

| Field                                     | Type   | Keterangan                                                                                    |
| ----------------------------------------- | ------ | --------------------------------------------------------------------------------------------- |
| id                                        | UUID   | `mother_profiles.id`                                                                          |
| user_id                                   | UUID   |                                                                                               |
| full_name, email, phone_number, photo_url | string | join ke `users` — `mother_profiles` sendiri tidak punya kolom selain `id`/`user_id`/timestamp |

- **Response Error**:

| Status | Kasus           |
| ------ | --------------- |
| 403    | bukan pemilik   |
| 404    | tidak ditemukan |

- **Catatan**: paralel dengan endpoint 56 (`GET /caregiver-profiles/{id}`) yang sebelumnya sudah ada tapi belum ada versi mother-nya — celah ini ditutup di v3. endpoint 56 juga di-[FIX] responsnya (lihat Modul 7) supaya konsisten (ikut join ke `users`).

---

## 2. Modul: Child Profile

### 7. Tambah profil anak

- **Method & Path**: `POST /children`
- **Authorization**: `mother_profile_id` diambil dari `user_id` yang login (bukan dari body)
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{
    "full_name": "ibor",
    "birth_date": "22-07-2026",
    "gender": "male",
    "photo_url": "https://docs.google.com/",
    "allergies": {
        "food": ["Kacang", "Telur"],
        "medicine": ["Penisilin"],
        "animal": ["Kucing"],
        "others": ["Debu"]
    },
    "chronic_diseases": ["Diabetes", "PJB"],
    "diets": ["TETP"],
    "favorite_foods": ["Ayam", "Telur", "Naspad"],
    "favorite_textures": ["Finger food", "Makanan keluarga"],
    "notes": "Saya mau menu yang murah"
}
```

| Field                                          | Type            | Wajib | Keterangan                                                             |
| ---------------------------------------------- | --------------- | ----- | ---------------------------------------------------------------------- |
| full_name                                      | string          | Ya    | Maks 150 karakter                                                      |
| birth_date                                     | date            | Ya    | `DD-MM-YYYY`                                                           |
| gender                                         | enum            | Ya    | `male` / `female`                                                      |
| photo_url                                      | string          | Tidak |                                                                        |
| allergies.food / .medicine / .animal / .others | array\<string\> | Tidak | tiap item -> row `child_allergy_profiles` dengan `category` sesuai key |
| chronic_diseases                               | array\<string\> | Tidak | -> row `child_chronic_disease_profiles.disease_name`                   |
| diets                                          | array\<string\> | Tidak | -> row `child_diet_profiles.diet_name`                                 |
| favorite_foods                                 | array\<string\> | Tidak | -> row `favorite_food_profiles.food_name`                              |
| favorite_textures                              | array\<string\> | Tidak | -> row `favorite_texture_profiles.texture_name`                        |
| notes                                          | string          | Tidak | -> `child.notes_profile`                                               |

- **Response Body (201)**: `{ "id": "<child_id>" }`
- **Response Error**:

| Status | Kasus                                                       |
| ------ | ----------------------------------------------------------- |
| 400    | `birth_date` bukan tanggal valid, `gender` bukan enum valid |
| 422    | `full_name`/`birth_date`/`gender` kosong                    |

- **Catatan**: Field diganti dari Indonesia (`name`, `date_of_birth`, `alergi`, `kondisi_kronis`, dst) ke English snake_case sesuai konvensi.

### 8. Update profil anak

- **Method & Path**: `PATCH /children/{child_id}`
- **Auth [FIX]**: divalidasi lewat `child_id` — pastikan `children.mother_profile_id` terhubung ke `mother_profiles.user_id` yang login. (Dokumentasi lama menyebut "berdasarkan user_id", tapi tabel `children` tidak punya kolom `user_id` sama sekali.)
- **Path Params**:

| Field    | Type | Keterangan |
| -------- | ---- | ---------- |
| child_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: sama seperti endpoint 7, semua field opsional (hanya kirim yang mau diubah)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                      |
| ------ | ------------------------------------------ |
| 403    | `child_id` bukan milik mother yang login   |
| 404    | `child_id` tidak ditemukan / sudah dihapus |

- **Catatan**: untuk field turunan (`allergies`, `chronic_diseases`, `diets`, `favorite_foods`, `favorite_textures`), strategi update disarankan replace-all (hapus semua row lama untuk child tsb, insert ulang dari array baru) supaya tidak perlu diff per-item.

### 9. [BARU] Hapus profil anak (soft delete)

- **Method & Path**: `DELETE /children/{child_id}`
- **Authorization**: pastikan `child.mother_profile_id` milik `user_id` yang login
- **Path Params**:

| Field    | Type | Keterangan |
| -------- | ---- | ---------- |
| child_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                           |
| ------ | ------------------------------- |
| 403    | bukan milik mother yang login   |
| 404    | tidak ditemukan / sudah dihapus |

- **Catatan**: soft delete — set `child.deleted_at = now()`.

### 10. [BARU] Get detail profil anak

- **Method & Path**: `GET /children/{child_id}/profile`
- **Authorization**: pastikan `child.mother_profile_id` milik `user_id` yang login (mother) atau ada `caregiver_engagements` aktif untuk `child_id` ini (caregiver)
- **Path Params**:

| Field    | Type | Keterangan |
| -------- | ---- | ---------- |
| child_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: bentuk sama persis dengan request body endpoint 7, plus `id` — supaya bisa langsung dipakai sebagai initial state form edit (endpoint 8):

```json
{
    "id": "<child_id>",
    "full_name": "ibor",
    "birth_date": "22-07-2026",
    "gender": "male",
    "photo_url": "https://docs.google.com/",
    "streak_days": 12,
    "allergies": {
        "food": ["Kacang", "Telur"],
        "medicine": ["Penisilin"],
        "animal": ["Kucing"],
        "others": ["Debu"]
    },
    "chronic_diseases": ["Diabetes", "PJB"],
    "diets": ["TETP"],
    "favorite_foods": ["Ayam", "Telur", "Naspad"],
    "favorite_textures": ["Finger food", "Makanan keluarga"],
    "notes": "Saya mau menu yang murah"
}
```

| Field        | Type    | Keterangan                          |
| ------------ | ------- | ----------------------------------- |
| streak_days  | integer | `child.streak_days`                 |
| (field lain) | —       | sama seperti tabel field endpoint 7 |

- **Response Error**:

| Status | Kasus                                                              |
| ------ | ------------------------------------------------------------------ |
| 403    | bukan mother pemilik / caregiver yang tidak punya engagement aktif |
| 404    | tidak ditemukan / sudah dihapus                                    |

- **Catatan**: celah dari v2 — sebelumnya tidak ada endpoint untuk mengambil profil anak secara lengkap (yang ada hanya ringkasan dashboard di endpoint 65, yang fieldnya jauh lebih sedikit).

### 11. [BARU] GetListChild (berdasarkan user login (mother/caregiver))

- **Method & Path**: `GET /children`
- **Authorization**: Mengambil data anak milik mother profile atau yang sedang diasuh oleh caregiver profile dari `user_id` yang login
- **Path Params**: `-` (tidak ada)

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
[
    {
        "id": "...",
        "full_name": "ibor",
        "age": "2 tahun 3 bulan",
        "gender": "male",
        "photo_url": "..."
    }
]
```

| Field       | Type           | Keterangan                |
| ----------- | -------------- | ------------------------- |
| id          | UUID           |                           |
| full_name   | string         |                           |
| age         | string         | format: "X tahun Y bulan" |
| gender      | enum           |                           |
| photo_url   | string \| null |                           |
| mother_name | string         |                           |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

- **Catatan**: versi ringan dari endpoint 65 (dashboard) — dipakai untuk kebutuhan seperti dropdown/switcher pilih anak, tanpa perlu join ke growth/nutrition/KPSP yang lebih berat. Bentuk response-nya sengaja disamakan dengan endpoint 57 (list anak versi caregiver) untuk konsistensi antar-role.

### 12. [BARU] GetChild (ringan)

- **Method & Path**: `GET /children/{child_id}`
- **Authorization**: Bebas (selama token valid) — tidak ada cek ownership / engagement (agar caregiver bisa melihat profil anak saat validasi QR checkin)
- **Path Params**:

| Field    | Type | Keterangan |
| -------- | ---- | ---------- |
| child_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
{
    "id": "<child_id>",
    "full_name": "ibor",
    "birth_date": "22-07-2026",
    "gender": "male",
    "photo_url": "https://docs.google.com/",
    "streak_days": 12,
    "mother_name": "Siti Rahayu"
}
```

- **Response Error**:

| Status | Kasus                                                              |
| ------ | ------------------------------------------------------------------ |
| 403    | bukan mother pemilik / caregiver yang tidak punya engagement aktif |
| 404    | tidak ditemukan / sudah dihapus                                    |

---

## 3. Modul: Child Growth

### 📖 Algoritma Perhitungan Status Pertumbuhan Anak (WHO LMS)

#### Ringkasan Algoritma LMS

Sistem menggunakan metode LMS (Lambda, Median, Coefficient of Variation) dari dataset WHO untuk menghitung Z-Score pertumbuhan anak.

**Rumus Z-Score:**

- Jika L ≠ 0: `Z = (((X / M)^L) - 1) / (L × S)`
- Jika L = 0: `Z = ln(X / M) / S`
  _(X = hasil pengukuran anak)_

> **Catatan (v4):** Z-Score kini **dihitung secara on-the-fly di server (service layer)** saat data diminta, bukan lagi disimpan di database. Tabel `child_growth_analyses` telah dihapus. Ini memastikan data grafik selalu valid tanpa perlu update row tersendiri setiap kali ada perbaikan data historis.

**Klasifikasi Status:**

- **Normal**: -2 ≤ Z ≤ +2
- **Berisiko**: -3 ≤ Z < -2 atau +2 < Z ≤ +3
- **Sangat Buruk**: Z < -3 atau Z > +3

---

### Indikator & Dataset WHO

| Indikator | Input Parameter      | Nilai X (Pengukuran) | Dataset                       | Parameter Pencarian  |
| --------- | -------------------- | -------------------- | ----------------------------- | -------------------- |
| **BB/U**  | Usia, BB             | Berat Badan          | WFA                           | Usia                 |
| **BB/TB** | Usia, BB, TB/PB      | Berat Badan          | WFL (<24 bln) / WFH (≥24 bln) | Tinggi/Panjang Badan |
| **TB/U**  | Usia, TB/PB          | Tinggi/Panjang Badan | LHFA                          | Usia                 |
| **LK/U**  | Usia, Lingkar Kepala | Lingkar Kepala       | HCFA                          | Usia                 |
| **IMT/U** | Usia, BB, TB         | BMI (BB / TB²)       | BFA                           | Usia                 |

_(Semua indikator juga membutuhkan parameter Jenis Kelamin untuk memilih tabel)_

**Alur:** Input → Tentukan Dataset → Cari Baris (berdasarkan Usia/Tinggi) → Ambil L, M, S → Hitung Z-Score → Klasifikasi Status

**Referensi:** WHO Child Growth Standards (2006)

### 13. Get child_growth_reports terbaru

- **Method & Path**: `GET /children/{child_id}/growth-reports/latest`
- **Authorization**: `child_growth_reports` sesuai `child_id`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
{
    "id": "...",
    "measured_at": "04-05-2026",
    "weight_kg": 12.4,
    "height_cm": 89,
    "head_circumference_cm": 47,
    "status": "Berisiko",
    "description": "Pertumbuhan anak berada di rentang berisiko pada salah satu atau lebih aspek pengukuran. Konsultasikan dengan dokter atau ahli gizi untuk evaluasi lebih lanjut."
}
```

| Field                 | Type   | Keterangan                                                                                |
| --------------------- | ------ | ----------------------------------------------------------------------------------------- |
| id                    | UUID   |                                                                                           |
| measured_at           | date   |                                                                                           |
| weight_kg             | number |                                                                                           |
| height_cm             | number |                                                                                           |
| head_circumference_cm | number |                                                                                           |
| status                | string | Hasil agregat terburuk dari 5 indikator WHO: `Normal`, `Berisiko`, `Gangguan Pertumbuhan` |

- **Catatan [FIX v5]**: Field `status` dihitung on-the-fly menggunakan Z-Score WHO LMS dari 5 indikator sekaligus. Status yang dikembalikan adalah yang terburuk (_worst-case_) di antara semua indikator:
    - **`Gangguan Pertumbuhan`**: Jika terdapat _minimal satu_ indikator dengan Z-Score < -3 atau > +3.
    - **`Berisiko`**: Jika terdapat _minimal satu_ indikator dengan Z-Score berada di rentang [-3, -2) atau (2, 3].
    - **`Normal`**: Jika _semua_ indikator memiliki Z-Score di dalam rentang aman [-2, +2].

    Indikator yang dievaluasi:
    1. Weight-for-Age (BB/U)
    2. Height-for-Age (PB/U atau TB/U)
    3. Weight-for-Height (BB/PB atau BB/TB)
    4. BMI-for-Age (IMT/U)
    5. Head-Circumference-for-Age (LK/U)

- **Response Error**:

| Status | Kasus                                                                       |
| ------ | --------------------------------------------------------------------------- |
| 403    | `child_id` bukan milik user yang login                                      |
| 404    | `child_id` tidak ditemukan, atau belum ada record growth report sama sekali |

### 14. Get child_growth_analyses (filter age-range & analysis_type) [FIX]

- **Method & Path**: `GET /children/{child_id}/growth-analyses`
- **Authorization**: `child_growth_analyses` yang `child_growth_report_id`-nya milik `child_id` ybs
- **Path Params**: `child_id` (UUID)
- **Query Params [FIX v4: dijadikan wajib]**:

| Field         | Type | Wajib | Keterangan                                                                                                                      |
| ------------- | ---- | ----- | ------------------------------------------------------------------------------------------------------------------------------- |
| analysis_type | enum | Ya    | salah satu dari `weight_for_age`, `height_for_age`, `weight_for_height`, `bmi_for_age`, `head_circumference_for_age`            |
| age_range     | enum | Ya    | `0-2`, `0-12`, atau `0-60` (satuan **bulan usia anak**, dihitung dari `child.birth_date` ke `child_growth_reports.measured_at`) |

- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**:

```json
{
    "data_points": [
        { "x": 2, "y": 5.5 },
        { "x": 3, "y": 6.2 },
        { "x": 4, "y": 7.1 }
    ],
    "message": {
        "title": "Normal",
        "description": "Kondisi anak berada dalam rentang normal. Pertahankan pola asuh dan nutrisi yang baik."
    }
}
```

| Field               | Type            | Keterangan                                                                                 |
| ------------------- | --------------- | ------------------------------------------------------------------------------------------ |
| data_points         | array of object | Berisi koordinat grafik pertumbuhan                                                        |
| data_points[].x     | number          | Sumbu X: biasanya usia (bulan), kecuali `weight_for_height` (menggunakan tinggi badan, cm) |
| data_points[].y     | number          | Sumbu Y: hasil pengukuran (kg atau cm) sesuai `analysis_type`                              |
| message             | object \| null  | Kesimpulan status pertumbuhan dari titik terakhir (latest measurement)                     |
| message.title       | string          | Klasifikasi status (contoh: "Normal", "Berisiko")                                          |
| message.description | string          | Deskripsi atau saran singkat terkait status pertumbuhan saat ini                           |

- **Response Error**:

| Status | Kasus                                                                        |
| ------ | ---------------------------------------------------------------------------- |
| 400    | `analysis_type`/`age_range` tidak dikirim, atau bukan salah satu nilai valid |
| 403    | `child_id` bukan milik user yang login                                       |
| 404    | `child_id` tidak ditemukan                                                   |

- **Catatan [FIX v4]**: Z-Score sekarang dihitung secara on-the-fly dan menghasilkan respons berbentuk JSON object (`data_points` dan `message`). Sumbu X pada `data_points` menyesuaikan indikator: untuk indikator BB/TB, sumbu X adalah panjang/tinggi badan (`height_cm`), sedangkan indikator lain menggunakan usia dalam bulan (`age_months`). Pesan (`message`) didapatkan dengan menghitung Z-Score dari titik pengukuran paling terakhir.

### 15. Get riwayat child_growth_reports

- **Method & Path**: `GET /children/{child_id}/growth-reports`
- **Authorization**: sesuai `child_id`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada — dikembalikan semua, diurutkan `measured_at` terbaru dulu)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
[
    {
        "id": "...",
        "measured_at": "04-05-2026",
        "weight_kg": 12.4,
        "height_cm": 89,
        "head_circumference_cm": 47,
        "status": "Normal",
        "description": "Pertumbuhan anak berada dalam rentang normal di semua aspek. Pertahankan pola makan seimbang dan pemantauan rutin."
    }
]
```

- **Catatan [FIX v5]**: Setiap item dalam array kini menyertakan `status` dan `description` yang dihitung on-the-fly (aturan sama seperti endpoint 13). Tidak ada pagination — semua riwayat dikembalikan sekaligus.

- **Response Error**:

| Status | Kasus                                  |
| ------ | -------------------------------------- |
| 403    | `child_id` bukan milik user yang login |
| 404    | `child_id` tidak ditemukan             |

### 16. Tambah child_growth_reports

- **Method & Path**: `POST /children/{child_id}/growth-reports`
- **Authorization**: hanya menyimpan `child_id` yang sesuai
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX]**:

```json
{
    "measured_at": "04-05-2026",
    "height_cm": 89,
    "weight_kg": 12.4,
    "head_circumference_cm": 47
}
```

| Field                 | Type   | Wajib | Keterangan     |
| --------------------- | ------ | ----- | -------------- |
| measured_at           | date   | Ya    | `DD-MM-YYYY`   |
| height_cm             | number | Tidak | > 0 jika diisi |
| weight_kg             | number | Tidak | > 0 jika diisi |
| head_circumference_cm | number | Tidak | > 0 jika diisi |

- **Response Body (201)**: `{ "id": "<growth_report_id>" }`
- **Response Error**:

| Status | Kasus                                                                            |
| ------ | -------------------------------------------------------------------------------- |
| 400    | `measured_at` bukan tanggal valid                                                |
| 422    | salah satu dari `height_cm`/`weight_kg`/`head_circumference_cm` diisi angka <= 0 |
| 403    | `child_id` bukan milik user yang login                                           |

### 17. [BARU] Edit child_growth_reports

- **Method & Path**: `PATCH /growth-reports/{child_growth_report_id}`
- **Authorization**: `child_growth_report_id` -> `child_id` -> `mother_profile_id` harus milik user yang login
- **Path Params**:

| Field                  | Type | Keterangan |
| ---------------------- | ---- | ---------- |
| child_growth_report_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: sama seperti endpoint 16, semua field opsional
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |
| 422    | nilai <= 0                  |

- **Catatan**: dipakai untuk koreksi salah input (mis. salah ketik berat badan). Pembaruan data pengukuran akan otomatis tercermin pada status yang dihitung ulang saat data berikutnya diminta.

### 18. [BARU] Hapus child_growth_reports

- **Method & Path**: `DELETE /growth-reports/{child_growth_report_id}`
- **Authorization**: sama seperti endpoint 17
- **Path Params**:

| Field                  | Type | Keterangan |
| ---------------------- | ---- | ---------- |
| child_growth_report_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: hard delete (`child_growth_reports` tidak punya `deleted_at`).

---

## 4. Modul: Child Development (KPSP)

**Algoritma bersama untuk endpoint 19(GET), 23(CREATE), 24(UPDATE)**.

```
KPSP_PERIODS = [3, 6, 9, 12, 15, 18, 21, 24, 30, 36, 42, 48, 54, 60]  // bulan
age_months   = umur anak sekarang, dihitung dari birth_date

month_target = periode TERBESAR di KPSP_PERIODS yang <= age_months
               jika age_months < 3 -> month_target = 3

kpsp_score   = jumlah jawaban "true" dari 10 pertanyaan (skala 0-10)

status =
    jika kpsp_score 9-10 -> "Sesuai Usia"
    jika kpsp_score 7-8  -> "Perkembangan meragukan"
    jika kpsp_score <= 6 -> "Kemungkinan penyimpangan"

next_period  = periode TERKECIL di KPSP_PERIODS yang > age_months
next_check_date =
    jika next_period ada  -> birth_date + next_period bulan   (format DD-MM-YYYY)
    jika tidak ada (age_months >= 60, periode terakhir sudah dikerjakan)
                          -> literal string "done"
```

Contoh: usia 7 bulan -> `month_target=6`, `next_check_date` = usia 9 bulan (2 bulan lagi). Usia 2 bulan -> `month_target=3`, `next_check_date` = usia 3 bulan (bulan depan).

### 19. Get child_development_reports terbaru + 4 domain

- **Method & Path**: `GET /children/{child_id}/development-reports/latest`
- **Authorization**: sesuai `child_id`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
{
    "id": "...",
    "kpsp_score": 8,
    "next_check_date": "25-08-2026",
    "status": "Sesuai Usia",
    "created_at": "...",
    "domains": [
        {
            "developmental_domain": "Gross motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "Fine motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "Speech and language",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "Socialization",
            "total_question": 3,
            "true_answer": 2
        }
    ]
}
```

| Field                          | Type           | Keterangan                                                                                 |
| ------------------------------ | -------------- | ------------------------------------------------------------------------------------------ |
| id                             | UUID           |                                                                                            |
| kpsp_score                     | integer        | 0-10                                                                                       |
| next_check_date                | date \| "done" |                                                                                            |
| status                         | string         | "Sesuai Usia", "Perkembangan meragukan", atau "Kemungkinan penyimpangan"                   |
| created_at                     | datetime       |                                                                                            |
| domains[].developmental_domain | enum           | 4 aspek: `Gross motor skills`, `Fine motor skills`, `Speech and language`, `Socialization` |
| domains[].total_question       | integer        | jumlah pertanyaan di domain ini untuk `month_target` ybs                                   |
| domains[].true_answer          | integer        | jumlah jawaban `true` di domain ini                                                        |

- **Response Error**:

| Status | Kasus                                             |
| ------ | ------------------------------------------------- |
| 403    | `child_id` bukan milik user yang login            |
| 404    | belum ada `child_development_reports` sama sekali |

- **Catatan [FIX v4]**: v3 mengembalikan raw list jawaban per domain lalu FE yang menghitung `total_question`/`true_answer`. Diubah jadi agregasi di query (`GROUP BY developmental_domain`, `COUNT(*)` untuk `total_question`, `COUNT(*) FILTER (WHERE answer = true)` untuk `true_answer`) — lebih baik dari sisi performa karena payload jauh lebih kecil dan logic agregasi cukup ditulis sekali di backend, tidak diulang di tiap platform client (mobile/web).

### 20. Get riwayat asesmen KPSP

- **Method & Path**: `GET /children/{child_id}/development-reports`
- **Authorization**: sesuai `child_id`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: Array dari riwayat hasil asesmen (list `child_development_reports`)

```json
[
    {
        "id": "...",
        "kpsp_score": 8,
        "month_target": 24,
        "status": "Sesuai Usia",
        "created_at": "..."
    }
]
```

| Field        | Type     | Keterangan                                                               |
| ------------ | -------- | ------------------------------------------------------------------------ |
| id           | UUID     |                                                                          |
| kpsp_score   | integer  |                                                                          |
| month_target | integer  |                                                                          |
| status       | string   | "Sesuai Usia", "Perkembangan meragukan", atau "Kemungkinan penyimpangan" |
| created_at   | datetime |                                                                          |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | `child_id` tidak ditemukan  |

### 21. Get detail hasil asesmen KPSP

- **Method & Path**: `GET /development-reports/{child_development_report_id}`
- **Authorization**: `child_development_reports` sesuai `child_id` milik mother/caregiver yang login
- **Path Params**: `child_development_report_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**:

```json
{
    "id": "...",
    "kpsp_score": 8,
    "month_target": 24,
    "status": "Sesuai Usia",
    "recommended_actions": [{ "id": "...", "title": "...", "action_text": "..." }],
    "domains": [
        {
            "developmental_domain": "Gross motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "Fine motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "Speech and language",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "developmental_domain": "Socialization",
            "total_question": 3,
            "true_answer": 2
        }
    ]
}
```

| Field                                                     | Type                 | Keterangan                                                                                    |
| --------------------------------------------------------- | -------------------- | --------------------------------------------------------------------------------------------- |
| id                                                        | UUID                 |                                                                                               |
| kpsp_score                                                | integer              |                                                                                               |
| month_target                                              | integer              | **[FIX v4]** menggantikan `next_check_date` (dihapus — tidak relevan di halaman detail hasil) |
| status                                                    | string               | "Sesuai Usia", "Perkembangan meragukan", atau "Kemungkinan penyimpangan"                      |
| recommended_actions[].id/title/action_text                | UUID/string/string   |                                                                                               |
| domains[].developmental_domain/total_question/true_answer | enum/integer/integer | **[FIX v4]** disamakan dengan format agregat endpoint 19                                      |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan [FIX v4]**: `recommended_actions` yang tampil **hanya rekomendasi aksi untuk domain (`developmental_domain`) yang memiliki performa terburuk** (yaitu domain dengan rasio `true_answer / total_question` paling kecil). Jika sebelumnya mengambil semua rekomendasi, kini difilter secara spesifik untuk domain yang paling butuh stimulasi.

### 22. Get pertanyaan assessment_kpsp_questions

- **Method & Path**: `GET /assessment-kpsp-questions`
- **Path Params**: `-` (tidak ada)
- **Query Params [FIX]**:

| Field        | Type    | Wajib | Keterangan                                                                                                                                   |
| ------------ | ------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| month_target | integer | Ya    | salah satu dari `3,6,9,12,15,18,21,24,30,36,42,48,54,60` — sebelumnya tidak ada filter di dokumentasi lama, padahal soal berbeda per periode |

- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: Array dari `assessment_kpsp_questions` — field `month_target` ditambahkan ke tiap item response

```json
[
    {
        "id": "...",
        "month_target": 3,
        "developmental_domain": "Gross motor skills",
        "question_text": "..."
    }
]
```

- **Response Error**:

| Status | Kasus                                       |
| ------ | ------------------------------------------- |
| 400    | `month_target` bukan salah satu nilai valid |

### 23. Simpan hasil asesmen KPSP (create)

- **Method & Path**: `POST /children/{child_id}/development-reports`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: `assessment_kpsp_answer` -> `answer`; `month_target` dikirim kembali oleh client]**:

```json
{
    "month_target": 24,
    "list_answer": [
        {
            "assessment_kpsp_question_id": "021f1136-d24a-4e32-aaa8-cca24b9e608d",
            "answer": true
        },
        {
            "assessment_kpsp_question_id": "03b996b5-1377-4a91-a0ba-5632526e1c09",
            "answer": false
        }
    ]
}
```

| Field                                     | Type            | Wajib | Keterangan                                             |
| ----------------------------------------- | --------------- | ----- | ------------------------------------------------------ |
| month_target                              | integer         | Ya    | salah satu dari 3,6,9,12,15,18,21,24,30,36,42,48,54,60 |
| list_answer                               | array\<object\> | Ya    | wajib 10 item (10 pertanyaan)                          |
| list_answer[].assessment_kpsp_question_id | UUID            | Ya    |                                                        |
| list_answer[].answer                      | boolean         | Ya    |                                                        |

- **Response Body (201)**: sama seperti endpoint 21 — field `id` di response **adalah** `child_development_report_id` yang baru dibuat
- **Response Error**:

| Status | Kasus                                                                                                                     |
| ------ | ------------------------------------------------------------------------------------------------------------------------- |
| 422    | `list_answer` bukan 10 item, atau ada `assessment_kpsp_question_id` yang tidak match periode `month_target` anak saat ini |
| 403    | `child_id` bukan milik user yang login                                                                                    |

- **Catatan [FIX v4]**: `month_target` dikirim oleh client di body request, lalu divalidasi terhadap `child.birth_date` memakai algoritma di atas modul ini; `kpsp_score` & `next_check_date` tetap dihitung 100% di server. Alur pembuatan ID: server men-generate UUID baru untuk `child_development_report_id` **sebelum** proses insert, lalu UUID yang sama itu dipakai sebagai foreign key `child_development_report_id` di setiap row `assessment_kpsp_answers` yang diinsert dalam satu transaksi — jadi client tidak perlu (dan tidak boleh) generate/kirim ID ini. Setelah semua 10 soal dijawab, sistem juga insert ke `development_report_recommendations` untuk tiap pertanyaan yang dijawab `false`.

### 24. Update hasil asesmen KPSP terbaru

- **Method & Path**: `PATCH /development-reports/{child_development_report_id}`
- **Path Params**: `child_development_report_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{
    "list_answer": [
        {
            "assessment_kpsp_question_id": "021f1136-d24a-4e32-aaa8-cca24b9e608d",
            "assessment_kpsp_answer": true
        }
    ]
}
```

| Field       | Type            | Wajib | Keterangan                                 |
| ----------- | --------------- | ----- | ------------------------------------------ |
| list_answer | array\<object\> | Ya    | boleh partial (hanya jawaban yang berubah) |

- **Response Body (200)**: sama seperti endpoint 21
- **Response Error**:

| Status | Kasus                                         |
| ------ | --------------------------------------------- |
| 403    | bukan milik user yang login                   |
| 404    | `child_development_report_id` tidak ditemukan |

- **Catatan**: `kpsp_score`, `next_check_date`, dan `development_report_recommendations` dihitung ulang. Dipanggil saat mother mengulangi asesmen.

### 25. Get checklist_milestone_tasks

- **Method & Path**: `GET /checklist-milestone-tasks`
- **Path Params**: `-` (tidak ada)
- **Query Params**:

| Field        | Type    | Wajib                     | Keterangan                           |
| ------------ | ------- | ------------------------- | ------------------------------------ |
| month_target | integer | Ya                        | salah satu dari 14 nilai valid       |
| child_id     | UUID    | Ya **[FIX: ditambahkan]** | dipakai untuk left join `is_checked` |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari `checklist_milestone_tasks`

```json
[
    {
        "id": "...",
        "developmental_domain": "Gross motor skills",
        "task_description": "...",
        "is_checked": false
    }
]
```

- **Response Error**:

| Status | Kasus                                  |
| ------ | -------------------------------------- |
| 400    | `month_target` invalid                 |
| 403    | `child_id` bukan milik user yang login |

- **Catatan**: `is_checked` didapat dari left join ke `checklist_milestone_progress` dengan `child_id` yang dikirim, supaya FE tahu item mana yang sudah dicentang tanpa request terpisah.

### 26. Get rekomendasi aksi KPSP berdasarkan hasil asesmen

- **Method & Path**: `GET /children/{child_id}/development-reports/{report_id}/recommendations`
- **Authorization**: sesuai `child_id`
- **Path Params**:

| Field     | Type | Keterangan                    |
| --------- | ---- | ----------------------------- |
| child_id  | UUID | ID anak                       |
| report_id | UUID | `child_development_report_id` |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
[
    {
        "developmental_domain": "Gross motor skills",
        "action_text": "Perlu latihan lebih sering untuk menyusun 3-4 balok mainan."
    },
    {
        "developmental_domain": "Speech and language",
        "action_text": "Perlu latihan lebih sering untuk mencoret kertas menggunakan krayon/pensil."
    }
]
```

| Field                | Type   | Keterangan                                                          |
| -------------------- | ------ | ------------------------------------------------------------------- |
| developmental_domain | string | Aspek perkembangan: `Gross motor skills`, `Fine motor skills`, dll. |
| action_text          | string | Teks rekomendasi latihan yang harus dilakukan orang tua             |

- **Response Error**:

| Status | Kasus                                  |
| ------ | -------------------------------------- |
| 403    | `child_id` bukan milik user yang login |
| 404    | `report_id` tidak ditemukan            |

- **Catatan**: Data diambil dari junction table `development_report_recommendations`, di-join ke `recommended_actions` dan `assessment_kpsp_questions` (alias `akq`) untuk mendapatkan `developmental_domain` dan `action_text`. Hanya pertanyaan yang dijawab `false` yang menghasilkan rekomendasi.

### 27. Tandai checklist milestone (Sync / UPSERT)

- **Method & Path**: `PATCH /children/{child_id}/checklist-milestone-progress`

- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: semantik jadi Sync/UPSERT, mencakup fungsi endpoint 55 lama]**:

```json
{
    "checklist_milestone_task_ids": ["b6f1e2d0-...-uuid-1", "8a2c3f10-...-uuid-2"]
}
```

| Field                        | Type          | Wajib | Keterangan                                                                                                                              |
| ---------------------------- | ------------- | ----- | --------------------------------------------------------------------------------------------------------------------------------------- |
| checklist_milestone_task_ids | array\<UUID\> | Ya    | daftar **lengkap** task yang tercentang saat ini untuk `child_id` ini (bukan hanya yang baru dicentang). Kirim `[]` untuk uncheck semua |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                                  |
| ------ | ------------------------------------------------------ |
| 403    | `child_id` bukan milik user yang login                 |
| 404    | ada `checklist_milestone_task_id` yang tidak ditemukan |

- **Catatan [FIX v4]**: Menggunakan semantik **Sync / UPSERT** (penyempurnaan dari pendekatan _replace-all_ mentah) — client selalu mengirim seluruh state checklist yang tercentang untuk `child_id` ini di tiap panggilan; server secara cerdas menghapus baris `checklist_milestone_progress` yang dilepas centangnya (uncheck) menggunakan logika selisih/diff, dan melakukan _Insert_ bersyarat (`ON CONFLICT DO NOTHING`) untuk baris baru. Hal ini merupakan standar _Best Practice_ karena menghemat beban _Write I/O_ pada database serta mempertahankan riwayat waktu asli pencapaian tugas (`created_at`) milik sang anak. Dengan ini endpoint terpisah untuk _uncheck_ tidak diperlukan lagi.

### 28. [BARU] Hapus development-report (salah input)

- **Method & Path**: `DELETE /development-reports/{child_development_report_id}`
- **Authorization**: sesuai `child_id` milik user yang login
- **Path Params**: `child_development_report_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: hard delete (`child_development_reports` tidak punya `deleted_at`). `assessment_kpsp_answers` & `development_report_recommendations` terkait ikut terhapus lewat `ON DELETE CASCADE`. Berbeda dari endpoint 24 (edit jawaban) — ini untuk kasus asesmen dibuat keliru sama sekali (mis. salah pilih anak).

---

## 5. Modul: Child Nutrition

> Catatan umum: `child_nutrition_reports` tidak punya kolom tanggal eksplisit — "hari ini" = record terbaru per `child_id` (`ORDER BY created_at DESC LIMIT 1`), dengan asumsi sistem hanya generate maksimal 1 record per anak per hari.

### 29. Get nutrition report + menu hari ini + shopping list (POV mother)

- **Method & Path**: `GET /children/{child_id}/nutrition/today`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `date` (Wajib, format `YYYY-MM-DD`)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
{
    "id": "...",
    "report_date": "2026-08-14",
    "can_regenerate": true,
    "calories": 0,
    "target_calories": 0,
    "protein": 0,
    "target_protein": 0,
    "fat": 0,
    "target_fat": 0,
    "carbohydrate": 0,
    "target_carbohydrate": 0,
    "status": "Normal",
    "menu": {
        "id": "...",
        "created_at": "...",
        "recipes": [
            {
                "id": "...",
                "name": "...",
                "meal_time": "breakfast",
                "meal_texture": "...",
                "calories": 0,
                "protein": 0
            }
        ]
    },
    "shopping_list": [
        {
            "name": "Dada ayam",
            "ingredient_id": "DA003",
            "unit": "300g"
        }
    ]
}
```

| Field                                                         | Type    | Keterangan                                                                                                                                                                          |
| ------------------------------------------------------------- | ------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| id                                                            | UUID    | `child_nutrition_reports.id`                                                                                                                                                        |
| report_date                                                   | string  | Tanggal berlakunya laporan gizi (YYYY-MM-DD)                                                                                                                                        |
| can_regenerate                                                | boolean | Indikator apakah menu hari ini bisa di-regenerate oleh AI (bernilai `false` jika resep hari ini adalah hasil 'reuse' dari hari sebelumnya)                                          |
| calories/protein/fat/carbohydrate                             | integer | aktual                                                                                                                                                                              |
| target_calories/target_protein/target_fat/target_carbohydrate | integer | target                                                                                                                                                                              |
| status                                                        | string  | Ditentukan dari jumlah makronutrisi (kalori, protein, lemak, karbo) yang nilainya >= 90% dari targetnya. (4 = "normal", 3 = "kurang optimal", 2 = "Berisiko", <=1 = "sangat buruk") |
| menu.id                                                       | UUID    | Identik dengan ID laporan (ini hanya properti penampung list recipes)                                                                                                               |
| menu.recipes[]                                                | array   | Daftar resep/makanan yang ada di dalam menu hari ini                                                                                                                                |
| shopping_list[]                                               | array   | Daftar bahan belanja dari menu (bahan utama / priority 1). Bahan yang sama (berdasarkan `ingredient_id`) akan digabung dan jumlah takarannya (unit) dijumlahkan secara otomatis     |

- **Response Error**:

| Status | Kasus                                               |
| ------ | --------------------------------------------------- |
| 400    | Parameter `date` tidak dikirim atau formatnya salah |
| 403    | `child_id` bukan milik user yang login              |
| 404    | belum ada `child_nutrition_reports` untuk hari ini  |

### 30. Reuse Recipe (Tukar Menu dengan Hari Lalu)

- **Method & Path**: `POST /children/{child_id}/nutrition/reuse-recipe`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{
    "source_recipe_id": "uuid-dari-resep-yang-ingin-dipakai",
    "report_date": "2026-08-14"
}
```

- **Response Body (200)**:

```json
{
    "message": "Recipe reused successfully"
}
```

- **Response Error**:

| Status | Kasus                                                                   |
| ------ | ----------------------------------------------------------------------- |
| 400    | format body salah / `child_id` tidak valid / format `report_date` salah |
| 403    | bukan milik user yang login                                             |
| 404    | laporan gizi tidak ditemukan / resep sumber tidak ditemukan             |

### 31. Update status selesai recipe

- **Method & Path**: `PATCH /recipes/{recipe_id}/complete`
- **Path Params**: `recipe_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{
    "portions_consumed": 0.33,
    "report_date": "2026-08-14"
}
```

| Field             | Type   | Wajib | Keterangan                          |
| ----------------- | ------ | ----- | ----------------------------------- |
| portions_consumed | float  | Ya    | Porsi makan (0, 0.33, 0.66, 1)      |
| report_date       | string | Ya    | Tanggal laporan format "YYYY-MM-DD" |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                             |
| ------ | ------------------------------------------------- |
| 403    | `recipe_id` bukan milik anak dari user yang login |
| 404    | `recipe_id` tidak ditemukan                       |

### 32. [BARU] Update bookmark recipe

- **Method & Path**: `PATCH /recipes/{recipe_id}/bookmark`
- **Path Params**: `recipe_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{ "is_bookmarked": true }
```

| Field         | Type    | Wajib | Keterangan |
| ------------- | ------- | ----- | ---------- |
| is_bookmarked | boolean | Ya    |            |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                             |
| ------ | ------------------------------------------------- |
| 403    | `recipe_id` bukan milik anak dari user yang login |
| 404    | `recipe_id` tidak ditemukan                       |

- **Catatan**: v3 menggabungkan status selesai & bookmark di satu endpoint (17) — dipisah karena dua aksi ini independen di UI (mis. tombol "selesai masak" vs ikon bookmark, dipicu di momen yang berbeda).

### 33. Get detail recipe

- **Method & Path**: `GET /recipes/{recipe_id}`
- **Path Params**: `recipe_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**:

```json
{
    "id": "0454ed87-b1ca-4a91-afda-a876015dba7b",
    "name": "Bubur Kentang Telur ayam ras",
    "meal_time": "breakfast",
    "meal_texture": "potongan kecil",
    "cooking_time": "20 menit",
    "description": "Bubur lembut berbahan dasar kentang dan telur ayam, cocok untuk bayi 6-12 bulan.",
    "calories": 269,
    "protein": 6,
    "is_bookmarked": true,
    "main_ingredients": [
        {
            "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
            "recipe_id": "0454ed87-b1ca-4a91-afda-a876015dba7b",
            "ingredient": {
                "ingredient_id": "AR001",
                "name": "Beras giling, mentah",
                "image_url": "https://unsplash.com/photos/a-single-carrot-on-a-white-background-eJJ2316gOak",
                "category": "Serealia",
                "price": "murah"
            },
            "unit": "1.4 batang (14.0g)",
            "priority": 1,
            "slot": "karbohidrat"
        }
    ],
    "recipe_spices": [
        {
            "id": "b9ee2955-a41c-4776-ac3d-6b6e92842918",
            "name": "Kunyit",
            "unit": "2 ruas (10.0g)"
        }
    ],
    "cooking_steps": [
        {
            "id": "14225ac9-25a3-4db7-95d2-229b40d1af23",
            "recipe_id": "0454ed87-b1ca-4a91-afda-a876015dba7b",
            "step_number": 1,
            "instruction": "Cuci bersih semua bahan.",
            "image_url": "https://placehold.co/600x400?text=Langkah+1",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        {
            "id": "1e06c453-7d0a-4e6f-8357-77e8addd8856",
            "recipe_id": "0454ed87-b1ca-4a91-afda-a876015dba7b",
            "step_number": 2,
            "instruction": "Haluskan/potong sesuai tekstur anak.",
            "image_url": "https://placehold.co/600x400?text=Langkah+2",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z"
        }
    ]
}
```

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan [FIX v4]**: v3 mengembalikan semua `main_ingredients` diurutkan `priority` (termasuk opsi substitusi priority 2, 3, dst). Diubah jadi hanya kembalikan `priority = 1` — opsi substitusi lain hanya relevan saat user memanggil endpoint 33 (tukar prioritas), tidak perlu ditampilkan di halaman detail resep.

### 34. Tukar prioritas main_ingredients

- **Method & Path**: `PATCH /recipes/main-ingredients/priority`
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: Array of objects untuk melakukan swap priority pada satu atau beberapa resep sekaligus.

```json
[
    {
        "recipe_id": "da1828c5-e979-4877-ac60-93c771ce160c",
        "slot": "karbohidrat",
        "priority": 2
    },
    {
        "recipe_id": "9f3c2e11-c1f0-4665-b1a8-c0b05b6b3e6c",
        "slot": "protein",
        "priority": 3
    }
]
```

| Field     | Type    | Wajib | Keterangan                                                             |
| --------- | ------- | ----- | ---------------------------------------------------------------------- |
| recipe_id | UUID    | Ya    | ID dari resep yang ingin diubah bahan utamanya                         |
| slot      | string  | Ya    | Kategori slot (mis. `karbohidrat`, `protein`, `sayur`)                 |
| priority  | integer | Ya    | Nilai priority target yang ingin ditukar menjadi priority 1 (harus >1) |

- **Response Body (200)**: `-` (body kosong / 200 OK)
- **Response Error**:

| Status | Kasus                                                                            |
| ------ | -------------------------------------------------------------------------------- |
| 400    | format body tidak valid / slot kosong / priority <= 1 / `recipe_id` invalid UUID |
| 403    | salah satu `recipe_id` bukan milik user yang login                               |
| 500    | Internal server error                                                            |

- **Catatan**:
    - Seluruh operasi swap ini bersifat **atomic** (dalam satu DB transaction). Jika salah satu gagal atau ditolak, seluruh operasi dibatalkan.
    - Bahan dengan `slot` dan `priority` yang dikirim akan di-set menjadi `priority = 1`, sedangkan bahan lainnya pada slot tersebut akan digeser posisinya secara berurutan.

### 35. Get bookmark menu

- **Method & Path**: `GET /children/{child_id}/recipes/bookmarked`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari `recipes` dengan `is_bookmarked = true`, di-join dari `child_nutrition_reports` -> `child_nutrition_recipes` -> `recipes` untuk memastikan milik `child_id` ybs

```json
[
    {
        "id": "...",
        "name": "...",
        "meal_time": "breakfast",
        "meal_texture": "...",
        "calories": 0,
        "protein": 0
    }
]
```

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

### 36. [BARU] Get riwayat nutrition report per bulan

- **Method & Path**: `GET /children/{child_id}/nutrition-reports`
- **Authorization**: sesuai `child_id`
- **Path Params**: `child_id` (UUID)
- **Query Params**:

| Field | Type    | Wajib | Keterangan |
| ----- | ------- | ----- | ---------- |
| month | integer | Ya    | 1-12       |
| year  | integer | Ya    | `YYYY`     |

- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: list `child_nutrition_reports` dalam bulan tsb, diurutkan `created_at` naik:

```json
[
    {
        "id": "...",
        "created_at": "...",
        "calories": 800,
        "protein": 20,
        "fat": 15,
        "carbohydrate": 100,
        "meal_times": ["breakfast", "lunch"]
    }
]
```

| Field                             | Type          | Keterangan                                                                                                                                                       |
| --------------------------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| id                                | UUID          | `child_nutrition_reports.id`                                                                                                                                     |
| created_at                        | datetime      |                                                                                                                                                                  |
| calories/protein/fat/carbohydrate | integer       | aktual — **[FIX v4]** field `target_*` dihapus dari response ini                                                                                                 |
| meal_times                        | array\<enum\> | **[FIX v4]** `recipes.meal_time` dari resep yang `portions_consumed > 0` pada hari itu, join `child_nutrition_reports` -> `child_nutrition_recipes` -> `recipes` |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 400    | `month` bukan 1-12          |
| 403    | bukan milik user yang login |

- **Catatan**: celah dari v2 — sebelumnya hanya ada "hari ini" (#15), tidak ada cara melihat tren gizi mingguan/bulanan (mis. untuk grafik protein vs target selama sebulan).

### 37. [BARU] Get list ingredient untuk daily shop (POV Mother & Caregiver)

- **Method & Path**: `GET /menu/daily-shop`
- **Authorization**: `user_id` diambil dari token (bisa mother atau caregiver)
- **Path Params**: `-` (tidak ada)
- **Query Params**: `date` (Wajib, format `YYYY-MM-DD`)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari object bahan utama (priority=1) dan bahan penggantinya (priority>=2) yang dikelompokkan per resep dan slot, digabung dari seluruh anak yang direlasikan ke user.
  _(Jika tidak ada menu hari ini, akan return array kosong `[]`)_

```json
[
    {
        "name": "Oat instan",
        "unit": "500g",
        "priority": 1,
        "slot": "karbohidrat",
        "recipe_id": "da1828c5-e979-4877-ac60-93c771ce160c",
        "child_name": "yasa",
        "substitutes": [
            {
                "name": "havermout",
                "unit": "500g",
                "priority": 2
            },
            {
                "name": "gandum",
                "unit": "300g",
                "priority": 3
            }
        ]
    }
]
```

- **Response Error**:

| Status | Kasus                               |
| ------ | ----------------------------------- |
| 401    | Token tidak ada/tidak valid         |
| 403    | User bukan mother ataupun caregiver |
| 500    | Internal server error               |

- **Catatan**: Melakukan query JOIN kompleks dari `children` sampai `ingredients`. Menggunakan pola _grouping_ di layer Go untuk menyusun list `substitutes`.

### 38. [BARU] Get menu hari ini + shopping list (POV caregiver)

- **Method & Path**: `GET /children/{child_id}/nutrition-reports/today/shopping`
- **Authorization**: sesuai `child_id` (mother/caregiver)
- **Path Params**: `child_id` (UUID)
- **Query Params**: `date` (Wajib, format `YYYY-MM-DD`)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Menu hari ini beserta daftar belanja, tanpa data kalori aktual/target `nutrition_report`.
  _(Jika tidak ada menu hari ini, akan return `menu` null/empty dan `shopping_list: []`)_

```json
{
    "menu": {
        "id": "...",
        "created_at": "...",
        "recipes": [
            {
                "id": "...",
                "name": "...",
                "meal_time": "breakfast",
                "meal_texture": "...",
                "calories": 0,
                "protein": 0
            }
        ]
    },
    "shopping_list": [
        {
            "name": "Oat instan",
            "ingredient_id": "OA005",
            "unit": "3.2 sdm (32g)"
        }
    ]
}
```

- **Response Error**:

| Status | Kasus                              |
| ------ | ---------------------------------- |
| 400    | Format UUID invalid                |
| 401    | Token tidak ada/tidak valid        |
| 403    | User tidak punya akses ke anak ini |
| 500    | Internal server error              |

### 39. [BARU] Generate menu dari AI (Food Engine)

- **Method & Path**: `POST /menu/generate`
- **Authorization**: `user_id` diambil dari token (atau sesuai kebutuhan arsitektur)
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{
    "report_date": "2026-08-14"
}
```

- **Response Body (201)**:

```json
{
    "status": "created",
    "message": "Menu berhasil di-generate"
}
```

- **Response Error**:

| Status | Kasus                                                               |
| ------ | ------------------------------------------------------------------- |
| 400    | format body tidak valid / `report_date` tidak ada atau salah format |
| 502    | Food Engine gagal memproses / timeout AI                            |
| 500    | Gagal menyimpan ke database                                         |

- **Catatan**: Endpoint ini adalah sebuah _Command_ yang memicu server untuk mem-_fetch_ profil anak dari DB, merakit payload `all_child.json`, mengirimnya ke AI (Food Engine), dan menyimpan hasilnya ke database (strategi _Replace_ menu hari ini). Jika ada anak yang menu hari ininya berstatus `can_regenerate=false` (resep hasil reuse), backend otomatis akan men-skip anak tersebut secara aman dan melanjutkan generate menu untuk anak lainnya. Setelah menerima respons sukses (201) dari endpoint ini, _client_ (Flutter) bertugas memanggil ulang **Endpoint 29** (`GET /menu/today`) untuk mengambil data terbaru dan me-render UI.

### 40. [BARU] Get report menu by ID

- **Method & Path**: `GET /children/{child_id}/nutrition-reports/{report_id}`
- **Authorization**: `user_id` dari token (Mother/Caregiver) yang memiliki akses ke anak
- **Path Params**:
    - `child_id` (UUID)
    - `report_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: (Sama persis dengan Endpoint 30)

```json
{
    "id": "...",
    "created_at": "...",
    "recipes": [
        {
            "id": "...",
            "name": "Nasi Tim Ikan Mas Wortel",
            "meal_time": "lunch",
            "meal_texture": "potongan kecil",
            "calories": 320,
            "protein": 18,
            "portions_consumed": 1
        }
    ]
}
```

- **Response Error**:

| Status | Kasus                              |
| ------ | ---------------------------------- |
| 400    | Format UUID invalid                |
| 401    | Token tidak ada/tidak valid        |
| 403    | User tidak punya akses ke anak ini |
| 404    | `report_id` tidak ditemukan        |
| 500    | Internal server error              |

---

## 6. Modul: Photos & Contacts

### 41. Get list teman (contacts)

- **Method & Path**: `GET /mother-profiles/contacts`
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX] (200)**: list contacts, di-join ke `mother_profiles` -> `users` untuk menampilkan nama/foto:

```json
[
    {
        "contact_id": "...",
        "related_mother_profile_id": "...",
        "full_name": "...",
        "photo_url": "..."
    }
]
```

- **Response Error**:

| Status | Kasus                                           |
| ------ | ----------------------------------------------- |
| 403    | `mother_profile_id` bukan milik user yang login |

### 42. Hapus teman

- **Method & Path**: `DELETE /contacts/{contact_id}`
- **Path Params**: `contact_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                    |
| ------ | ---------------------------------------- |
| 403    | `contact_id` bukan milik user yang login |
| 404    | tidak ditemukan                          |

- **Catatan**: hard delete (tabel `contacts` tidak punya `deleted_at`). Jika implementasi endpoint 49 membuat baris dua arah (lihat catatan di 62), pertimbangkan apakah delete ini perlu ikut menghapus baris pasangannya juga (baris di sisi teman) — perlu didiskusikan terpisah.

### 43. Get semua foto anak milik mother

- **Method & Path**: `GET /mother-profiles/child-photos`
- **Path Params**: `-` (tidak ada)
- **Query Params [FIX v4: dihapus, kembali seperti v2]**: `-` (tidak ada — selalu kembalikan foto semua anak milik mother ini)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari `child_photos`

```json
[
    {
        "id": "...",
        "child_id": "...",
        "photo_url": "...",
        "is_review_required": false
    }
]
```

- **Response Error**:

| Status | Kasus                                           |
| ------ | ----------------------------------------------- |
| 403    | `mother_profile_id` bukan milik user yang login |

### 44. Get foto anak dari suatu contact

- **Method & Path**: `GET /contacts/{contact_id}/child-photos`
- **Path Params**: `contact_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari `child_photos` via `photo_shares`

```json
[
    {
        "id": "...",
        "photo_url": "..."
    }
]
```

- **Response Error**:

| Status | Kasus                                    |
| ------ | ---------------------------------------- |
| 403    | `contact_id` bukan milik user yang login |

### 45. Get semua foto (gabungan 39+40)

- **Method & Path**: `GET /mother-profiles/child-photos/all`
- **Path Params**: `-` (tidak ada)
- **Query Params [FIX v4: dihapus]**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: Array dari `child_photos` — hanya yang `is_review_required = false` (foto caregiver yang masih pending review tidak ikut tampil di gabungan ini)

```json
[
    {
        "id": "...",
        "photo_url": "..."
    }
]
```

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

### 46. Get detail foto

- **Method & Path**: `GET /child-photos/{child_photo_id}`
- **Path Params**: `child_photo_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: Detail `child_photos` — field `visibility`, `is_review_required`, `created_at` dihapus dari response (tidak dibutuhkan di halaman detail foto)

```json
{
    "id": "...",
    "child_id": "...",
    "photo_url": "...",
    "caption": "..."
}
```

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

### 47. Tambah foto (POV mother)

- **Method & Path**: `POST /children/{child_id}/photos`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX]**:

```json
{
    "url": "https://example.com",
    "caption": "razky breakfast telur pagi ini",
    "visibility": "only",
    "list_visibility": ["<contact_id_1>", "<contact_id_2>", "<contact_id_3>"],
    "is_review_required": false
}
```

| Field              | Type          | Wajib                          | Keterangan                                          |
| ------------------ | ------------- | ------------------------------ | --------------------------------------------------- |
| url                | string        | Ya                             | -> `child_photos.photo_url`                         |
| caption            | string        | Ya **[FIX v4]**                |                                                     |
| visibility         | enum          | Ya                             | `all` / `private` / `only`                          |
| list_visibility    | array\<UUID\> | Wajib jika `visibility="only"` | tiap item `contact_id` -> insert row `photo_shares` |
| is_review_required | boolean       | Tidak (default `false`)        |                                                     |

- **Response Body (201)**: `{ "id": "<child_photo_id>" }`
- **Response Error**:

| Status | Kasus                                                                    |
| ------ | ------------------------------------------------------------------------ |
| 422    | `caption` kosong, atau `visibility="only"` tapi `list_visibility` kosong |
| 403    | `child_id` bukan milik user yang login                                   |

- **Catatan**: `visibility` diisi enum (`all`/`private`/`only`). `list_visibility` (array `contact_id`) hanya dipakai & wajib diisi ketika `visibility = "only"` — masing-masing di-insert sebagai row di `photo_shares`. `captions` -> `caption` (sesuai nama kolom, singular).

### 48. Tambah foto (POV caregiver)

- **Method & Path**: `POST /children/{child_id}/photos`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{ "url": "https://example.com", "is_review_required": true }
```

| Field              | Type    | Wajib            | Keterangan                                           |
| ------------------ | ------- | ---------------- | ---------------------------------------------------- |
| url                | string  | Ya               |                                                      |
| is_review_required | boolean | Ya, harus `true` | validasi server: request dari caregiver wajib `true` |

- **Response Body (201)**: `{ "id": "<child_photo_id>" }`
- **Response Error**:

| Status | Kasus                                                                           |
| ------ | ------------------------------------------------------------------------------- |
| 403    | `child_id` bukan target `caregiver_engagements` aktif dari caregiver yang login |
| 422    | `is_review_required` dikirim `false`                                            |

- **Catatan**: `is_review_required` wajib `true` untuk upload dari caregiver; `visibility` pakai default kolom (`private`).

### 49. [BARU] Tambah contact/teman

- **Method & Path**: `POST /mother-profiles/contacts`
- **Authorization**: `mother_profile_id` didapat dari token
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{ "related_mother_profile_id": "<mother_profile_id_teman>" }
```

| Field                     | Type | Wajib | Keterangan                                                                                      |
| ------------------------- | ---- | ----- | ----------------------------------------------------------------------------------------------- |
| related_mother_profile_id | UUID | Ya    | didapat dari hasil scan QR/deep-link teman (mekanisme sama seperti endpoint 52 untuk caregiver) |

- **Response Body (201)**: `{ "id": "<contact_id>" }`
- **Response Error**:

| Status | Kasus                                                                            |
| ------ | -------------------------------------------------------------------------------- |
| 404    | `related_mother_profile_id` tidak ditemukan                                      |
| 409    | contact sudah ada sebelumnya untuk pasangan ini                                  |
| 422    | `related_mother_profile_id` == `mother_profile_id` (tidak bisa add diri sendiri) |

- **Catatan [ASUMSI — perlu konfirmasi]**: celah terbesar di v2 — sama sekali tidak ada endpoint untuk membuat `contacts`, padahal GET (37) dan DELETE (38) sudah ada. Skema `contacts` tidak punya tabel invite/QR-token terpisah, jadi diasumsikan mekanismenya mirip endpoint 52 (share/scan kode yang meng-encode `mother_profile_id`). Diasumsikan juga relasi pertemanan bersifat **mutual/dua arah** — satu panggilan endpoint ini akan insert **dua baris** (`(mother_profile_id, related_mother_profile_id)` dan sebaliknya `(related_mother_profile_id, mother_profile_id)`) supaya kedua mother saling melihat satu sama lain di list contacts masing-masing (endpoint 41). Kalau ternyata pertemanan dimaksud satu arah saja (mis. perlu approval dulu), desain ini perlu direvisi.

### 50. [BARU] Edit foto

- **Method & Path**: `PATCH /child-photos/{child_photo_id}`
- **Authorization**: mother pemilik anak di foto tsb, ATAU (khusus field `is_review_required`) mother yang sedang me-review foto dari caregiver
- **Path Params**: `child_photo_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: semua field opsional

```json
{
    "caption": "razky breakfast telur pagi ini",
    "visibility": "only",
    "list_visibility": ["<contact_id_1>"],
    "is_review_required": false
}
```

| Field              | Type          | Wajib                                  | Keterangan                                                           |
| ------------------ | ------------- | -------------------------------------- | -------------------------------------------------------------------- |
| caption            | string        | Tidak                                  |                                                                      |
| visibility         | enum          | Tidak                                  | `all`/`private`/`only`                                               |
| list_visibility    | array\<UUID\> | Wajib jika `visibility="only"` dikirim | replace-all `photo_shares` untuk foto ini                            |
| is_review_required | boolean       | Tidak                                  | dipakai mother untuk set `false` setelah approve foto dari caregiver |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                             |
| ------ | ------------------------------------------------- |
| 403    | bukan milik user yang login                       |
| 404    | tidak ditemukan                                   |
| 422    | `visibility="only"` tapi `list_visibility` kosong |

- **Catatan**: celah dari v2 — tidak ada cara mengubah caption/visibility setelah upload, atau menandai foto dari caregiver sudah di-review.

### 51. [BARU] Hapus foto

- **Method & Path**: `DELETE /child-photos/{child_photo_id}`
- **Authorization**: mother pemilik anak di foto tsb
- **Path Params**: `child_photo_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: hard delete (`child_photos` tidak punya `deleted_at`); `photo_shares` terkait ikut terhapus lewat `ON DELETE CASCADE`.

---

## 7. Modul: Caregiver

### 52. Get list akses caregiver (POV mother)

- **Method & Path**: `GET /mother-profiles/caregiver-engagements`
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**:

```json
[
    {
        "caregiver_engagement_id": "...",
        "child_id": "...",
        "child_name": "...",
        "caregiver_profile_id": "...",
        "caregiver_name": "...",
        "phone_number": "081234567899"
    }
]
```

| Field                   | Type           | Keterangan                                                                             |
| ----------------------- | -------------- | -------------------------------------------------------------------------------------- |
| caregiver_engagement_id | UUID           | `caregiver_engagements.id`                                                             |
| child_id, child_name    | UUID, string   | `child.id`, `child.full_name`                                                          |
| caregiver_profile_id    | UUID           | `caregiver_profiles.id`                                                                |
| caregiver_name          | string         | `users.full_name` (lewat `caregiver_profiles.user_id`)                                 |
| phone_number            | string \| null | **[FIX v4]** dari `users.phone_number` — nomor yang bisa dihubungi untuk caregiver ybs |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

- **Contoh query [FIX v4 — tambah `phone_number`]**:

```sql
SELECT
  c.id            AS child_id,
  c.full_name     AS child_name,
  ce.id           AS caregiver_engagement_id,
  ce.caregiver_profile_id,
  u.full_name     AS caregiver_name,
  u.phone_number  AS phone_number
FROM child c
JOIN caregiver_engagements ce ON ce.child_id = c.id
JOIN caregiver_profiles cp    ON cp.id = ce.caregiver_profile_id
JOIN users u                  ON u.id = cp.user_id
WHERE c.mother_profile_id = ?
  AND ce.deleted_at IS NULL;
```

### 53. [BARU] Get riwayat akses caregiver yang sudah dicabut (POV mother)

- **Method & Path**: `GET /mother-profiles/caregiver-engagements/revoked`
- **Authorization**: sama seperti endpoint 53
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: sama seperti endpoint 53 (termasuk `phone_number`), tapi filter `ce.deleted_at IS NOT NULL`
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

- **Contoh query**: sama seperti endpoint 53, ganti baris terakhir jadi `AND ce.deleted_at IS NOT NULL;`
- **Catatan**: celah dari v3 — endpoint 53 hanya menampilkan akses yang masih aktif, tidak ada cara melihat riwayat caregiver yang aksesnya sudah dicabut.

### 54. Hapus akses caregiver

- **Method & Path**: `DELETE /caregiver-engagements/{caregiver_engagement_id}`
- **Path Params**: `caregiver_engagement_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: soft delete — set `caregiver_engagements.deleted_at = now()`.

### 55. Get profil caregiver [FIX]

- **Method & Path**: `GET /caregiver-profiles`
- **Authorization**: `user_id` diambil dari token
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX] (200)**: satu record `caregiver_profiles`, di-join ke `users` (konsisten dengan endpoint 6 versi mother):

```json
{
    "id": "...",
    "user_id": "...",
    "full_name": "...",
    "email": "...",
    "phone_number": "...",
    "photo_url": "..."
}
```

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan [FIX v3]**: sebelumnya hanya "satu record `caregiver_profiles`" (yang tabelnya cuma punya `id`/`user_id`/timestamp, tidak berguna untuk UI profil). Sekarang eksplisit di-join ke `users`, konsisten dengan endpoint 6. Relasi `caregiver_profiles`<->`users` bersifat 1:1, jadi tidak ada makna "terbaru" — cukup 1 record per user.

## 8. Modul: Medical

### 56. Get list medical_notes aktif per anak [FIX]

- **Method & Path**: `GET /medical-notes`
- **Path Params**: `-` (tidak ada)
- **Query Params**:

| Field      | Type    | Wajib              | Keterangan                                                                                                          |
| ---------- | ------- | ------------------ | ------------------------------------------------------------------------------------------------------------------- |
| status     | enum    | Ya                 | `active` (default filter `valid_until >= now()`) atau `history`                                                     |
| child_name | string  | Tidak **[FIX v5]** | Filter berdasarkan nama anak (exact match). Jika tidak dikirim, menampilkan semua anak.                             |
| month      | integer | Tidak **[FIX v4]** | 1-12. Default: bulan berjalan. Filter berdasarkan **nomor bulan saja** (`EXTRACT(MONTH FROM valid_until) = month`). |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari `medical_notes` aktif

```json
[
    {
        "id": "...",
        "valid_until": "22-07-2026",
        "recommendation": "...",
        "created_at": "...",
        "child_name": "...",
        "prohibition_count": 0,
        "allergy_count": 0
    }
]
```

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

- **Catatan [FIX v5]**: Filter `child_name` ditambahkan untuk menyortir berdasarkan nama anak langsung. Filter `month` tetap dipertahankan. `prohibition_count`/`allergy_count` tetap dihitung dari `medical_restrictions` per `medical_note_id`.

### 57. Get riwayat medical_notes (sudah tidak aktif) [FIX]

- **Method & Path**: `GET /medical-notes`
- **Path Params**: `-` (tidak ada)
- **Query Params**:

| Field      | Type    | Wajib              | Keterangan                                                                              |
| ---------- | ------- | ------------------ | --------------------------------------------------------------------------------------- |
| status     | enum    | Ya                 | `history` (filter `valid_until < now()`)                                                |
| child_name | string  | Tidak **[FIX v5]** | Filter berdasarkan nama anak (exact match). Jika tidak dikirim, menampilkan semua anak. |
| month      | integer | Tidak **[FIX v4]** | 1-12. Default: bulan berjalan. Filter berdasarkan **nomor bulan saja**.                 |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari riwayat `medical_notes` (sama seperti endpoint 58)

```json
[
    {
        "id": "...",
        "valid_date": "22-07-2026",
        "recommendation": "...",
        "created_at": "...",
        "child_name": "...",
        "prohibition_count": 0,
        "allergy_count": 0
    }
]
```

- **Response Error**: sama seperti endpoint 56
- **Catatan**: sama seperti catatan endpoint 56.

### 58. Get rincian medical_notes

- **Method & Path [FIX]**: `GET /medical-notes/{medical_note_id}`
- **Auth**: divalidasi lewat `medical_note_id` -> pastikan `child_id` pemilik note tsb dimiliki oleh mother/caregiver yang login.
- **Path Params**: `medical_note_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Detail `medical_notes` beserta `medical_restrictions`, `daily_nutrition_targets`, dan `child_name`

```json
{
    "id": "...",
    "child_name": "...",
    "doctor_name": "...",
    "facility_name": "...",
    "recommendation": "...",
    "valid_until": "22-07-2026",
    "created_at": "...",
    "prohibitions": ["Makanan Keras", "Serat Tinggi"],
    "allergies": ["Kacang Tanah", "Susu Sapi"],
    "daily_nutrition_targets": [
        {
            "nutrient": "calorie",
            "quantity": 10
        }
    ]
}
```

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

### 59. Tambah catatan dokter

- **Method & Path [FIX v4]**: `POST /medical-notes`
- **Path Params [FIX v4]**: `-` (tidak ada — `child_id` sekarang dikirim di body, bukan di path)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: tambah `child_id`, `facility_name` jadi opsional]**:

```json
{
    "child_id": "2c8c5e18-0eca-4bf8-9ef0-831dbc3358cb",
    "doctor_name": "dr. Tirta",
    "facility_name": "RSSA",
    "recommendation": "Berikan makanan bertekstur lembut dan tinggi kalori...",
    "daily_nutrition_targets": [
        { "nutrient": "calorie", "quantity": 10 },
        { "nutrient": "protein", "quantity": 15 },
        { "nutrient": "fat", "quantity": 20 },
        { "nutrient": "carbohydrate", "quantity": 80 }
    ],
    "prohibitions": ["Makanan Keras", "Serat Tinggi"],
    "allergies": ["Kacang Tanah", "Susu Sapi"],
    "valid_until": "22-07-2026"
}
```

| Field                                     | Type            | Wajib              | Keterangan                                                                |
| ----------------------------------------- | --------------- | ------------------ | ------------------------------------------------------------------------- |
| child_id                                  | UUID            | Ya **[FIX v4]**    | anak yang dituju catatan ini                                              |
| doctor_name                               | string          | Ya                 | maks 150 karakter, teks bebas (bukan FK)                                  |
| facility_name                             | string          | Tidak **[FIX v4]** | maks 150 karakter                                                         |
| recommendation                            | string          | Ya                 |                                                                           |
| daily_nutrition_targets[].medical_note_id | string (uuid)   | Ya                 | ID rujukan ke _medical note_                                              |
| daily_nutrition_targets[].nutrient        | string (enum)   | Ya                 | Hanya diperbolehkan: `"calories"`, `"protein"`, `"fat"`, `"carbohydrate"` |
| daily_nutrition_targets[].quantity        | float           | Ya                 | Jumlah target nutrisi per hari (misal 150.5)                              |
| prohibitions                              | array\<string\> | Tidak              | -> row `medical_restrictions` dengan `type='prohibition'`                 |
| allergies                                 | array\<string\> | Tidak              | -> row `medical_restrictions` dengan `type='allergy'`                     |
| valid_until                               | date            | Ya                 | `DD-MM-YYYY`                                                              |

- **Response Body (201)**: `{ "id": "<medical_note_id>" }`
- **Response Error**:

| Status | Kasus                                                         |
| ------ | ------------------------------------------------------------- |
| 400    | `valid_date` bukan tanggal valid                              |
| 422    | `child_id`/`doctor_name`/`recommendation`/`valid_date` kosong |
| 403    | `child_id` bukan milik user yang login                        |
| 404    | `child_id` tidak ditemukan                                    |

- **Catatan**: `prohibitions[]` -> row `medical_restrictions` dengan `type = 'prohibition'`; `allergies[]` -> row `medical_restrictions` dengan `type = 'allergy'`.

### 60. [BARU] Edit catatan dokter

- **Method & Path**: `PATCH /medical-notes/{medical_note_id}`
- **Authorization**: sama seperti endpoint 60
- **Path Params**: `medical_note_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: dijabarkan eksplisit]**: semua field opsional, kirim hanya yang mau diubah

```json
{
    "doctor_name": "dr. Tirta",
    "facility_name": "RSSA",
    "recommendation": "Berikan makanan bertekstur lembut dan tinggi kalori...",
    "daily_nutrition_targets": [{ "nutrient": "calorie", "quantity": 10 }],
    "prohibitions": ["Makanan Keras"],
    "allergies": ["Kacang Tanah"],
    "valid_until": "22-07-2026"
}
```

| Field                   | Type            | Wajib | Keterangan                                                        |
| ----------------------- | --------------- | ----- | ----------------------------------------------------------------- |
| doctor_name             | string          | Tidak | maks 150 karakter                                                 |
| facility_name           | string          | Tidak | maks 150 karakter                                                 |
| recommendation          | string          | Tidak |                                                                   |
| daily_nutrition_targets | array\<object\> | Tidak | replace-all — kirim seluruh daftar baru, bukan hanya yang berubah |
| prohibitions            | array\<string\> | Tidak | replace-all                                                       |
| allergies               | array\<string\> | Tidak | replace-all                                                       |
| valid_until             | date            | Tidak | `DD-MM-YYYY`                                                      |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |
| 422    | field yang diisi kosong     |

- **Catatan**: `child_id` tidak bisa diubah lewat endpoint ini (catatan medis tidak bisa dipindah ke anak lain — kalau salah anak, hapus lalu buat ulang lewat endpoint 61). Untuk `daily_nutrition_targets`/`prohibitions`/`allergies`, strategi update replace-all (sama seperti pola di endpoint 8). Celah dari v2 — sebelumnya hanya bisa create (33), tidak bisa edit kalau ada salah input (mis. salah ketik `valid_date`).

### 61. [BARU] Hapus catatan dokter

- **Method & Path**: `DELETE /medical-notes/{medical_note_id}`
- **Authorization**: sama seperti endpoint 60
- **Path Params**: `medical_note_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: hard delete (`medical_notes` tidak punya `deleted_at`). `medical_restrictions` & `daily_nutrition_targets` terkait ikut terhapus lewat `ON DELETE CASCADE`.

---

## 9. Modul: Notification

### 62. Get seluruh notifikasi [FIX]

- **Method & Path [FIX]**: `GET /notifications`
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Array dari `notifications`

```json
[
    {
        "id": "...",
        "title": "...",
        "message": "...",
        "notification_type": "meal_reminder",
        "created_at": "..."
    }
]
```

### 63. Get notifikasi terbaru

- **Method & Path**: `GET /notifications/latest`
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: Objek `notifications` tunggal terbaru

```json
{
    "id": "...",
    "title": "...",
    "message": "...",
    "notification_type": "meal_reminder",
    "created_at": "..."
}
```

- **Response Error**:

| Status | Kasus                               |
| ------ | ----------------------------------- |
| 404    | Tidak ada notifikasi yang ditemukan |

- **Catatan [FIX v3]**: **Keterbatasan skema** — tabel `notifications` tidak punya kolom `is_read`, sehingga fitur "tandai sudah dibaca"/badge unread **tidak bisa diimplementasikan** dengan skema saat ini. Kalau dibutuhkan, perlu penambahan kolom `is_read BOOLEAN NOT NULL DEFAULT false` di migration terpisah (di luar cakupan dokumen ini).

## 10. Modul: Dashboard (gabungan)

### 64. Get seluruh anak + tinggi, berat, skor KPSP, protein

- **Method & Path**: `GET /dashboard/children`
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list, tiap item:

```json
[
    {
        "id": "...",
        "full_name": "...",
        "photo_url": "...",
        "age": "2 tahun 3 bulan",
        "height_cm": 89,
        "weight_kg": 12.4,
        "kpsp_score": 8,
        "kpsp_answers_count": 10,
        "protein": 30,
        "target_protein": 35,
        "streak": 5,
        "status_growth": "Normal",
        "status_development": "Sesuai",
        "status_nutrition": "Normal",
        "status": "Tumbuh Optimal"
    }
]
```

| Field                   | Type            | Keterangan                                                  |
| ----------------------- | --------------- | ----------------------------------------------------------- |
| id, full_name           | UUID, string    | `children`                                                  |
| birth_date              | date            | **[FIX v4]** `child.birth_date`                             |
| height_cm, weight_kg    | number \| null  | dari `child_growth_reports` **terbaru** per `child_id`      |
| kpsp_score              | integer \| null | dari `child_development_reports` **terbaru** per `child_id` |
| protein, target_protein | integer \| null | dari `child_nutrition_reports` **terbaru** per `child_id`   |
| streak                  | integer         | **[FIX v4]** `child.streak_days`                            |
| status                  | string          | Status Gabungan. Lihat aturan penentuan di bawah.           |

**Aturan Penentuan Status Gabungan (Tumbuh, Kembang, Nutrisi):**
Berdasarkan parameter warna dari ketiga status (Hijau = Normal/Sesuai/Baik, Kuning = Berisiko/Meragukan/Kurang Optimal, Merah = Sangat Buruk/Penyimpangan):

- Jika terdapat 3 indikator Merah ➔ **"Perlu Konsultasi"**
- Jika terdapat 2 indikator Merah ➔ **"Perlu Pendampingan"**
- Jika terdapat 1 indikator Merah ➔ **"Perlu Perhatian"**
- Jika tidak ada Merah, dan terdapat 3 indikator Kuning ➔ **"Perlu Perhatian"**
- Jika tidak ada Merah, dan terdapat 1 atau 2 indikator Kuning ➔ **"Perkembangan Baik"**
- Jika semuanya (3 indikator) bernilai Hijau ➔ **"Tumbuh Optimal"**
- Jika salah satu/lebih dari ketiga status tersebut belum memiliki data ("Tidak Diketahui" / "Tidak Tersedia") ➔ **"Data Belum Lengkap"**

- **Response Error**:

| Status | Kasus                                           |
| ------ | ----------------------------------------------- |
| 403    | `mother_profile_id` bukan milik user yang login |

- **Catatan [FIX]**: `height_cm`/`weight_kg` diambil dari `child_growth_reports` **terbaru**, `kpsp_score` dari `child_development_reports` **terbaru**, `protein` (+`target_protein`) dari `child_nutrition_reports` **terbaru** — masing-masing per `child_id` (dokumentasi lama tidak menyebut "terbaru" secara eksplisit padahal relasinya 1:N).

---

## 11. Modul: Checkin

### 65. Generate Checkin QR Token (POV Ibu)

- **Method & Path [BARU]**: `POST /checkin/generate`
- **Authorization**: `Bearer <token>` (Ibu)
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{
    "childId": "123e4567-e89b-12d3-a456-426614174000"
}
```

- **Response Body (200)**:

```json
{
    "token": "4a2b8c9d1e3f...",
    "expiresAt": "2026-08-04T12:05:00Z"
}
```

### 66. Validate Checkin QR (POV Caregiver)

- **Method & Path [BARU]**: `POST /checkin/validate`
- **Authorization**: `Bearer <token>` (Caregiver)
- **Path Params**: `-` (tidak ada)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{
    "token": "4a2b8c9d1e3f..."
}
```

- **Response Body (200)**:

```json
{
    "valid": true,
    "engagementId": "...",
    "checkedInAt": "2026-08-04T12:02:00Z"
}
```

- **Response Error**:

| Status | Kasus                                 |
| ------ | ------------------------------------- |
| 404    | `token not found` / `child not found` |
| 403    | `token expired` (token kedaluwarsa)   |
| 403    | User tidak punya profile caregiver    |

> **Catatan**: Jika `active engagement already exists`, API tetap akan mengembalikan status **200 OK** (karena secara bisnis ini tetap dianggap sebagai check-in yang sukses) namun tanpa ID _engagement_ baru.

## Ringkasan Semua Endpoint (1-66)

| #   | Modul                    | Method & Path                                                              | Ringkasan                                                            |
| --- | ------------------------ | -------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| 1   | Auth & User Profile      | `PATCH /users`                                                             | Update profil user                                                   |
| 2   | Auth & User Profile      | `DELETE /users`                                                            | Hapus user (soft delete)                                             |
| 3   | Auth & User Profile      | `GET /users`                                                               | [BARU] Get profil user sendiri                                       |
| 4   | Auth & User Profile      | `POST /mother-profiles`                                                    | [BARU] Daftar sebagai mother (pilih role)                            |
| 5   | Auth & User Profile      | `POST /caregiver-profiles`                                                 | [BARU] Daftar sebagai caregiver (pilih role)                         |
| 6   | Auth & User Profile      | `GET /mother-profiles`                                                     | [BARU] Get profil mother sendiri                                     |
| 7   | Child Profile            | `POST /children`                                                           | Tambah profil anak                                                   |
| 8   | Child Profile            | `PATCH /children/{child_id}`                                               | Update profil anak                                                   |
| 9   | Child Profile            | `DELETE /children/{child_id}`                                              | [BARU] Hapus profil anak (soft delete)                               |
| 10  | Child Profile            | `GET /children/{child_id}`                                                 | [BARU] Get detail profil anak                                        |
| 11  | Child Profile            | `GET /children`                                                            | [BARU] Get list anak milik mother (ringan)                           |
| 12  | Child Profile            | `GET /children/{child_id}/simple`                                          | [BARU] Get profil anak (ringan)                                      |
| 13  | Child Growth             | `GET /children/{child_id}/growth-reports/latest`                           | Get child_growth_reports terbaru                                     |
| 14  | Child Growth             | `GET /children/{child_id}/growth-analyses`                                 | Get child_growth_analyses (filter age-range & analysis_type) [FIX]   |
| 15  | Child Growth             | `GET /children/{child_id}/growth-reports`                                  | Get riwayat child_growth_reports                                     |
| 16  | Child Growth             | `POST /children/{child_id}/growth-reports`                                 | Tambah child_growth_reports                                          |
| 17  | Child Growth             | `PATCH /growth-reports/{child_growth_report_id}`                           | [BARU] Edit child_growth_reports                                     |
| 18  | Child Growth             | `DELETE /growth-reports/{child_growth_report_id}`                          | [BARU] Hapus child_growth_reports                                    |
| 19  | Child Development (KPSP) | `GET /children/{child_id}/development-reports/latest`                      | Get child_development_reports terbaru + 4 domain                     |
| 20  | Child Development (KPSP) | `GET /children/{child_id}/development-reports`                             | Get riwayat asesmen KPSP                                             |
| 21  | Child Development (KPSP) | `GET /development-reports/{child_development_report_id}`                   | Get detail hasil asesmen KPSP                                        |
| 22  | Child Development (KPSP) | `GET /assessment-kpsp-questions`                                           | Get pertanyaan assessment_kpsp_questions                             |
| 23  | Child Development (KPSP) | `POST /children/{child_id}/development-reports`                            | Simpan hasil asesmen KPSP (create)                                   |
| 24  | Child Development (KPSP) | `PATCH /development-reports/{child_development_report_id}`                 | Update hasil asesmen KPSP terbaru                                    |
| 25  | Child Development (KPSP) | `GET /checklist-milestone-tasks`                                           | Get checklist_milestone_tasks                                        |
| 26  | Child Development (KPSP) | `GET /children/{child_id}/development-reports/{report_id}/recommendations` | Get rekomendasi aksi KPSP berdasarkan hasil asesmen                  |
| 27  | Child Development (KPSP) | `PATCH /children/{child_id}/checklist-milestone-progress`                  | Tandai checklist milestone (Sync / UPSERT)                           |
| 28  | Child Development (KPSP) | `DELETE /development-reports/{child_development_report_id}`                | [BARU] Hapus development-report (salah input)                        |
| 29  | Child Nutrition          | `GET /children/{child_id}/nutrition/today`                                 | Get nutrition report + menu hari ini + shopping list (POV mother)    |
| 30  | Child Nutrition          | `POST /children/{child_id}/nutrition/reuse-recipe`                         | Reuse Recipe (Tukar Menu dengan Hari Lalu)                           |
| 31  | Child Nutrition          | `PATCH /recipes/{recipe_id}/complete`                                      | Update status selesai recipe                                         |
| 32  | Child Nutrition          | `PATCH /recipes/{recipe_id}/bookmark`                                      | [BARU] Update bookmark recipe                                        |
| 33  | Child Nutrition          | `GET /recipes/{recipe_id}`                                                 | Get detail recipe                                                    |
| 34  | Child Nutrition          | `PATCH /recipes/{recipe_id}/main-ingredients/{main_ingredient_id}`         | Tukar prioritas main_ingredients                                     |
| 35  | Child Nutrition          | `GET /children/{child_id}/recipes/bookmarked`                              | Get bookmark menu                                                    |
| 36  | Child Nutrition          | `GET /children/{child_id}/nutrition-reports`                               | [BARU] Get riwayat nutrition report per bulan                        |
| 37  | Child Nutrition          | `GET /menu/daily-shop`                                                     | [BARU] Get list ingredient untuk daily shop (POV Mother & Caregiver) |
| 38  | Child Nutrition          | `GET /children/{child_id}/nutrition-reports/today/shopping`                | [BARU] Get menu hari ini + shopping list (POV caregiver)             |
| 39  | Child Nutrition          | `POST /menu/generate`                                                      | [BARU] Generate menu dari AI (Food Engine)                           |
| 40  | Child Nutrition          | `GET /children/{child_id}/nutrition-reports/{report_id}`                   | [BARU] Get report menu by ID                                         |
| 41  | Photos & Contacts        | `GET /mother-profiles/contacts`                                            | Get list teman (contacts)                                            |
| 42  | Photos & Contacts        | `DELETE /contacts/{contact_id}`                                            | Hapus teman                                                          |
| 43  | Photos & Contacts        | `GET /mother-profiles/child-photos`                                        | Get semua foto anak milik mother                                     |
| 44  | Photos & Contacts        | `GET /contacts/{contact_id}/child-photos`                                  | Get foto anak dari suatu contact                                     |
| 45  | Photos & Contacts        | `GET /mother-profiles/child-photos/all`                                    | Get semua foto (gabungan 39+40)                                      |
| 46  | Photos & Contacts        | `GET /child-photos/{child_photo_id}`                                       | Get detail foto                                                      |
| 47  | Photos & Contacts        | `POST /children/{child_id}/photos`                                         | Tambah foto (POV mother)                                             |
| 48  | Photos & Contacts        | `POST /children/{child_id}/photos`                                         | Tambah foto (POV caregiver)                                          |
| 49  | Photos & Contacts        | `POST /mother-profiles/contacts`                                           | [BARU] Tambah contact/teman                                          |
| 50  | Photos & Contacts        | `PATCH /child-photos/{child_photo_id}`                                     | [BARU] Edit foto                                                     |
| 51  | Photos & Contacts        | `DELETE /child-photos/{child_photo_id}`                                    | [BARU] Hapus foto                                                    |
| 52  | Caregiver                | `GET /mother-profiles/caregiver-engagements`                               | Get list akses caregiver (POV mother)                                |
| 53  | Caregiver                | `GET /mother-profiles/caregiver-engagements/revoked`                       | [BARU] Get riwayat akses caregiver yang sudah dicabut (POV mother)   |
| 54  | Caregiver                | `DELETE /caregiver-engagements/{caregiver_engagement_id}`                  | Hapus akses caregiver                                                |
| 55  | Caregiver                | `GET /caregiver-profiles`                                                  | Get profil caregiver [FIX]                                           |
| 56  | Medical                  | `GET /mother-profiles/medical-notes`                                       | Get list medical_notes aktif per anak [FIX]                          |
| 57  | Medical                  | `GET /mother-profiles/medical-notes`                                       | Get riwayat medical_notes (sudah tidak aktif) [FIX]                  |
| 58  | Medical                  | `GET /medical-notes/{medical_note_id}`                                     | Get rincian medical_notes                                            |
| 59  | Medical                  | `POST /medical-notes`                                                      | Tambah catatan dokter                                                |
| 60  | Medical                  | `PATCH /medical-notes/{medical_note_id}`                                   | [BARU] Edit catatan dokter                                           |
| 61  | Medical                  | `DELETE /medical-notes/{medical_note_id}`                                  | [BARU] Hapus catatan dokter                                          |
| 62  | Notification             | `GET /notifications`                                                       | Get seluruh notifikasi [FIX]                                         |
| 63  | Notification             | `GET /notifications/latest`                                                | Get notifikasi terbaru                                               |
| 64  | Dashboard (gabungan)     | `GET /dashboard/children`                                                  | Get seluruh anak + tinggi, berat, skor KPSP, protein                 |
| 65  | Checkin                  | `POST /checkin/generate`                                                   | Generate Checkin QR Token (POV Ibu)                                  |
| 66  | Checkin                  | `POST /checkin/validate`                                                   | Validate Checkin QR (POV Caregiver)                                  |

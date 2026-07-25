# NusaGizi — Dokumentasi Endpoint API (v4)

Dokumen ini adalah revisi dari `NusaGizi_Endpoint_Documentation_v3.md`, berdasarkan review langsung per endpoint. Ringkasan perubahan v3 -> v4 ada di tabel paling bawah dokumen. Riwayat sebelumnya:

1. **v3**: standar dokumentasi per endpoint dibuat konsisten penuh (_Path Params_, _Query Params_, _Request Body_, _Response Body_, _Response Error_ di semua endpoint), endpoint baru nomor 47–67 ditambahkan untuk menutup celah CRUD.
2. **Tidak ada pagination** — riwayat/list difilter per-bulan (`month`[+`year`]) atau age-range, bukan `page`/`limit`.
3. **Endpoint autentikasi (login/register/refresh/logout) sengaja tidak dibuat** — sepenuhnya ditangani Auth0 di luar API ini.

Nomor endpoint asli dari v1/v2 dipertahankan agar mudah dibandingkan; endpoint yang isinya diubah ditandai **[FIX]**, endpoint baru ditandai **[BARU]**.

---

## 0. Konvensi Umum

- **Auth**: ditangani provider eksternal (Auth0). Setiap request terautentikasi membawa `user_id` dari token. Role (mother/caregiver) ditentukan dari ada-tidaknya row terkait di `mother_profiles`/`caregiver_profiles` untuk `user_id` tsb.
- **Di luar cakupan dokumen ini**: login, register, refresh token, logout (semua ditangani Auth0). Pembuatan row `users` diasumsikan terjadi otomatis (mis. via Auth0 Action/webhook) saat signup pertama kali — dokumen ini mulai dari endpoint pertama yang dipanggil aplikasi **setelah** row `users` ada, yaitu pemilihan role (lihat endpoint 47/48 di bawah).
- **Format ID**: semua ID adalah UUID (string), termasuk di request body — bukan integer.
- **Format tanggal** (`date`, contoh kolom `birth_date`, `measured_at`, `valid_date`): `DD-MM-YYYY` (contoh: `25-07-2026`).
- **Format timestamp** (`datetime`, contoh kolom `created_at`, `updated_at`): ISO 8601 UTC (contoh: `2026-07-25T09:30:00Z`).
- **Penamaan field**: English snake_case, mengikuti nama kolom database (contoh: `full_name`, `birth_date`).
- **Soft delete**: hanya untuk tabel yang punya kolom `deleted_at` (`users`, `child`, `caregiver_engagements`) — endpoint DELETE terkait melakukan `UPDATE ... SET deleted_at = now()`. Tabel lain (`child_growth_reports`, `child_development_reports`, `child_photos`, `medical_notes`, `contacts`, `notification`, `ingredient_shopping_items`, dst) **tidak** punya `deleted_at`, sehingga endpoint DELETE-nya melakukan hard delete (baris benar-benar dihapus).
- **Response POST**: kecuali dinyatakan lain, endpoint POST yang membuat resource baru mengembalikan minimal `{ "id": "<uuid>" }` (status `201 Created`), supaya client tidak perlu re-fetch.
- **Response PATCH/DELETE tanpa body**: status `200 OK` dengan body kosong (`-`), kecuali dinyatakan lain.
- **Tidak ada pagination**: endpoint list/riwayat difilter per periode (`month` + `year`) atau age-range, bukan `page`/`limit`. Jika suatu saat volume data jadi masalah, pagination bisa ditambahkan sebagai perubahan terpisah.
- **Base path**: contoh path di dokumen ini pakai konvensi resource-oriented (`/children/{child_id}/...`), sesuaikan dengan prefix API Anda (mis. `/api/v1/...`).

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

> Alur singkat: user signup/login via Auth0 → row `users` tersedia (di luar cakupan) → aplikasi tanya "Anda ibu atau caregiver?" → panggil endpoint 47 atau 48 untuk membuat profile sesuai role.

### 23. Update profil user

- **Method & Path**: `PATCH /users/{user_id}`
- **Authorization**: hanya bisa update `user_id` milik sendiri (harus sama dengan `user_id` di token)
- **Path Params**:

| Field   | Type | Keterangan            |
| ------- | ---- | --------------------- |
| user_id | UUID | ID user yang diupdate |

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

| Status | Kasus                                  |
| ------ | -------------------------------------- |
| 400    | `email` bukan format email valid       |
| 403    | `user_id` di path ≠ `user_id` di token |
| 404    | `user_id` tidak ditemukan              |
| 409    | `email` sudah dipakai user lain        |

### 25. Hapus user (soft delete)

- **Method & Path**: `DELETE /users/{user_id}`
- **Authorization**: hanya bisa delete `user_id` milik sendiri
- **Path Params**:

| Field   | Type | Keterangan           |
| ------- | ---- | -------------------- |
| user_id | UUID | ID user yang dihapus |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                     |
| ------ | ----------------------------------------- |
| 403    | `user_id` di path ≠ `user_id` di token    |
| 404    | `user_id` tidak ditemukan / sudah dihapus |

- **Catatan**: soft delete — set `users.deleted_at = now()`. Perlu didiskusikan terpisah (di luar scope dokumen ini) apakah `child`/data terkait ikut di-cascade-soft-delete atau tetap ada.

### 45. [BARU] Get profil user sendiri

- **Method & Path**: `GET /users/{user_id}` (atau `GET /users/me`)
- **Authorization**: hanya bisa GET `user_id` milik sendiri
- **Path Params**:

| Field   | Type | Keterangan                   |
| ------- | ---- | ---------------------------- |
| user_id | UUID | ID user (boleh diganti `me`) |

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

| Status | Kasus                     |
| ------ | ------------------------- |
| 403    | `user_id` di path ≠ token |
| 404    | `user_id` tidak ditemukan |

- **Catatan**: ditambahkan agar form edit di endpoint 23 punya data awal untuk di-load. Field `has_mother_profile`/`has_caregiver_profile` ditambahkan di v3 supaya app tahu kapan harus mengarahkan user ke endpoint 47/48.

### 47. [BARU] Daftar sebagai mother (pilih role)

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

### 48. [BARU] Daftar sebagai caregiver (pilih role)

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

### 49. [BARU] Get profil mother sendiri

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}` (atau `/mother-profiles/me`)
- **Authorization**: hanya bisa GET `mother_profile_id` milik sendiri
- **Path Params**:

| Field             | Type | Keterangan         |
| ----------------- | ---- | ------------------ |
| mother_profile_id | UUID | boleh diganti `me` |

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

- **Catatan**: paralel dengan endpoint 29 (`GET /caregiver-profiles/{id}`) yang sebelumnya sudah ada tapi belum ada versi mother-nya — celah ini ditutup di v3. Endpoint 29 juga di-[FIX] responsnya (lihat Modul 7) supaya konsisten (ikut join ke `users`).

---

## 2. Modul: Child Profile

### 22. Tambah profil anak

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
    "food_frequency": 4,
    "food_goal": "Menyesuaikan kondisi kesehatan",
    "notes": "Saya mau menu yang murah"
}
```

| Field                                          | Type            | Wajib | Keterangan                                                            |
| ---------------------------------------------- | --------------- | ----- | --------------------------------------------------------------------- |
| full_name                                      | string          | Ya    | Maks 150 karakter                                                     |
| birth_date                                     | date            | Ya    | `DD-MM-YYYY`                                                          |
| gender                                         | enum            | Ya    | `male` / `female`                                                     |
| photo_url                                      | string          | Tidak |                                                                       |
| allergies.food / .medicine / .animal / .others | array\<string\> | Tidak | tiap item -> row `child_allergy_profile` dengan `category` sesuai key |
| chronic_diseases                               | array\<string\> | Tidak | -> row `child_chronic_disease_profile.disease_name`                   |
| diets                                          | array\<string\> | Tidak | -> row `child_diet_profile.diet_name`                                 |
| favorite_foods                                 | array\<string\> | Tidak | -> row `favorite_food_profile.food_name`                              |
| favorite_textures                              | array\<string\> | Tidak | -> row `favorite_texture_profile.texture_name`                        |
| food_frequency                                 | integer         | Ya    | -> `child.food_frequency_profile`                                     |
| food_goal                                      | string          | Ya    | -> `child.food_goal_profile`                                          |
| notes                                          | string          | Tidak | -> `child.notes_profile`                                              |

- **Response Body (201)**: `{ "id": "<child_id>" }`
- **Response Error**:

| Status | Kasus                                                                 |
| ------ | --------------------------------------------------------------------- |
| 400    | `birth_date` bukan tanggal valid, `gender` bukan enum valid           |
| 422    | `full_name`/`birth_date`/`gender`/`food_frequency`/`food_goal` kosong |

- **Catatan**: Field diganti dari Indonesia (`name`, `date_of_birth`, `alergi`, `kondisi_kronis`, dst) ke English snake_case sesuai konvensi.

### 24. Update profil anak

- **Method & Path**: `PATCH /children/{child_id}`
- **Auth [FIX]**: divalidasi lewat `child_id` — pastikan `child.mother_profile_id` terhubung ke `mother_profiles.user_id` yang login. (Dokumentasi lama menyebut "berdasarkan user_id", tapi tabel `child` tidak punya kolom `user_id` sama sekali.)
- **Path Params**:

| Field    | Type | Keterangan |
| -------- | ---- | ---------- |
| child_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: sama seperti endpoint 22, semua field opsional (hanya kirim yang mau diubah)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                      |
| ------ | ------------------------------------------ |
| 403    | `child_id` bukan milik mother yang login   |
| 404    | `child_id` tidak ditemukan / sudah dihapus |

- **Catatan**: untuk field turunan (`allergies`, `chronic_diseases`, `diets`, `favorite_foods`, `favorite_textures`), strategi update disarankan replace-all (hapus semua row lama untuk child tsb, insert ulang dari array baru) supaya tidak perlu diff per-item.

### 46. [BARU] Hapus profil anak (soft delete)

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

### 50. [BARU] Get detail profil anak

- **Method & Path**: `GET /children/{child_id}`
- **Authorization**: pastikan `child.mother_profile_id` milik `user_id` yang login (mother) atau ada `caregiver_engagements` aktif untuk `child_id` ini (caregiver)
- **Path Params**:

| Field    | Type | Keterangan |
| -------- | ---- | ---------- |
| child_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: bentuk sama persis dengan request body endpoint 22, plus `id` — supaya bisa langsung dipakai sebagai initial state form edit (endpoint 24):

```json
{
    "id": "<child_id>",
    "full_name": "ibor",
    "birth_date": "22-07-2026",
    "gender": "male",
    "photo_url": "https://docs.google.com/",
    "upload_streak_days": 12,
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
    "food_frequency": 4,
    "food_goal": "Menyesuaikan kondisi kesehatan",
    "notes": "Saya mau menu yang murah"
}
```

| Field              | Type    | Keterangan                           |
| ------------------ | ------- | ------------------------------------ |
| upload_streak_days | integer | `child.upload_streak_days`           |
| (field lain)       | —       | sama seperti tabel field endpoint 22 |

- **Response Error**:

| Status | Kasus                                                              |
| ------ | ------------------------------------------------------------------ |
| 403    | bukan mother pemilik / caregiver yang tidak punya engagement aktif |
| 404    | tidak ditemukan / sudah dihapus                                    |

- **Catatan**: celah dari v2 — sebelumnya tidak ada endpoint untuk mengambil profil anak secara lengkap (yang ada hanya ringkasan dashboard di endpoint 1, yang fieldnya jauh lebih sedikit).

### 51. [BARU] Get list anak milik mother (ringan)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/children`
- **Authorization**: `mother_profile_id` sesuai `user_id` yang login
- **Path Params**:

| Field             | Type | Keterangan |
| ----------------- | ---- | ---------- |
| mother_profile_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
[
    {
        "id": "...",
        "full_name": "ibor",
        "birth_date": "22-07-2026",
        "gender": "male",
        "photo_url": "..."
    }
]
```

| Field      | Type           | Keterangan |
| ---------- | -------------- | ---------- |
| id         | UUID           |            |
| full_name  | string         |            |
| birth_date | date           |            |
| gender     | enum           |            |
| photo_url  | string \| null |            |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

- **Catatan**: versi ringan dari endpoint 1 (dashboard) — dipakai untuk kebutuhan seperti dropdown/switcher pilih anak, tanpa perlu join ke growth/nutrition/KPSP yang lebih berat. Bentuk response-nya sengaja disamakan dengan endpoint 34 (list anak versi caregiver) untuk konsistensi antar-role.

---

## 3. Modul: Child Growth

### 3. Get child_growth_reports terbaru

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
    "head_circumference_cm": 47
}
```

| Field                 | Type           | Keterangan |
| --------------------- | -------------- | ---------- |
| id                    | UUID           |            |
| measured_at           | date           |            |
| weight_kg             | number \| null |            |
| height_cm             | number \| null |            |
| head_circumference_cm | number \| null |            |

- **Response Error**:

| Status | Kasus                                                                       |
| ------ | --------------------------------------------------------------------------- |
| 403    | `child_id` bukan milik user yang login                                      |
| 404    | `child_id` tidak ditemukan, atau belum ada record growth report sama sekali |

### 4. Get child_growth_analyses (filter age-range & analysis_type) [FIX]

- **Method & Path**: `GET /children/{child_id}/growth-analyses`
- **Authorization**: `child_growth_analyses` yang `child_growth_report_id`-nya milik `child_id` ybs
- **Path Params**: `child_id` (UUID)
- **Query Params [FIX v4: dijadikan wajib]**:

| Field         | Type | Wajib | Keterangan                                                                                                                     |
| ------------- | ---- | ----- | ------------------------------------------------------------------------------------------------------------------------------ |
| analysis_type | enum | Ya    | salah satu dari `weight_for_age`, `height_for_age`, `weight_for_height`, `bmi_for_age`, `head_circumference_for_age`           |
| age_range     | enum | Ya    | `0-2`, `0-6`, atau `0-60` (satuan **bulan usia anak**, dihitung dari `child.birth_date` ke `child_growth_reports.measured_at`) |

- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**:

```json
[{ "z_score": -1.2, "age_months": 5 }]
```

| Field      | Type           | Keterangan                                                                               |
| ---------- | -------------- | ---------------------------------------------------------------------------------------- |
| z_score    | number \| null |                                                                                          |
| age_months | integer        | dihitung: umur anak (bulan) saat `measured_at` — dipakai FE sebagai titik sumbu-X grafik |

- **Response Error**:

| Status | Kasus                                                                        |
| ------ | ---------------------------------------------------------------------------- |
| 400    | `analysis_type`/`age_range` tidak dikirim, atau bukan salah satu nilai valid |
| 403    | `child_id` bukan milik user yang login                                       |
| 404    | `child_id` tidak ditemukan                                                   |

- **Catatan [FIX v4]**: `analysis_type` & `age_range` dijadikan wajib (v3 masih opsional, tidak ada gunanya untuk kasus pakai grafik yang selalu butuh 1 tipe + 1 rentang umur spesifik). Response dipangkas jadi cuma `z_score`+`age_months` — cukup untuk mapping langsung ke titik grafik di FE. Field `id` & `analysis_type` dibuang karena sudah pasti seragam (difilter di query); `measured_at` diganti `age_months` karena itu yang langsung dipakai sebagai sumbu-X grafik pertumbuhan WHO, bukan tanggal kalender mentah. Query tetap join ke `child_growth_reports.measured_at` untuk hitung `age_months` dan untuk filter `age_range`.

### 5. Get riwayat child_growth_reports

- **Method & Path**: `GET /children/{child_id}/growth-reports`
- **Authorization**: sesuai `child_id`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada — dikembalikan semua, diurutkan `measured_at` terbaru dulu)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list `child_growth_reports`, field sama seperti endpoint 3
- **Response Error**:

| Status | Kasus                                  |
| ------ | -------------------------------------- |
| 403    | `child_id` bukan milik user yang login |
| 404    | `child_id` tidak ditemukan             |

### 6. Tambah child_growth_reports

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

### 52. [BARU] Edit child_growth_reports

- **Method & Path**: `PATCH /growth-reports/{child_growth_report_id}`
- **Authorization**: `child_growth_report_id` -> `child_id` -> `mother_profile_id` harus milik user yang login
- **Path Params**:

| Field                  | Type | Keterangan |
| ---------------------- | ---- | ---------- |
| child_growth_report_id | UUID |            |

- **Query Params**: `-` (tidak ada)
- **Request Body**: sama seperti endpoint 6, semua field opsional
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |
| 422    | nilai <= 0                  |

- **Catatan**: dipakai untuk koreksi salah input (mis. salah ketik berat badan), tanpa perlu hapus lalu buat ulang (yang akan mengacak-acak grafik pertumbuhan karena `child_growth_analyses` terhubung via `child_growth_report_id`).

### 53. [BARU] Hapus child_growth_reports

- **Method & Path**: `DELETE /growth-reports/{child_growth_report_id}`
- **Authorization**: sama seperti endpoint 52
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

- **Catatan**: hard delete (`child_growth_reports` tidak punya `deleted_at`) — `child_growth_analyses` terkait ikut terhapus otomatis lewat `ON DELETE CASCADE`.

---

## 4. Modul: Child Development (KPSP)

**Algoritma bersama untuk endpoint 7, 9, 11, 12** (dikonfirmasi user):

```
KPSP_PERIODS = [3, 6, 9, 12, 15, 18, 21, 24, 30, 36, 42, 48, 54, 60]  // bulan
age_months   = umur anak sekarang, dihitung dari birth_date

month_target = periode TERBESAR di KPSP_PERIODS yang <= age_months
               jika age_months < 3 -> month_target = 3

kpsp_score   = jumlah jawaban "true" dari 10 pertanyaan (skala 0-10)

next_period  = periode TERKECIL di KPSP_PERIODS yang > age_months
next_check_date =
    jika next_period ada  -> birth_date + next_period bulan   (format DD-MM-YYYY)
    jika tidak ada (age_months >= 60, periode terakhir sudah dikerjakan)
                          -> literal string "done"
```

Contoh: usia 7 bulan -> `month_target=6`, `next_check_date` = usia 9 bulan (2 bulan lagi). Usia 2 bulan -> `month_target=3`, `next_check_date` = usia 3 bulan (bulan depan).

### 7. Get child_development_reports terbaru + 4 domain

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
    "created_at": "...",
    "domains": [
        {
            "nerve_name": "Gross motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "nerve_name": "Fine motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "nerve_name": "Speech and language",
            "total_question": 3,
            "true_answer": 2
        },
        { "nerve_name": "Socialization", "total_question": 3, "true_answer": 2 }
    ]
}
```

| Field                    | Type           | Keterangan                                                                                 |
| ------------------------ | -------------- | ------------------------------------------------------------------------------------------ |
| id                       | UUID           |                                                                                            |
| kpsp_score               | integer        | 0-10                                                                                       |
| next_check_date          | date \| "done" |                                                                                            |
| created_at               | datetime       |                                                                                            |
| domains[].nerve_name     | enum           | 4 aspek: `Gross motor skills`, `Fine motor skills`, `Speech and language`, `Socialization` |
| domains[].total_question | integer        | jumlah pertanyaan di domain ini untuk `month_target` ybs                                   |
| domains[].true_answer    | integer        | jumlah jawaban `true` di domain ini                                                        |

- **Response Error**:

| Status | Kasus                                             |
| ------ | ------------------------------------------------- |
| 403    | `child_id` bukan milik user yang login            |
| 404    | belum ada `child_development_reports` sama sekali |

- **Catatan [FIX v4]**: v3 mengembalikan raw list jawaban per domain lalu FE yang menghitung `total_question`/`true_answer`. Diubah jadi agregasi di query (`GROUP BY nerve_name`, `COUNT(*)` untuk `total_question`, `COUNT(*) FILTER (WHERE answer = true)` untuk `true_answer`) — lebih baik dari sisi performa karena payload jauh lebih kecil dan logic agregasi cukup ditulis sekali di backend, tidak diulang di tiap platform client (mobile/web).

### 8. Get riwayat asesmen KPSP

- **Method & Path**: `GET /children/{child_id}/development-reports`
- **Authorization**: sesuai `child_id`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: list `child_development_reports` (`id, kpsp_score, month_target, created_at`)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | `child_id` tidak ditemukan  |

### 9. Get detail hasil asesmen KPSP

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
    "recommended_actions": [
        { "id": "...", "title": "...", "action_text": "..." }
    ],
    "domains": [
        {
            "nerve_name": "Gross motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "nerve_name": "Fine motor skills",
            "total_question": 3,
            "true_answer": 2
        },
        {
            "nerve_name": "Speech and language",
            "total_question": 3,
            "true_answer": 2
        },
        { "nerve_name": "Socialization", "total_question": 3, "true_answer": 2 }
    ]
}
```

| Field                                           | Type                 | Keterangan                                                                                    |
| ----------------------------------------------- | -------------------- | --------------------------------------------------------------------------------------------- |
| id                                              | UUID                 |                                                                                               |
| kpsp_score                                      | integer              |                                                                                               |
| month_target                                    | integer              | **[FIX v4]** menggantikan `next_check_date` (dihapus — tidak relevan di halaman detail hasil) |
| recommended_actions[].id/title/action_text      | UUID/string/string   |                                                                                               |
| domains[].nerve_name/total_question/true_answer | enum/integer/integer | **[FIX v4]** disamakan dengan format agregat endpoint 7                                       |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: `recommended_actions` yang tampil hanya yang terhubung ke pertanyaan yang dijawab `false` (indikasi keterlambatan), diambil lewat `assessment_attempt_recomendations`.

### 10. Get pertanyaan assessment_kpsp_questions

- **Method & Path**: `GET /assessment-kpsp-questions`
- **Path Params**: `-` (tidak ada)
- **Query Params [FIX]**:

| Field        | Type    | Wajib | Keterangan                                                                                                                                   |
| ------------ | ------- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| month_target | integer | Ya    | salah satu dari `3,6,9,12,15,18,21,24,30,36,42,48,54,60` — sebelumnya tidak ada filter di dokumentasi lama, padahal soal berbeda per periode |

- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: list `assessment_kpsp_questions` (`id, month_target, nerve_name, question_text, description`) — field `month_target` ditambahkan ke tiap item response
- **Response Error**:

| Status | Kasus                                       |
| ------ | ------------------------------------------- |
| 400    | `month_target` bukan salah satu nilai valid |

### 11. Simpan hasil asesmen KPSP (create)

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

- **Response Body (201)**: sama seperti endpoint 9 — field `id` di response **adalah** `child_development_report_id` yang baru dibuat
- **Response Error**:

| Status | Kasus                                                                                                                     |
| ------ | ------------------------------------------------------------------------------------------------------------------------- |
| 422    | `list_answer` bukan 10 item, atau ada `assessment_kpsp_question_id` yang tidak match periode `month_target` anak saat ini |
| 403    | `child_id` bukan milik user yang login                                                                                    |

- **Catatan [FIX v4]**: `month_target` dikirim oleh client di body request, lalu divalidasi terhadap `child.birth_date` memakai algoritma di atas modul ini; `kpsp_score` & `next_check_date` tetap dihitung 100% di server. Alur pembuatan ID: server men-generate UUID baru untuk `child_development_report_id` **sebelum** proses insert, lalu UUID yang sama itu dipakai sebagai foreign key `child_development_report_id` di setiap row `assessment_kpsp_answers` yang diinsert dalam satu transaksi — jadi client tidak perlu (dan tidak boleh) generate/kirim ID ini. Setelah semua 10 soal dijawab, sistem juga insert ke `assessment_attempt_recomendations` untuk tiap pertanyaan yang dijawab `false`.

### 12. Update hasil asesmen KPSP terbaru

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

- **Response Body (200)**: sama seperti endpoint 9
- **Response Error**:

| Status | Kasus                                         |
| ------ | --------------------------------------------- |
| 403    | bukan milik user yang login                   |
| 404    | `child_development_report_id` tidak ditemukan |

- **Catatan**: `kpsp_score`, `next_check_date`, dan `assessment_attempt_recomendations` dihitung ulang. Dipanggil saat mother mengulangi asesmen.

### 13. Get checklist_milestone_tasks

- **Method & Path**: `GET /checklist-milestone-tasks`
- **Path Params**: `-` (tidak ada)
- **Query Params**:

| Field        | Type    | Wajib                     | Keterangan                           |
| ------------ | ------- | ------------------------- | ------------------------------------ |
| month_target | integer | Ya                        | salah satu dari 14 nilai valid       |
| child_id     | UUID    | Ya **[FIX: ditambahkan]** | dipakai untuk left join `is_checked` |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list `checklist_milestone_tasks` (`id, nerve_name, task_description, is_checked`)
- **Response Error**:

| Status | Kasus                                  |
| ------ | -------------------------------------- |
| 400    | `month_target` invalid                 |
| 403    | `child_id` bukan milik user yang login |

- **Catatan**: `is_checked` didapat dari left join ke `checklist_milestone_progress` dengan `child_id` yang dikirim, supaya FE tahu item mana yang sudah dicentang tanpa request terpisah.

### 14. Tandai checklist milestone (replace-all)

- **Method & Path**: `PATCH /children/{child_id}/checklist-milestone-progress`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: semantik jadi replace-all, mencakup fungsi endpoint 55 lama]**:

```json
{
    "checklist_milestone_task_ids": [
        "b6f1e2d0-...-uuid-1",
        "8a2c3f10-...-uuid-2"
    ]
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

- **Catatan [FIX v4]**: semantik diubah dari "tambah centang saja" (insert-only) menjadi **replace-all** — client selalu mengirim seluruh state checklist yang tercentang untuk `child_id` ini di tiap panggilan; server menghapus row `checklist_milestone_progress` yang tidak ada di list baru, dan insert yang belum ada. Dengan ini endpoint terpisah untuk uncheck tidak diperlukan lagi — uncheck satu item cukup dengan mengirim ulang list tanpa task tsb.

### 54. [BARU] Hapus development-report (salah input)

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

- **Catatan**: hard delete (`child_development_reports` tidak punya `deleted_at`). `assessment_kpsp_answers` & `assessment_attempt_recomendations` terkait ikut terhapus lewat `ON DELETE CASCADE`. Berbeda dari endpoint 12 (edit jawaban) — ini untuk kasus asesmen dibuat keliru sama sekali (mis. salah pilih anak).

---

## 5. Modul: Child Nutrition

> Catatan umum: `child_nutrition_reports` tidak punya kolom tanggal eksplisit — "hari ini" = record terbaru per `child_id` (`ORDER BY created_at DESC LIMIT 1`), dengan asumsi sistem hanya generate maksimal 1 record per anak per hari (dikonfirmasi user).

### 15. Get nutrition report + shopping + menu hari ini (POV mother)

- **Method & Path**: `GET /children/{child_id}/nutrition/today`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
{
    "id": "...",
    "calories": 0,
    "target_calories": 0,
    "protein": 0,
    "target_protein": 0,
    "fat": 0,
    "target_fat": 0,
    "carbohydrate": 0,
    "target_carbohydrate": 0,
    "shopping": {
        "id": "...",
        "is_completed": false,
        "items": [
            {
                "id": "...",
                "ingredient_name": "Bawang Merah",
                "quantity": 0,
                "unit": "g"
            }
        ]
    },
    "menu": {
        "id": "...",
        "recipes": [
            {
                "id": "...",
                "name": "...",
                "meal_time": "sarapan",
                "meal_texture": "...",
                "is_alergen": false,
                "calories": 0,
                "protein": 0
            }
        ]
    }
}
```

| Field                                                         | Type            | Keterangan                                                                                                                                                                    |
| ------------------------------------------------------------- | --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| id                                                            | UUID            | `child_nutrition_reports.id`                                                                                                                                                  |
| calories/protein/fat/carbohydrate                             | integer         | aktual                                                                                                                                                                        |
| target_calories/target_protein/target_fat/target_carbohydrate | integer         | target                                                                                                                                                                        |
| shopping.id                                                   | UUID            | `daily_shoppings.id`                                                                                                                                                          |
| shopping.is_completed                                         | boolean         |                                                                                                                                                                               |
| shopping.items[]                                              | array\<object\> | `id, ingredient_name, quantity, unit` — dari `ingredient_shopping_items`, di-join ke `ingredients` untuk ambil `name` **[FIX v4]** (tidak lagi expose `ingredient_id` mentah) |
| menu.id                                                       | UUID            | `daily_menus.id`                                                                                                                                                              |
| menu.recipes[]                                                | array\<object\> | dari `recipes`                                                                                                                                                                |

- **Response Error**:

| Status | Kasus                                              |
| ------ | -------------------------------------------------- |
| 403    | `child_id` bukan milik user yang login             |
| 404    | belum ada `child_nutrition_reports` untuk hari ini |

- **Catatan [FIX]**: field `allergen` di dokumentasi lama -> nama kolom sebenarnya `is_alergen`.

### 16. Get menu hari ini

- **Method & Path**: `GET /children/{child_id}/daily-menus/today`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: satu record `daily_menus` beserta list `recipes` (field sama seperti #15)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | belum ada menu hari ini     |

### 17. Update status selesai recipe [FIX v4: dipisah dari toggle bookmark]

- **Method & Path**: `PATCH /recipes/{recipe_id}/complete`
- **Path Params**: `recipe_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{ "is_completed": true }
```

| Field        | Type    | Wajib | Keterangan |
| ------------ | ------- | ----- | ---------- |
| is_completed | boolean | Ya    |            |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                             |
| ------ | ------------------------------------------------- |
| 403    | `recipe_id` bukan milik anak dari user yang login |
| 404    | `recipe_id` tidak ditemukan                       |

### 68. [BARU] Update bookmark recipe

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

### 18. Get detail recipe

- **Method & Path**: `GET /recipes/{recipe_id}`
- **Path Params**: `recipe_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: satu record `recipes` beserta `main_ingredients` (**hanya** row dengan `priority = 1` — yaitu bahan yang sedang dipakai sekarang, bukan daftar opsi substitusi) dan `cooking_steps` (diurutkan `step_number`)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan [FIX v4]**: v3 mengembalikan semua `main_ingredients` diurutkan `priority` (termasuk opsi substitusi priority 2, 3, dst). Diubah jadi hanya kembalikan `priority = 1` — opsi substitusi lain hanya relevan saat user memanggil endpoint 19 (tukar prioritas), tidak perlu ditampilkan di halaman detail resep.

### 19. Tukar prioritas main_ingredients

- **Method & Path [FIX]**: `PATCH /recipes/{recipe_id}/main-ingredients/{main_ingredient_id}`
- **Path Params**:

| Field              | Type | Keterangan                                                                                                                |
| ------------------ | ---- | ------------------------------------------------------------------------------------------------------------------------- |
| recipe_id          | UUID |                                                                                                                           |
| main_ingredient_id | UUID | **[FIX]** sebelumnya tidak ada di spesifikasi endpoint — ditambahkan karena wajib untuk tahu ingredient mana yang ditukar |

- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: dihapus]**: `-` (tidak ada — cukup panggil endpoint ini, tidak perlu kirim apa pun)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                                   |
| ------ | ------------------------------------------------------- |
| 403    | `recipe_id` bukan milik user yang login                 |
| 404    | `main_ingredient_id` tidak ditemukan di `recipe_id` tsb |

- **Catatan [FIX v4]**: `main_ingredient_id` di path otomatis di-set jadi `priority = 1`. `main_ingredients` lain di recipe yang sama otomatis digeser +1 (yang sebelumnya `priority=1` jadi `2`, yang `2` jadi `3`, dst).

### 21. Get bookmark menu

- **Method & Path**: `GET /children/{child_id}/recipes/bookmarked`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list `recipes` dengan `is_bookmarked = true`, di-join dari `child_nutrition_reports` -> `daily_menus` -> `recipes` untuk memastikan milik `child_id` ybs
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

### 56. [BARU] Get riwayat nutrition report per bulan

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
        "meal_times": ["sarapan", "makan_siang"]
    }
]
```

| Field                             | Type          | Keterangan                                                                                                                                                       |
| --------------------------------- | ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| id                                | UUID          | `child_nutrition_reports.id`                                                                                                                                     |
| created_at                        | datetime      |                                                                                                                                                                  |
| calories/protein/fat/carbohydrate | integer       | aktual — **[FIX v4]** field `target_*` dihapus dari response ini                                                                                                 |
| meal_times                        | array\<enum\> | **[FIX v4]** `recipes.meal_time` dari resep yang `is_completed = true` pada `daily_menus` hari itu, join `child_nutrition_reports` -> `daily_menus` -> `recipes` |

- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 400    | `month` bukan 1-12          |
| 403    | bukan milik user yang login |

- **Catatan**: celah dari v2 — sebelumnya hanya ada "hari ini" (#15), tidak ada cara melihat tren gizi mingguan/bulanan (mis. untuk grafik protein vs target selama sebulan).

### 57. [BARU] Tandai belanja selesai

- **Method & Path**: `PATCH /daily-shoppings/{daily_shopping_id}`
- **Authorization**: `daily_shopping_id` -> `child_nutrition_report_id` -> `child_id` harus milik user yang login
- **Path Params**: `daily_shopping_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{ "is_completed": true }
```

| Field        | Type    | Wajib | Keterangan |
| ------------ | ------- | ----- | ---------- |
| is_completed | boolean | Ya    |            |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: celah dari v2 — `is_completed` sudah muncul di response endpoint 15, tapi tidak ada endpoint untuk mengubahnya.

### 58. [BARU] Tambah item belanja manual

- **Method & Path**: `POST /daily-shoppings/{daily_shopping_id}/items`
- **Authorization**: sama seperti endpoint 57
- **Path Params**: `daily_shopping_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{ "ingredient_id": "...", "quantity": 250, "unit": "g" }
```

| Field         | Type   | Wajib | Keterangan                                            |
| ------------- | ------ | ----- | ----------------------------------------------------- |
| ingredient_id | UUID   | Ya    | harus ada di master `ingredients` (lihat endpoint 61) |
| quantity      | number | Ya    | > 0                                                   |
| unit          | string | Ya    | maks 20 karakter                                      |

- **Response Body (201)**: `{ "id": "<ingredient_shopping_item_id>" }`
- **Response Error**:

| Status | Kasus                                                    |
| ------ | -------------------------------------------------------- |
| 404    | `daily_shopping_id` atau `ingredient_id` tidak ditemukan |
| 422    | `quantity` <= 0                                          |

- **Catatan**: dipakai saat user mau menambah item belanja di luar yang otomatis di-generate dari resep (mis. tambah "tisu basah").

### 59. [BARU] Update item belanja

- **Method & Path**: `PATCH /ingredient-shopping-items/{ingredient_shopping_item_id}`
- **Authorization**: item -> `daily_shopping_id` -> ... harus milik user yang login
- **Path Params**: `ingredient_shopping_item_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: semua field opsional

```json
{ "quantity": 300, "unit": "g" }
```

| Field    | Type   | Wajib | Keterangan     |
| -------- | ------ | ----- | -------------- |
| quantity | number | Tidak | > 0 jika diisi |
| unit     | string | Tidak |                |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |
| 422    | `quantity` <= 0             |

### 60. [BARU] Hapus item belanja

- **Method & Path**: `DELETE /ingredient-shopping-items/{ingredient_shopping_item_id}`
- **Authorization**: sama seperti endpoint 59
- **Path Params**: `ingredient_shopping_item_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: hard delete (`ingredient_shopping_items` tidak punya `deleted_at`).

### 61. [BARU] Cari master ingredients (autocomplete)

- **Method & Path**: `GET /ingredients`
- **Authorization**: tidak perlu kepemilikan khusus (data master, sama untuk semua user)
- **Path Params**: `-` (tidak ada)
- **Query Params**:

| Field  | Type   | Wajib | Keterangan                                                                                            |
| ------ | ------ | ----- | ----------------------------------------------------------------------------------------------------- |
| search | string | Tidak | filter `name ILIKE '%search%'`. Jika kosong, kembalikan semua (atau batasi N teratas di implementasi) |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**:

```json
[{ "id": "...", "name": "Bawang Merah", "image_url": "https://..." }]
```

- **Response Error**: `-` (tidak ada kasus error khusus selain 400 jika `search` terlalu panjang)
- **Catatan**: celah dari v2 — endpoint 58 (tambah item belanja manual) dan `main_ingredients` butuh `ingredient_id` yang valid, tapi tidak ada cara mencarinya dari sisi client sebelum ini.

---

## 6. Modul: Photos & Contacts

### 37. Get list teman (contacts)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/contacts`
- **Path Params**: `mother_profile_id` (UUID)
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

### 38. Hapus teman

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

- **Catatan**: hard delete (tabel `contacts` tidak punya `deleted_at`). Jika implementasi endpoint 62 membuat baris dua arah (lihat catatan di 62), pertimbangkan apakah delete ini perlu ikut menghapus baris pasangannya juga (baris di sisi teman) — perlu didiskusikan terpisah.

### 39. Get semua foto anak milik mother

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/child-photos`
- **Path Params**: `mother_profile_id` (UUID)
- **Query Params [FIX v4: dihapus, kembali seperti v2]**: `-` (tidak ada — selalu kembalikan foto semua anak milik mother ini)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list `child_photos` (`id, child_id, photo_url, is_review_required`)
- **Response Error**:

| Status | Kasus                                           |
| ------ | ----------------------------------------------- |
| 403    | `mother_profile_id` bukan milik user yang login |

### 40. Get foto anak dari suatu contact

- **Method & Path**: `GET /contacts/{contact_id}/child-photos`
- **Path Params**: `contact_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list `child_photos` (`id, photo_url`), via `photo_shared_with`
- **Response Error**:

| Status | Kasus                                    |
| ------ | ---------------------------------------- |
| 403    | `contact_id` bukan milik user yang login |

### 41. Get semua foto (gabungan 39+40)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/child-photos/all`
- **Path Params**: `mother_profile_id` (UUID)
- **Query Params [FIX v4: dihapus]**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: list `child_photos` (`id, photo_url`) — hanya yang `is_review_required = false` (foto caregiver yang masih pending review tidak ikut tampil di gabungan ini)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

### 42. Get detail foto

- **Method & Path**: `GET /child-photos/{child_photo_id}`
- **Path Params**: `child_photo_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: `id, child_id, photo_url, caption` — field `visibility`, `is_review_required`, `created_at` dihapus dari response (tidak dibutuhkan di halaman detail foto)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

### 43. Tambah foto (POV mother)

- **Method & Path**: `POST /children/{child_id}/photos`
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX]**:

```json
{
    "url": "https://example.com",
    "caption": "razky sarapan telur pagi ini",
    "visibility": "only",
    "list_visibility": ["<contact_id_1>", "<contact_id_2>", "<contact_id_3>"],
    "is_review_required": false
}
```

| Field              | Type          | Wajib                          | Keterangan                                               |
| ------------------ | ------------- | ------------------------------ | -------------------------------------------------------- |
| url                | string        | Ya                             | -> `child_photos.photo_url`                              |
| caption            | string        | Ya **[FIX v4]**                |                                                          |
| visibility         | enum          | Ya                             | `all` / `private` / `only`                               |
| list_visibility    | array\<UUID\> | Wajib jika `visibility="only"` | tiap item `contact_id` -> insert row `photo_shared_with` |
| is_review_required | boolean       | Tidak (default `false`)        |                                                          |

- **Response Body (201)**: `{ "id": "<child_photo_id>" }`
- **Response Error**:

| Status | Kasus                                                                    |
| ------ | ------------------------------------------------------------------------ |
| 422    | `caption` kosong, atau `visibility="only"` tapi `list_visibility` kosong |
| 403    | `child_id` bukan milik user yang login                                   |

- **Catatan**: `visibility` diisi enum (`all`/`private`/`only`). `list_visibility` (array `contact_id`) hanya dipakai & wajib diisi ketika `visibility = "only"` — masing-masing di-insert sebagai row di `photo_shared_with`. `captions` -> `caption` (sesuai nama kolom, singular).

### 44. Tambah foto (POV caregiver)

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

### 62. [BARU] Tambah contact/teman

- **Method & Path**: `POST /mother-profiles/{mother_profile_id}/contacts`
- **Authorization**: `mother_profile_id` di path harus milik user yang login
- **Path Params**: `mother_profile_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**:

```json
{ "related_mother_profile_id": "<mother_profile_id_teman>" }
```

| Field                     | Type | Wajib | Keterangan                                                                                      |
| ------------------------- | ---- | ----- | ----------------------------------------------------------------------------------------------- |
| related_mother_profile_id | UUID | Ya    | didapat dari hasil scan QR/deep-link teman (mekanisme sama seperti endpoint 26 untuk caregiver) |

- **Response Body (201)**: `{ "id": "<contact_id>" }`
- **Response Error**:

| Status | Kasus                                                                            |
| ------ | -------------------------------------------------------------------------------- |
| 404    | `related_mother_profile_id` tidak ditemukan                                      |
| 409    | contact sudah ada sebelumnya untuk pasangan ini                                  |
| 422    | `related_mother_profile_id` == `mother_profile_id` (tidak bisa add diri sendiri) |

- **Catatan [ASUMSI — perlu konfirmasi]**: celah terbesar di v2 — sama sekali tidak ada endpoint untuk membuat `contacts`, padahal GET (37) dan DELETE (38) sudah ada. Skema `contacts` tidak punya tabel invite/QR-token terpisah, jadi diasumsikan mekanismenya mirip endpoint 26 (share/scan kode yang meng-encode `mother_profile_id`). Diasumsikan juga relasi pertemanan bersifat **mutual/dua arah** — satu panggilan endpoint ini akan insert **dua baris** (`(mother_profile_id, related_mother_profile_id)` dan sebaliknya `(related_mother_profile_id, mother_profile_id)`) supaya kedua mother saling melihat satu sama lain di list contacts masing-masing (endpoint 37). Kalau ternyata pertemanan dimaksud satu arah saja (mis. perlu approval dulu), desain ini perlu direvisi.

### 63. [BARU] Edit foto

- **Method & Path**: `PATCH /child-photos/{child_photo_id}`
- **Authorization**: mother pemilik anak di foto tsb, ATAU (khusus field `is_review_required`) mother yang sedang me-review foto dari caregiver
- **Path Params**: `child_photo_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: semua field opsional

```json
{
    "caption": "razky sarapan telur pagi ini",
    "visibility": "only",
    "list_visibility": ["<contact_id_1>"],
    "is_review_required": false
}
```

| Field              | Type          | Wajib                                  | Keterangan                                                           |
| ------------------ | ------------- | -------------------------------------- | -------------------------------------------------------------------- |
| caption            | string        | Tidak                                  |                                                                      |
| visibility         | enum          | Tidak                                  | `all`/`private`/`only`                                               |
| list_visibility    | array\<UUID\> | Wajib jika `visibility="only"` dikirim | replace-all `photo_shared_with` untuk foto ini                       |
| is_review_required | boolean       | Tidak                                  | dipakai mother untuk set `false` setelah approve foto dari caregiver |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                                             |
| ------ | ------------------------------------------------- |
| 403    | bukan milik user yang login                       |
| 404    | tidak ditemukan                                   |
| 422    | `visibility="only"` tapi `list_visibility` kosong |

- **Catatan**: celah dari v2 — tidak ada cara mengubah caption/visibility setelah upload, atau menandai foto dari caregiver sudah di-review.

### 64. [BARU] Hapus foto

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

- **Catatan**: hard delete (`child_photos` tidak punya `deleted_at`); `photo_shared_with` terkait ikut terhapus lewat `ON DELETE CASCADE`.

---

## 7. Modul: Caregiver

### 26. Tambah akses caregiver (scan QR)

- **Method & Path [FIX]**: `POST /children/{child_id}/caregiver-engagements`
- **Authorization**: `caregiver_profile_id` diambil dari `user_id` caregiver yang login; `child_id` dari path (hasil scan QR)
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (201)**: `{ "id": "<caregiver_engagement_id>" }`
- **Response Error**:

| Status | Kasus                                                                           |
| ------ | ------------------------------------------------------------------------------- |
| 404    | `child_id` tidak ditemukan                                                      |
| 409    | engagement aktif untuk pasangan `caregiver_profile_id`+`child_id` ini sudah ada |

### 27. Get list akses caregiver (POV mother)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/caregiver-engagements`
- **Path Params**: `mother_profile_id` (UUID)
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

### 28. Hapus akses caregiver

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

### 69. [BARU] Get riwayat akses caregiver yang sudah dicabut (POV mother)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/caregiver-engagements/revoked`
- **Authorization**: sama seperti endpoint 27
- **Path Params**: `mother_profile_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: sama seperti endpoint 27 (termasuk `phone_number`), tapi filter `ce.deleted_at IS NOT NULL`
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

- **Contoh query**: sama seperti endpoint 27, ganti baris terakhir jadi `AND ce.deleted_at IS NOT NULL;`
- **Catatan**: celah dari v3 — endpoint 27 hanya menampilkan akses yang masih aktif, tidak ada cara melihat riwayat caregiver yang aksesnya sudah dicabut.

### 29. Get profil caregiver [FIX]

- **Method & Path**: `GET /caregiver-profiles/{caregiver_profile_id}` (atau `/caregiver-profiles/me`)
- **Path Params**: `caregiver_profile_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX] (200)**: satu record `caregiver_profiles`, di-join ke `users` (konsisten dengan endpoint 49 versi mother):

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

- **Catatan [FIX v3]**: sebelumnya hanya "satu record `caregiver_profiles`" (yang tabelnya cuma punya `id`/`user_id`/timestamp, tidak berguna untuk UI profil). Sekarang eksplisit di-join ke `users`, konsisten dengan endpoint 49. Relasi `caregiver_profiles`<->`users` bersifat 1:1, jadi tidak ada makna "terbaru" — cukup 1 record per user.

### 34. Get list child dari caregiver_engagement

- **Method & Path**: `GET /caregiver-profiles/{caregiver_profile_id}/children`
- **Path Params**: `caregiver_profile_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: list `child` (`id, full_name, photo_url`) — `birth_date`/`gender` dihapus dari response ini
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |

- **Contoh query [FIX v4]**:

```sql
SELECT c.id, c.full_name, c.photo_url
FROM caregiver_engagements ce
JOIN child c ON c.id = ce.child_id
WHERE ce.caregiver_profile_id = ?
  AND ce.deleted_at IS NULL;
```

### 35. Get nutrition report + shopping + menu hari ini (POV caregiver)

- **Method & Path**: `GET /children/{child_id}/nutrition/today` (endpoint sama dengan #15, akses via caregiver engagement bukan mother_profile_id)
- **Path Params**: `child_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body [FIX v4] (200)**: subset dari endpoint 15 — hanya `id, shopping, menu` (field nutrisi aktual/target di top-level tidak ikut disertakan, karena POV caregiver tidak perlu lihat angka target gizi)
- **Response Error**: sama seperti endpoint 15, ditambah 403 jika tidak ada `caregiver_engagements` aktif untuk `child_id` ini

### 36. Get detail recipe (POV caregiver)

- **Method & Path**: `GET /recipes/{recipe_id}` (endpoint sama dengan #18)
- **Path Params**: `recipe_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: sama seperti #18, tanpa `calories, protein, carbohydrate, fat`
- **Response Error**: sama seperti endpoint 18

---

## 8. Modul: Medical

### 30. Get list medical_notes aktif per anak (POV mother) [FIX]

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/medical-notes`
- **Path Params**: `mother_profile_id` (UUID)
- **Query Params**:

| Field  | Type    | Wajib              | Keterangan                                                                                                                                                                                                             |
| ------ | ------- | ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| status | enum    | Ya                 | `active` (default filter `valid_date >= now()`)                                                                                                                                                                        |
| month  | integer | Tidak **[FIX v4]** | 1-12. Default: bulan berjalan. Filter berdasarkan **nomor bulan saja** (`EXTRACT(MONTH FROM valid_date) = month`) — tidak peduli tahun, jadi mother bisa lihat "semua catatan bulan Juli" dari tahun manapun sekaligus |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list, tiap item: `{ id, valid_date, recommendation, created_at, child_name, prohibition_count, allergy_count }`
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 400    | `month` bukan 1-12          |
| 403    | bukan milik user yang login |

- **Catatan [FIX v4]**: `prohibition_count`/`allergy_count` dihitung dari `medical_restrictions` per `medical_note_id`, dikelompokkan per `type`. Param `year` dihapus (v3 mewajibkan `year` bersama `month`) — kini `month` berdiri sendiri dan opsional, default ke bulan berjalan jika tidak dikirim, dan mencocokkan nomor bulan tanpa syarat tahun.

### 31. Get riwayat medical_notes (sudah tidak aktif) [FIX]

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/medical-notes`
- **Path Params**: `mother_profile_id` (UUID)
- **Query Params**:

| Field  | Type    | Wajib              | Keterangan                                                                                            |
| ------ | ------- | ------------------ | ----------------------------------------------------------------------------------------------------- |
| status | enum    | Ya                 | `history` (filter `valid_date < now()`)                                                               |
| month  | integer | Tidak **[FIX v4]** | 1-12. Default: bulan berjalan. Sama seperti endpoint 30 — filter nomor bulan saja, tidak peduli tahun |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: sama seperti #30
- **Response Error**: sama seperti #30
- **Catatan**: sama seperti catatan endpoint 30.

### 32. Get rincian medical_notes

- **Method & Path [FIX]**: `GET /medical-notes/{medical_note_id}`
- **Auth**: divalidasi lewat `medical_note_id` -> pastikan `child_id` pemilik note tsb dimiliki oleh mother/caregiver yang login.
- **Path Params**: `medical_note_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: satu record `medical_notes` beserta list `medical_restrictions`, `daily_nutrition_targets`, dan `child.full_name`
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

### 33. Tambah catatan dokter

- **Method & Path [FIX v4]**: `POST /medical-notes`
- **Path Params [FIX v4]**: `-` (tidak ada — `child_id` sekarang dikirim di body, bukan di path)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: tambah `child_id`, `facility_location` jadi opsional]**:

```json
{
    "child_id": "2c8c5e18-0eca-4bf8-9ef0-831dbc3358cb",
    "doctor_name": "dr. Tirta",
    "facility_location": "RSSA",
    "recommendation": "Berikan makanan bertekstur lembut dan tinggi kalori...",
    "daily_nutrition_targets": [
        { "nutrient_name": "calories", "quantity": 10 },
        { "nutrient_name": "protein", "quantity": 15 },
        { "nutrient_name": "fat", "quantity": 20 },
        { "nutrient_name": "carbohydrate", "quantity": 80 }
    ],
    "prohibitions": ["Makanan Keras", "Serat Tinggi"],
    "allergies": ["Kacang Tanah", "Susu Sapi"],
    "valid_date": "22-07-2026"
}
```

| Field                                   | Type            | Wajib              | Keterangan                                                |
| --------------------------------------- | --------------- | ------------------ | --------------------------------------------------------- |
| child_id                                | UUID            | Ya **[FIX v4]**    | anak yang dituju catatan ini                              |
| doctor_name                             | string          | Ya                 | maks 150 karakter, teks bebas (bukan FK)                  |
| facility_location                       | string          | Tidak **[FIX v4]** | maks 150 karakter                                         |
| recommendation                          | string          | Ya                 |                                                           |
| daily_nutrition_targets                 | array\<object\> | Tidak              | tiap item -> row `daily_nutrition_targets`                |
| daily_nutrition_targets[].nutrient_name | string          | Ya (jika item ada) |                                                           |
| daily_nutrition_targets[].quantity      | number          | Ya (jika item ada) | > 0                                                       |
| prohibitions                            | array\<string\> | Tidak              | -> row `medical_restrictions` dengan `type='prohibition'` |
| allergies                               | array\<string\> | Tidak              | -> row `medical_restrictions` dengan `type='allergy'`     |
| valid_date                              | date            | Ya                 | `DD-MM-YYYY`                                              |

- **Response Body (201)**: `{ "id": "<medical_note_id>" }`
- **Response Error**:

| Status | Kasus                                                         |
| ------ | ------------------------------------------------------------- |
| 400    | `valid_date` bukan tanggal valid                              |
| 422    | `child_id`/`doctor_name`/`recommendation`/`valid_date` kosong |
| 403    | `child_id` bukan milik user yang login                        |
| 404    | `child_id` tidak ditemukan                                    |

- **Catatan**: `prohibitions[]` -> row `medical_restrictions` dengan `type = 'prohibition'`; `allergies[]` -> row `medical_restrictions` dengan `type = 'allergy'`.

### 65. [BARU] Edit catatan dokter

- **Method & Path**: `PATCH /medical-notes/{medical_note_id}`
- **Authorization**: sama seperti endpoint 32
- **Path Params**: `medical_note_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body [FIX v4: dijabarkan eksplisit]**: semua field opsional, kirim hanya yang mau diubah

```json
{
    "doctor_name": "dr. Tirta",
    "facility_location": "RSSA",
    "recommendation": "Berikan makanan bertekstur lembut dan tinggi kalori...",
    "daily_nutrition_targets": [
        { "nutrient_name": "calories", "quantity": 10 }
    ],
    "prohibitions": ["Makanan Keras"],
    "allergies": ["Kacang Tanah"],
    "valid_date": "22-07-2026"
}
```

| Field                   | Type            | Wajib | Keterangan                                                        |
| ----------------------- | --------------- | ----- | ----------------------------------------------------------------- |
| doctor_name             | string          | Tidak | maks 150 karakter                                                 |
| facility_location       | string          | Tidak | maks 150 karakter                                                 |
| recommendation          | string          | Tidak |                                                                   |
| daily_nutrition_targets | array\<object\> | Tidak | replace-all — kirim seluruh daftar baru, bukan hanya yang berubah |
| prohibitions            | array\<string\> | Tidak | replace-all                                                       |
| allergies               | array\<string\> | Tidak | replace-all                                                       |
| valid_date              | date            | Tidak | `DD-MM-YYYY`                                                      |

- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |
| 422    | field yang diisi kosong     |

- **Catatan**: `child_id` tidak bisa diubah lewat endpoint ini (catatan medis tidak bisa dipindah ke anak lain — kalau salah anak, hapus lalu buat ulang lewat endpoint 33). Untuk `daily_nutrition_targets`/`prohibitions`/`allergies`, strategi update replace-all (sama seperti pola di endpoint 24). Celah dari v2 — sebelumnya hanya bisa create (33), tidak bisa edit kalau ada salah input (mis. salah ketik `valid_date`).

### 66. [BARU] Hapus catatan dokter

- **Method & Path**: `DELETE /medical-notes/{medical_note_id}`
- **Authorization**: sama seperti endpoint 32
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

### 2. Get seluruh notifikasi [FIX]

- **Method & Path**: `GET /users/{user_id}/notifications`
- **Path Params**: `user_id` (UUID)
- **Query Params [BARU]**:

| Field             | Type | Wajib | Keterangan                                                                                                                                                                                                                                                     |
| ----------------- | ---- | ----- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| notification_type | enum | Tidak | salah satu dari 8 nilai (`meal_reminder`, `growth_development_reminder`, `doctor_activity`, `photos_activity`, `shop_activity`, `invitation_activity`, `new_menu_reminder`, `achievement`); kalau tidak dikirim, return semua notifikasi tanpa filter kategori |

- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list `notification` (`id, title, message, notification_type, created_at`)
- **Response Error**:

| Status | Kasus                                                                |
| ------ | -------------------------------------------------------------------- |
| 400    | `notification_type` dikirim tapi bukan salah satu dari 8 nilai valid |
| 403    | `user_id` di path ≠ token                                            |

- **Catatan [FIX v3]**: filter `notification_type` ditambahkan untuk kebutuhan tab/kategori di UI notifikasi. **Keterbatasan skema** — tabel `notification` tidak punya kolom `is_read`, sehingga fitur "tandai sudah dibaca"/badge unread **tidak bisa diimplementasikan** dengan skema saat ini. Kalau dibutuhkan, perlu penambahan kolom `is_read BOOLEAN NOT NULL DEFAULT false` di migration terpisah (di luar cakupan dokumen ini).

### 67. [BARU] Hapus notifikasi

- **Method & Path**: `DELETE /notifications/{notification_id}`
- **Authorization**: `notification_id` -> `user_id` harus sama dengan token
- **Path Params**: `notification_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: `-` (body kosong)
- **Response Error**:

| Status | Kasus                       |
| ------ | --------------------------- |
| 403    | bukan milik user yang login |
| 404    | tidak ditemukan             |

- **Catatan**: hard delete (`notification` tidak punya `deleted_at`). Celah dari v2 — sebelumnya hanya bisa GET list, tidak ada cara menghapus notifikasi individual (mis. swipe-to-delete di UI).

---

## 10. Modul: Dashboard (gabungan)

### 1. Get seluruh anak + tinggi, berat, skor KPSP, protein

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/children/summary`
- **Path Params**: `mother_profile_id` (UUID)
- **Query Params**: `-` (tidak ada)
- **Request Body**: `-` (tidak ada)
- **Response Body (200)**: list, tiap item:

```json
{
    "id": "...",
    "full_name": "...",
    "birth_date": "22-07-2026",
    "height_cm": 89,
    "weight_kg": 12.4,
    "kpsp_score": 8,
    "protein": 30,
    "target_protein": 35,
    "streak": 5
}
```

| Field                   | Type            | Keterangan                                                  |
| ----------------------- | --------------- | ----------------------------------------------------------- |
| id, full_name           | UUID, string    | `child`                                                     |
| birth_date              | date            | **[FIX v4]** `child.birth_date`                             |
| height_cm, weight_kg    | number \| null  | dari `child_growth_reports` **terbaru** per `child_id`      |
| kpsp_score              | integer \| null | dari `child_development_reports` **terbaru** per `child_id` |
| protein, target_protein | integer \| null | dari `child_nutrition_reports` **terbaru** per `child_id`   |
| streak                  | integer         | **[FIX v4]** `child.upload_streak_days`                     |

- **Response Error**:

| Status | Kasus                                           |
| ------ | ----------------------------------------------- |
| 403    | `mother_profile_id` bukan milik user yang login |

- **Catatan [FIX]**: `height_cm`/`weight_kg` diambil dari `child_growth_reports` **terbaru**, `kpsp_score` dari `child_development_reports` **terbaru**, `protein` (+`target_protein`) dari `child_nutrition_reports` **terbaru** — masing-masing per `child_id` (dokumentasi lama tidak menyebut "terbaru" secara eksplisit padahal relasinya 1:N).

---

## Ringkasan Semua Endpoint Baru (nomor 47–69)

| #   | Modul             | Method & Path                                             | Ringkasan                                                |
| --- | ----------------- | --------------------------------------------------------- | -------------------------------------------------------- |
| 47  | Auth & User       | `POST /mother-profiles`                                   | Daftar sebagai mother (role selection)                   |
| 48  | Auth & User       | `POST /caregiver-profiles`                                | Daftar sebagai caregiver (role selection)                |
| 49  | Auth & User       | `GET /mother-profiles/{id}`                               | Get profil mother sendiri                                |
| 50  | Child Profile     | `GET /children/{child_id}`                                | Get detail profil anak (lengkap)                         |
| 51  | Child Profile     | `GET /mother-profiles/{id}/children`                      | Get list anak (ringan, tanpa stats)                      |
| 52  | Child Growth      | `PATCH /growth-reports/{id}`                              | Edit record pengukuran                                   |
| 53  | Child Growth      | `DELETE /growth-reports/{id}`                             | Hapus record pengukuran                                  |
| 54  | Child Development | `DELETE /development-reports/{id}`                        | Hapus record asesmen KPSP                                |
| 56  | Child Nutrition   | `GET /children/{child_id}/nutrition-reports`              | Riwayat nutrition report per bulan                       |
| 57  | Child Nutrition   | `PATCH /daily-shoppings/{id}`                             | Tandai belanja selesai                                   |
| 58  | Child Nutrition   | `POST /daily-shoppings/{id}/items`                        | Tambah item belanja manual                               |
| 59  | Child Nutrition   | `PATCH /ingredient-shopping-items/{id}`                   | Update qty/unit item belanja                             |
| 60  | Child Nutrition   | `DELETE /ingredient-shopping-items/{id}`                  | Hapus item belanja                                       |
| 61  | Child Nutrition   | `GET /ingredients`                                        | Search master ingredients (autocomplete)                 |
| 62  | Photos & Contacts | `POST /mother-profiles/{id}/contacts`                     | Tambah contact/teman                                     |
| 63  | Photos & Contacts | `PATCH /child-photos/{id}`                                | Edit caption/visibility/review status foto               |
| 64  | Photos & Contacts | `DELETE /child-photos/{id}`                               | Hapus foto                                               |
| 65  | Medical           | `PATCH /medical-notes/{id}`                               | Edit catatan dokter                                      |
| 66  | Medical           | `DELETE /medical-notes/{id}`                              | Hapus catatan dokter                                     |
| 67  | Notification      | `DELETE /notifications/{id}`                              | Hapus 1 notifikasi                                       |
| 68  | Child Nutrition   | `PATCH /recipes/{id}/bookmark`                            | **[BARU v4]** Toggle bookmark recipe (dipisah dari 17)   |
| 69  | Caregiver         | `GET /mother-profiles/{id}/caregiver-engagements/revoked` | **[BARU v4]** Riwayat akses caregiver yang sudah dicabut |

**Endpoint 55 (v3) dihapus di v4** — fungsinya digabung ke endpoint 14 (lihat catatan di endpoint 14).

## Ringkasan Perubahan v3 → v4

| #       | Perubahan                                                                                                                                                                                                  |
| ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1       | Tambah field `birth_date`, `streak` di response                                                                                                                                                            |
| 2       | `notification_type` tetap opsional (revert dari "jadi wajib" di FIX v4)                                                                                                                                    |
| 4       | `analysis_type` & `age_range` jadi wajib; response dipangkas jadi cuma `z_score`+`age_months`                                                                                                              |
| 7       | `domains` diubah dari raw list jawaban jadi agregat `total_question`/`true_answer` per domain (dihitung di query, bukan di FE)                                                                             |
| 8       | `next_check_date` dihapus dari response                                                                                                                                                                    |
| 9       | `next_check_date` → `month_target`; `domains` disamakan formatnya dengan endpoint 7                                                                                                                        |
| 10      | Tambah field `month_target` di response                                                                                                                                                                    |
| 11      | Field `assessment_kpsp_answer` → `answer`; `month_target` dihapus dari request (dihitung server); diperjelas alur generate `child_development_report_id`                                                   |
| 14      | Semantik jadi **replace-all** (menggantikan endpoint 55 yang dihapus)                                                                                                                                      |
| 15      | `shopping.items[].ingredient_id` → `ingredient_name` (join ke `ingredients`)                                                                                                                               |
| 17      | Dipecah jadi 2 endpoint: 17 (`is_completed` saja) + 68 baru (`is_bookmarked` saja)                                                                                                                         |
| 18      | `main_ingredients` di response cuma yang `priority=1`, bukan semua diurutkan                                                                                                                               |
| 19      | Request body dihapus — panggil endpoint saja, otomatis set `priority=1`                                                                                                                                    |
| 20      | **Dihapus** (tidak dibutuhkan)                                                                                                                                                                             |
| 27      | Tambah field `phone_number` (dari `users.phone_number`) di response & contoh query                                                                                                                         |
| 30, 31  | `year` dihapus; `month` opsional (default bulan berjalan), filter nomor bulan tanpa syarat tahun                                                                                                           |
| 32      | Hapus catatan referensi `doctor_profiles`/`medical_relationship_id` lama dari deskripsi                                                                                                                    |
| 33      | Path jadi flat (`POST /medical-notes`, tanpa `child_id` di path); `child_id` ditambahkan ke body; `facility_location` jadi opsional; wajib hanya `child_id`, `doctor_name`, `recommendation`, `valid_date` |
| 34      | Response cuma `id, full_name, photo_url` (hapus `birth_date`, `gender`)                                                                                                                                    |
| 35      | Response cuma `id, shopping, menu` (hapus field nutrisi top-level)                                                                                                                                         |
| 39      | Query params dihapus (balik seperti v2, tanpa filter)                                                                                                                                                      |
| 41      | Query params dihapus; response difilter `is_review_required=false`                                                                                                                                         |
| 42      | Response cuma `id, child_id, photo_url, caption` (hapus `visibility`, `is_review_required`, `created_at`)                                                                                                  |
| 43      | `caption` jadi wajib                                                                                                                                                                                       |
| 56      | Field `target_*` dihapus; tambah `meal_times` (join ke `recipes` yang `is_completed=true`)                                                                                                                 |
| 65      | Field request body dijabarkan eksplisit satu per satu, semua opsional                                                                                                                                      |
| Baru    | Endpoint 68 (bookmark recipe), 69 (riwayat akses caregiver dicabut)                                                                                                                                        |
| Dihapus | Endpoint 20, 55                                                                                                                                                                                            |

**Perubahan dari v1/v2 (dipertahankan, lihat detail di narasi tiap endpoint)**: penamaan field → English snake_case, ID → UUID, format tanggal → `DD-MM-YYYY`, fix path/auth/request-body di berbagai endpoint, serta penambahan endpoint 45–67 di v2/v3.

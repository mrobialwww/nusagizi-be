# NusaGizi — Dokumentasi Endpoint API (v2, hasil review terhadap skema database)

Dokumen ini adalah revisi dari `NusaGizi_Endpoint_Documentation_Launchpad.pdf`, dicocokkan dengan skema database di migration `000005_simplify_medical_module`. Semua perbaikan sudah dikonfirmasi lewat sesi klarifikasi sebelumnya. Nomor endpoint asli dipertahankan agar mudah dibandingkan; endpoint tambahan diberi nomor baru (45, 46).

---

## 0. Konvensi Umum

- **Auth**: ditangani provider eksternal (Auth0). Setiap request terautentikasi membawa `user_id` dari token. Role (mother/caregiver) ditentukan dari ada-tidaknya row terkait di `mother_profiles`/`caregiver_profiles` untuk `user_id` tsb.
- **Format ID**: semua ID adalah UUID (string), termasuk di request body — bukan integer.
- **Format tanggal** di request/response body: `DD-MM-YYYY` (contoh: `25-07-2026`), kecuali dinyatakan lain.
- **Penamaan field**: English snake_case, mengikuti nama kolom database (contoh: `full_name`, `birth_date`).
- **Soft delete**: tabel dengan kolom `deleted_at` (`users`, `child`, `caregiver_engagements`) — endpoint DELETE terkait melakukan `UPDATE ... SET deleted_at = now()`, bukan hapus fisik.
- **Response POST**: kecuali dinyatakan lain, endpoint POST yang membuat resource baru mengembalikan minimal `{ "id": "<uuid>" }` di response body (bukan body kosong), supaya client tidak perlu re-fetch.
- **Base path**: contoh path di dokumen ini pakai konvensi resource-oriented (`/children/{child_id}/...`), sesuaikan dengan prefix API Anda (mis. `/api/v1/...`).

---

## 1. Modul: Auth & User Profile

### 23. Update profil user

- **Method & Path**: `PATCH /users/{user_id}`
- **Auth**: hanya bisa update `user_id` milik sendiri
- **Request Body**:

```json
{
    "full_name": "ibor alwan",
    "email": "robialwan8@gmail.com",
    "phone_number": "081234567899",
    "gender": "male"
}
```

- **Response Body**: `-`

### 25. Hapus user (soft delete)

- **Method & Path**: `DELETE /users/{user_id}`
- **Auth**: hanya bisa delete `user_id` milik sendiri
- **Request Body**: `-`
- **Response Body**: `-`
- **Catatan**: soft delete — set `users.deleted_at = now()`. Perlu didiskusikan terpisah (di luar scope dokumen ini) apakah child/data terkait ikut di-cascade-soft-delete atau tetap ada.

### 45. [BARU] Get profil user sendiri

- **Method & Path**: `GET /users/{user_id}` (atau `GET /users/me`)
- **Auth**: hanya bisa GET `user_id` milik sendiri
- **Request Body**: `-`
- **Response Body**: satu record dari `users`: `id, full_name, gender, email, phone_number, photo_url`
- **Catatan**: ditambahkan agar form edit di endpoint 23 punya data awal untuk di-load.

---

## 2. Modul: Child Profile

### 22. Tambah profil anak

- **Method & Path**: `POST /children`
- **Auth**: `mother_profile_id` diambil dari `user_id` yang login (bukan dari body)
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

- **Response Body**: `{ "id": "<child_id>" }`
- **Catatan [FIX]**:
    - Setiap item array di `allergies.*` disimpan sebagai row terpisah di `child_allergy_profile` dengan `category` sesuai key (`food`/`medicine`/`animal`/`others`) dan `allergen_name` = isi item.
    - `chronic_diseases[]` → row terpisah di `child_chronic_disease_profile.disease_name`.
    - `diets[]` → row terpisah di `child_diet_profile.diet_name`.
    - `favorite_foods[]` → row terpisah di `favorite_food_profile.food_name`.
    - `favorite_textures[]` → row terpisah di `favorite_texture_profile.texture_name`.
    - `food_frequency`, `food_goal`, `notes` langsung map ke `child.food_frequency_profile`, `child.food_goal_profile`, `child.notes_profile`.
    - Field diganti dari Indonesia (`name`, `date_of_birth`, `alergi`, `kondisi_kronis`, dst) ke English snake_case sesuai konvensi.

### 24. Update profil anak

- **Method & Path**: `PATCH /children/{child_id}`
- **Auth [FIX]**: divalidasi lewat `child_id` — pastikan `child.mother_profile_id` terhubung ke `mother_profiles.user_id` yang login. (Dokumentasi lama menyebut "berdasarkan user_id", tapi tabel `child` tidak punya kolom `user_id` sama sekali.)
- **Request Body**: sama seperti endpoint 22 (field opsional, hanya kirim yang mau diubah)
- **Response Body**: `-`
- **Catatan**: untuk field turunan (`allergies`, `chronic_diseases`, `diets`, `favorite_foods`, `favorite_textures`), strategi update disarankan replace-all (hapus semua row lama untuk child tsb, insert ulang dari array baru) supaya tidak perlu diff per-item.

### 46. [BARU] Hapus profil anak (soft delete)

- **Method & Path**: `DELETE /children/{child_id}`
- **Auth**: pastikan `child.mother_profile_id` milik `user_id` yang login
- **Request Body**: `-`
- **Response Body**: `-`
- **Catatan**: soft delete — set `child.deleted_at = now()`.

---

## 3. Modul: Child Growth

### 3. Get child_growth_reports terbaru

- **Method & Path**: `GET /children/{child_id}/growth-reports/latest`
- **Auth**: `child_growth_reports` sesuai `child_id`
- **Response Body**: satu record `child_growth_reports` (`id, measured_at, weight_kg, height_cm, head_circumference_cm`)

### 4. Get child_growth_analyses (filter range & analysis_type)

- **Method & Path**: `GET /children/{child_id}/growth-analyses`
- **Auth**: `child_growth_analyses` yang `child_growth_report_id`-nya milik `child_id` ybs
- **Query Params [FIX]**: `analysis_type`, `date_from`, `date_to`
- **Response Body**: list `child_growth_analyses` (`id, analysis_type, z_score, measured_at`)
- **Catatan**: `child_growth_analyses` tidak punya kolom tanggal sendiri — filter `date_from`/`date_to` dilakukan dengan join ke `child_growth_reports.measured_at`.

### 5. Get riwayat child_growth_reports

- **Method & Path**: `GET /children/{child_id}/growth-reports`
- **Auth**: sesuai `child_id`
- **Response Body**: list `child_growth_reports`

### 6. Tambah child_growth_reports

- **Method & Path**: `POST /children/{child_id}/growth-reports`
- **Auth**: hanya menyimpan `child_id` yang sesuai
- **Request Body [FIX]**:

```json
{
    "measured_at": "04-05-2026",
    "height_cm": 89,
    "weight_kg": 12.4,
    "head_circumference_cm": 47
}
```

- **Response Body**: `{ "id": "<growth_report_id>" }`

---

## 4. Modul: Child Development (KPSP)

**Algoritma bersama untuk endpoint 7, 9, 11, 12** (dikonfirmasi user):

```
KPSP_PERIODS = [3, 6, 9, 12, 15, 18, 21, 24, 30, 36, 42, 48, 54, 60]  // bulan
age_months   = umur anak sekarang, dihitung dari birth_date

month_target = periode TERBESAR di KPSP_PERIODS yang <= age_months
               jika age_months < 3 → month_target = 3

kpsp_score   = jumlah jawaban "true" dari 10 pertanyaan (skala 0–10)

next_period  = periode TERKECIL di KPSP_PERIODS yang > age_months
next_check_date =
    jika next_period ada  → birth_date + next_period bulan   (format DD-MM-YYYY)
    jika tidak ada (age_months >= 60, periode terakhir sudah dikerjakan)
                          → literal string "done"
```

Contoh: usia 7 bulan → `month_target=6`, `next_check_date` = usia 9 bulan (2 bulan lagi). Usia 2 bulan → `month_target=3`, `next_check_date` = usia 3 bulan (bulan depan).

### 7. Get child_development_reports terbaru + 4 domain

- **Method & Path**: `GET /children/{child_id}/development-reports/latest`
- **Auth**: sesuai `child_id`
- **Response Body**:

```json
{
  "id": "...",
  "kpsp_score": 8,
  "next_check_date": "25-08-2026",
  "created_at": "...",
  "domains": [
    { "nerve_name": "Gross motor skills", "answers": [{ "question_id": "...", "answer": true }] },
    { "nerve_name": "Fine motor skills", "answers": [...] },
    { "nerve_name": "Speech and language", "answers": [...] },
    { "nerve_name": "Socialization", "answers": [...] }
  ]
}
```

- **Catatan**: join `child_development_reports` → `assessment_kpsp_answers` → `assessment_kpsp_questions`, dikelompokkan per `nerve_name`.

### 8. Get riwayat asesmen KPSP

- **Method & Path**: `GET /children/{child_id}/development-reports`
- **Auth**: sesuai `child_id`
- **Response Body**: list `child_development_reports` (`id, kpsp_score, month_target, next_check_date, created_at`)

### 9. Get detail hasil asesmen KPSP

- **Method & Path**: `GET /development-reports/{child_development_report_id}`
- **Auth**: `child_development_reports` sesuai `child_id` milik mother/caregiver yang login
- **Response Body**:

```json
{
    "id": "...",
    "kpsp_score": 8,
    "next_check_date": "25-08-2026",
    "recommended_actions": [
        { "id": "...", "title": "...", "action_text": "..." }
    ],
    "domains": [
        /* sama seperti endpoint 7 */
    ]
}
```

- **Catatan**: `recommended_actions` yang tampil hanya yang terhubung ke pertanyaan yang dijawab `false` (indikasi keterlambatan), diambil lewat `assessment_attempt_recomendations`.

### 10. Get pertanyaan assessment_kpsp_questions

- **Method & Path**: `GET /assessment-kpsp-questions`
- **Query Params [FIX]**: `month_target` (wajib — sebelumnya tidak ada filter di dokumentasi lama, padahal soal berbeda per periode)
- **Response Body**: list `assessment_kpsp_questions` (`id, nerve_name, question_text, description`)

### 11. Simpan hasil asesmen KPSP (create)

- **Method & Path**: `POST /children/{child_id}/development-reports`
- **Request Body [FIX]** (`child_development_report_id` dihapus — di-generate server):

```json
{
    "month_target": 24,
    "list_answer": [
        {
            "assessment_kpsp_question_id": "021f1136-d24a-4e32-aaa8-cca24b9e608d",
            "assessment_kpsp_answer": true
        },
        {
            "assessment_kpsp_question_id": "03b996b5-1377-4a91-a0ba-5632526e1c09",
            "assessment_kpsp_answer": false
        }
    ]
}
```

- **Response Body**: sama seperti endpoint 9 (termasuk `id` yang baru dibuat)
- **Catatan**: server menghitung `kpsp_score` & `next_check_date` sesuai algoritma di atas, lalu insert ke `assessment_attempt_recomendations` untuk pertanyaan yang dijawab `false`. Dipanggil setelah semua 10 soal dijawab.

### 12. Update hasil asesmen KPSP terbaru

- **Method & Path**: `PATCH /development-reports/{child_development_report_id}`
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

- **Response Body**: sama seperti endpoint 9
- **Catatan**: `kpsp_score`, `next_check_date`, dan `assessment_attempt_recomendations` dihitung ulang. Dipanggil saat mother mengulangi asesmen.

### 13. Get checklist_milestone_tasks

- **Method & Path**: `GET /checklist-milestone-tasks`
- **Query Params**: `month_target`, `child_id` **[FIX: ditambahkan]**
- **Response Body**: list `checklist_milestone_tasks` (`id, nerve_name, task_description, is_checked`)
- **Catatan**: `is_checked` didapat dari left join ke `checklist_milestone_progress` dengan `child_id` yang dikirim, supaya FE tahu item mana yang sudah dicentang tanpa request terpisah.

### 14. Tandai checklist milestone

- **Method & Path**: `PATCH /children/{child_id}/checklist-milestone-progress`
- **Request Body [FIX: UUID, bukan integer]**:

```json
{
    "checklist_milestone_task_ids": [
        "b6f1e2d0-...-uuid-1",
        "8a2c3f10-...-uuid-2"
    ]
}
```

- **Response Body**: `-`

---

## 5. Modul: Child Nutrition

> Catatan umum: `child_nutrition_reports` tidak punya kolom tanggal eksplisit — "hari ini" = record terbaru per `child_id` (`ORDER BY created_at DESC LIMIT 1`), dengan asumsi sistem hanya generate maksimal 1 record per anak per hari (dikonfirmasi user).

### 15. Get nutrition report + shopping + menu hari ini (POV mother)

- **Method & Path**: `GET /children/{child_id}/nutrition/today`
- **Response Body**:

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
            { "id": "...", "ingredient_id": "...", "quantity": 0, "unit": "g" }
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

- **Catatan [FIX]**: field `allergen` di dokumentasi lama → nama kolom sebenarnya `is_alergen`.

### 16. Get menu hari ini

- **Method & Path**: `GET /children/{child_id}/daily-menus/today`
- **Response Body**: satu record `daily_menus` beserta list `recipes` (field sama seperti #15)

### 17. Update status/bookmark recipe

- **Method & Path**: `PATCH /recipes/{recipe_id}`
- **Request Body [FIX: `is_done` → `is_completed`]**:

```json
{ "is_completed": true, "is_bookmarked": true }
```

- **Response Body**: `-`

### 18. Get detail recipe

- **Method & Path**: `GET /recipes/{recipe_id}`
- **Response Body**: satu record `recipes` beserta list `main_ingredients` (diurutkan `priority`) dan `cooking_steps` (diurutkan `step_number`)

### 19. Tukar prioritas main_ingredients

- **Method & Path [FIX]**: `PATCH /recipes/{recipe_id}/main-ingredients/{main_ingredient_id}`
- **Request Body**: `{ "priority": 1 }`
- **Response Body**: `-`
- **Catatan**: `main_ingredient_id` sebelumnya tidak ada di spesifikasi endpoint — ditambahkan di path karena wajib untuk tahu ingredient mana yang ditukar. `main_ingredients` lain di recipe yang sama dengan `priority >= priority baru` digeser +1.

### 20. Get riwayat daily_menus (1 bulan terakhir)

- **Method & Path**: `GET /children/{child_id}/daily-menus?range=1m`
- **Response Body**: list `daily_menus`
- **Catatan**: filter tanggal memakai `child_nutrition_reports.created_at` (lewat relasi `child_nutrition_report_id`), karena `daily_menus` sendiri tidak punya kolom tanggal.

### 21. Get bookmark menu

- **Method & Path**: `GET /children/{child_id}/recipes/bookmarked`
- **Response Body**: list `recipes` dengan `is_bookmarked = true`, di-join dari `child_nutrition_reports` → `daily_menus` → `recipes` untuk memastikan milik `child_id` ybs

---

## 6. Modul: Photos & Contacts

### 37. Get list teman (contacts)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/contacts`
- **Response Body [FIX]**: list contacts, di-join ke `mother_profiles` → `users` untuk menampilkan nama/foto:

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

### 38. Hapus teman

- **Method & Path**: `DELETE /contacts/{contact_id}`
- **Response Body**: `-`
- **Catatan**: hard delete (tabel `contacts` tidak punya `deleted_at`).

### 39. Get semua foto anak milik mother

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/child-photos`
- **Response Body**: list `child_photos` (hanya `photo_url`)

### 40. Get foto anak dari suatu contact

- **Method & Path**: `GET /contacts/{contact_id}/child-photos`
- **Response Body**: list `child_photos` (hanya `photo_url`), via `photo_shared_with`

### 41. Get semua foto (gabungan 39+40)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/child-photos/all`
- **Response Body**: list `child_photos` (hanya `photo_url`)

### 42. Get detail foto

- **Method & Path**: `GET /child-photos/{child_photo_id}`
- **Response Body**: satu record penuh `child_photos`

### 43. Tambah foto (POV mother)

- **Method & Path**: `POST /children/{child_id}/photos`
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

- **Response Body**: `{ "id": "<child_photo_id>" }`
- **Catatan**: `visibility` diisi enum (`all`/`private`/`only`). `list_visibility` (array `contact_id`) hanya dipakai & wajib diisi ketika `visibility = "only"` — masing-masing di-insert sebagai row di `photo_shared_with`. `captions` → `caption` (sesuai nama kolom, singular).

### 44. Tambah foto (POV caregiver)

- **Method & Path**: `POST /children/{child_id}/photos`
- **Request Body**:

```json
{ "url": "https://example.com", "is_review_required": true }
```

- **Response Body**: `{ "id": "<child_photo_id>" }`
- **Catatan**: `is_review_required` wajib `true` untuk upload dari caregiver; `visibility` pakai default kolom (`private`).

---

## 7. Modul: Caregiver

### 26. Tambah akses caregiver (scan QR)

- **Method & Path [FIX]**: `POST /children/{child_id}/caregiver-engagements`
- **Auth**: `caregiver_profile_id` diambil dari `user_id` caregiver yang login; `child_id` dari path (hasil scan QR)
- **Request Body**: `-`
- **Response Body**: `{ "id": "<caregiver_engagement_id>" }`

### 27. Get list akses caregiver (POV mother)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/caregiver-engagements`
- **Response Body**: list `caregiver_engagements`
- **Contoh query [FIX — PK diperbaiki jadi `id`]**:

```sql
SELECT
  c.id            AS child_id,
  c.full_name     AS child_name,
  ce.id           AS caregiver_engagement_id,
  ce.caregiver_profile_id,
  cp.id           AS caregiver_profile_id_check,
  u.full_name     AS caregiver_name
FROM child c
JOIN caregiver_engagements ce ON ce.child_id = c.id
JOIN caregiver_profiles cp    ON cp.id = ce.caregiver_profile_id
JOIN users u                  ON u.id = cp.user_id
WHERE c.mother_profile_id = ?
  AND ce.deleted_at IS NULL;
```

### 28. Hapus akses caregiver

- **Method & Path**: `DELETE /caregiver-engagements/{caregiver_engagement_id}`
- **Response Body**: `-`
- **Catatan**: soft delete — set `caregiver_engagements.deleted_at = now()`.

### 29. Get profil caregiver

- **Method & Path**: `GET /caregiver-profiles/{caregiver_profile_id}` (atau `/caregiver-profiles/me`)
- **Response Body**: satu record `caregiver_profiles`
- **Catatan**: relasi `caregiver_profiles`↔`users` bersifat 1:1, jadi tidak ada makna "terbaru" — cukup 1 record per user.

### 34. Get list child dari caregiver_engagement

- **Method & Path**: `GET /caregiver-profiles/{caregiver_profile_id}/children`
- **Response Body**: list `child` (`id, full_name, birth_date, gender, photo_url`)
- **Contoh query [FIX — `photo` → `photo_url`, PK diperbaiki]**:

```sql
SELECT c.id, c.full_name, c.birth_date, c.gender, c.photo_url
FROM caregiver_engagements ce
JOIN child c ON c.id = ce.child_id
WHERE ce.caregiver_profile_id = ?
  AND ce.deleted_at IS NULL;
```

### 35. Get nutrition report + shopping + menu hari ini (POV caregiver)

- **Method & Path**: `GET /children/{child_id}/nutrition/today` (endpoint sama dengan #15, akses via caregiver engagement bukan mother_profile_id)
- **Response Body**: sama seperti endpoint 15

### 36. Get detail recipe (POV caregiver)

- **Method & Path**: `GET /recipes/{recipe_id}` (endpoint sama dengan #18)
- **Response Body**: sama seperti #18, tanpa `calories, protein, carbohydrate, fat`

---

## 8. Modul: Medical

### 30. Get list medical_notes aktif per anak (POV mother)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/medical-notes?status=active`
- **Query Params [FIX: filter ditambahkan]**: `valid_date >= now()`
- **Response Body**: list, tiap item: `{ valid_date, recommendation, created_at, child_name, prohibition_count, allergy_count }`
- **Catatan**: `prohibition_count`/`allergy_count` dihitung dari `medical_restrictions` per `medical_note_id`, dikelompokkan per `type`.

### 31. Get riwayat medical_notes (sudah tidak aktif)

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/medical-notes?status=history`
- **Query Params**: filter `valid_date < now()`
- **Response Body**: sama seperti #30

### 32. Get rincian medical_notes

- **Method & Path [FIX]**: `GET /medical-notes/{medical_note_id}`
- **Auth [FIX]**: divalidasi lewat `medical_note_id` → pastikan `child_id` pemilik note tsb dimiliki oleh mother/caregiver yang login. (Dokumentasi lama merujuk `doctor_profiles`/`medical_relationship_id` yang sudah dihapus di migration `000005_simplify_medical_module`.)
- **Response Body**: satu record `medical_notes` beserta list `medical_restrictions`, `daily_nutrition_targets`, dan `child.full_name`

### 33. Tambah catatan dokter

- **Method & Path**: `POST /children/{child_id}/medical-notes`
- **Request Body**:

```json
{
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

- **Response Body**: `{ "id": "<medical_note_id>" }`
- **Catatan**: `prohibitions[]` → row `medical_restrictions` dengan `type = 'prohibition'`; `allergies[]` → row `medical_restrictions` dengan `type = 'allergy'`.

---

## 9. Modul: Notification

### 2. Get seluruh notifikasi

- **Method & Path**: `GET /users/{user_id}/notifications`
- **Response Body**: list `notification` (`id, title, message, notification_type, created_at`)

---

## 10. Modul: Dashboard (gabungan)

### 1. Get seluruh anak + tinggi, berat, skor KPSP, protein

- **Method & Path**: `GET /mother-profiles/{mother_profile_id}/children/summary`
- **Response Body**: list, tiap item:

```json
{
    "id": "...",
    "full_name": "...",
    "height_cm": 89,
    "weight_kg": 12.4,
    "kpsp_score": 8,
    "protein": 30,
    "target_protein": 35
}
```

- **Catatan [FIX]**: `height_cm`/`weight_kg` diambil dari `child_growth_reports` **terbaru**, `kpsp_score` dari `child_development_reports` **terbaru**, `protein` (+`target_protein`) dari `child_nutrition_reports` **terbaru** — masing-masing per `child_id` (dokumentasi lama tidak menyebut "terbaru" secara eksplisit padahal relasinya 1:N).

---

## Ringkasan Semua Perubahan dari Dokumentasi Lama

| #          | Perubahan                                                                                                                                           |
| ---------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| Global     | Penamaan field → English snake_case; ID → UUID; format tanggal → `DD-MM-YYYY`                                                                       |
| 4          | Tambah query param `date_from`/`date_to`, filter via join ke `child_growth_reports`                                                                 |
| 6          | Field body diperbaiki: `measurement_at/height/weight/head_circumference` → `measured_at/height_cm/weight_kg/head_circumference_cm`                  |
| 9, 11, 12  | Algoritma `kpsp_score` & `next_check_date` didefinisikan eksplisit; `recommended_actions` hanya dari jawaban `false`                                |
| 11         | `child_development_report_id` dihapus dari request (di-generate server)                                                                             |
| 13         | Tambah param `child_id`, tambah field `is_checked` di response                                                                                      |
| 14         | ID diperbaiki jadi UUID                                                                                                                             |
| 15, 16, 35 | `allergen` → `is_alergen`                                                                                                                           |
| 17         | `is_done` → `is_completed`                                                                                                                          |
| 19         | Tambah `main_ingredient_id` di path                                                                                                                 |
| 22, 24     | Field alergi/kondisi kronis/dll diubah jadi array eksplisit; field Indonesia → English                                                              |
| 24         | Authorization diperbaiki (bukan berbasis `user_id`, tapi `child_id` + relasi mother)                                                                |
| 27, 34     | Contoh SQL diperbaiki (PK sebenarnya bernama `id`, bukan `child_id`/`caregiver_engagement_id`/dst sebagai nama kolom PK); `c.photo` → `c.photo_url` |
| 30         | Tambah filter eksplisit `valid_date >= now()`                                                                                                       |
| 32         | Authorization diperbaiki total — `doctor_profiles`/`medical_relationship_id` sudah tidak ada di skema                                               |
| 37         | Response ditambah join ke `mother_profiles`/`users` supaya ada nama & foto                                                                          |
| 43         | `captions` → `caption`; `visibility` (daftar nama) dipisah jadi `visibility` (enum) + `list_visibility` (array contact_id)                          |
| 45         | **[Baru]** GET profil user sendiri                                                                                                                  |
| 46         | **[Baru]** DELETE (soft delete) profil anak                                                                                                         |

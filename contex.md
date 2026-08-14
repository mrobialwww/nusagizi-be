# Dokumentasi Desain Skema `recipes` & `child_nutrition_recipes`

_(Migration `000004_create_core_domain_tables`)_

---

## 1. Latar Belakang & Keputusan Awal

Perubahan skema ini berangkat dari evolusi relasi:

1. `daily_menus` dihapus karena relasi `child_nutrition_reports → daily_menus` sebenarnya **1:1**, sehingga tidak perlu tabel terpisah.
2. Efeknya, `child_nutrition_reports` jadi berelasi langsung ke `recipes` dengan kardinalitas **1:N**.
3. Kardinalitas ini kemudian diubah menjadi **M:N**, karena satu `recipes` bisa dipakai oleh lebih dari satu `child_nutrition_reports` (reuse resep antar hari/anak). Junction table `child_nutrition_recipes` menampung `child_nutrition_report_id`, `recipe_id`, `portions_consumed`, `created_at`, `updated_at`.
4. `meal_time` dan `is_bookmarked` **sengaja tetap di tabel master `recipes`**, bukan di junction, dengan alasan:
    - `is_bookmarked` → status bookmark ingin lintas anak (shared).
    - `meal_time` → bersifat **kaku/immutable** sejak resep pertama kali dibuat.

**Kesimpulan awal:** arah normalisasi ini sudah tepat dan best practice. Yang perlu diperbaiki hanya di detail constraint dan mekanisme lifecycle data.

---

## 2. Klarifikasi 3 Pertanyaan Awal

| #   | Topik                                            | Keputusan                                                                                                                                                                                                                                                                                                                                                                                                      |
| --- | ------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Access control parent → recipes anak             | Sudah diatur di backend (app-layer), bukan tanggung jawab DB schema.                                                                                                                                                                                                                                                                                                                                           |
| 2   | Tipe data `portions_consumed`                    | `FLOAT` → **`NUMERIC(3,2)`**. Alasan: `FLOAT` adalah _binary floating point_, sehingga nilai seperti `0.33`/`0.66` tidak eksak dan berisiko mismatch dengan `CHECK (... IN (0, 0.33, 0.66, 1))`. `NUMERIC` adalah _exact decimal_, aman untuk constraint pembanding nilai persis.                                                                                                                              |
| 3   | "CNR dihapus → junction & recipes ikut terhapus" | **Perlu diluruskan**: FK `recipe_id ... ON DELETE CASCADE` di skema awal justru mengalir **dari `recipes` ke junction**, bukan sebaliknya. Untuk "hapus recipes yang jadi _orphan_ setelah junction row hilang", ini **tidak bisa** diwakilkan ke FK constraint biasa — harus logika eksplisit (cek dulu apakah recipe masih punya relasi lain) yang dijalankan di dalam transaksi yang sama saat CNR dihapus. |

---

## 3. Edge Case & Keputusan Final (Q&A)

### Q1 — Konflik `is_bookmarked` vs auto-delete saat orphan

**Isu:** Resep yang di-bookmark bisa ikut terhapus tanpa sepengetahuan user kalau kebetulan jadi orphan (CNR terakhirnya dihapus karena regenerate di hari yang sama).

> **A1 (keputusan):** Tidak dipedulikan. Ikuti alur mekanisme apa adanya — kalau resep ber-bookmark ikut terhapus akibat salah satu dari 5 mekanisme, itu **diterima sebagai risiko yang disengaja**, tidak perlu pengecualian khusus.

**Dampak ke implementasi:** Query cleanup orphan **tidak perlu** filter `AND is_bookmarked = false`. Cukup cek keberadaan relasi di junction saja (lihat §5).

_Follow-up question yang perlu dipertimbangkan ke depan:_ karena keputusan ini final dan disengaja, apakah perlu ada _soft warning_ di level produk/UX (misal toast "resep di-bookmark akan hilang jika tidak direuse hari ini")? Ini bukan masalah skema, murni pertimbangan produk — silakan diabaikan kalau memang tidak relevan.

---

### Q2 — Uniqueness "satu `meal_time` per CNR" belum dijamin di level DB

**Isu:** `meal_time` melekat di `recipes`, sehingga tidak bisa dibuat `UNIQUE (child_nutrition_report_id, meal_time)` lintas tabel.

> **A2 (keputusan):** Dilaksanakan — denormalisasi `meal_time` ke `child_nutrition_recipes`.

**Kenapa ini aman dilakukan** (tidak menimbulkan risiko _data staleness_/anomali update):

Denormalisasi berbahaya _hanya_ ketika kolom yang diduplikasi bisa **berubah nilainya** setelah disalin — karena itu membuka celah dua sumber data saling tidak sinkron (update di satu tempat, lupa update di tempat lain). Dalam kasus ini:

- `meal_time` bersifat **immutable sejak resep dibuat** (sudah dikonfirmasi di awal diskusi — "kaku di waktu dimana dia dibuat pertama kali").
- Karena nilainya tidak pernah di-_update_ setelah insert, tidak ada window waktu di mana `recipes.meal_time` dan `child_nutrition_recipes.meal_time` bisa berbeda.
- Nilainya disalin **sekali saat insert** (saat recipe pertama kali dikaitkan ke suatu CNR) dan tidak pernah perlu disinkronkan ulang — sehingga secara teknis ini bukan "denormalisasi berisiko", melainkan lebih tepat disebut _"snapshot atribut yang immutable by design"_.
- Sebagai bonus, ini justru memungkinkan constraint DB-level (`UNIQUE`) yang tidak bisa dicapai kalau `meal_time` tetap hanya di tabel master.

_Follow-up question:_ apakah ada skenario di mana `meal_time` sebuah `recipes` **bisa** diubah lewat fitur edit resep di masa depan (misal admin/parent mengedit typo)? Kalau ya, maka asumsi "immutable" ini perlu ditinjau ulang — kemungkinan perlu trigger sync atau larangan edit `meal_time` setelah resep pernah dipakai di junction manapun.

---

### Q3 — Konkurensi saat generate CNR kedua di hari yang sama

**Isu:** Race condition potensial saat menentukan "CNR mana yang merupakan CNR hari ini yang harus digantikan".

> **A3 (keputusan):** Sudah ditangani di backend (locking/transaksi di application layer), bukan lewat query DB murni.

**Catatan:** Tidak ada perubahan skema yang diperlukan untuk ini. Disarankan tetap ada `UNIQUE (child_id, report_date)` (atau setara) di tabel `child_nutrition_reports` sebagai _safety net_ di level DB — tapi karena locking sudah ditangani di backend, ini opsional, bukan blocker.

_Follow-up question:_ apakah constraint `UNIQUE (child_id, report_date)` tersebut sudah ada di migration tabel `child_nutrition_reports`? Ini tidak terlihat di file migration 000004 (tabel tersebut didefinisikan di migration lain), jadi perlu dicek terpisah agar locking di backend punya pengaman ganda di level DB.

---

### Q4 — Apakah `recipes` di-dedup lintas anak?

**Isu:** Apakah bookmark "lintas anak" benar-benar bermakna kalau resep tidak di-_share_ secara fisik antar anak.

> **A4 (keputusan/klarifikasi):**
>
> - Recipes **bisa** dilihat/dipakai bersama antar anak-anak dari **parent yang sama** (meskipun awalnya dibuat untuk satu anak).
> - Recipes **tidak bisa** dilihat lintas parent yang berbeda.
> - Mekanisme "bookmark lintas anak" sesederhana `SELECT * FROM recipes WHERE is_bookmarked = true`, karena parent memang hanya bisa mengakses recipes yang terhubung ke salah satu anaknya (sudah dijamin di backend), sehingga tidak akan pernah terjadi sharing recipes antar parent.

**Dampak:** Desain M:N `recipes ↔ child_nutrition_reports` sudah cukup untuk mendukung ini — satu `recipes` row memang bisa dipakai lintas anak (lintas CNR yang berbeda `child_id`, selama parent sama), tidak perlu kolom `child_id` tambahan di `recipes` maupun `child_nutrition_recipes`.

_Follow-up question penting:_ karena `is_bookmarked` adalah **satu kolom boolean pada satu row `recipes`** yang bisa dipakai bersama oleh beberapa anak (siblings), berarti **status bookmark otomatis shared** — kalau anak A membookmark suatu resep, resep itu otomatis akan tampak "ter-bookmark" juga ketika dilihat dari konteks anak B (karena row fisiknya sama). Apakah ini memang perilaku yang diinginkan (bookmark benar-benar _shared per keluarga_, bukan per anak)? Kalau iya, tidak perlu perubahan apa pun — desain saat ini sudah otomatis mendukungnya. Kalau ternyata yang diinginkan adalah bookmark _per anak per resep_, maka `is_bookmarked` **harus dipindah** dari `recipes` ke `child_nutrition_recipes` (atau tabel bookmark terpisah `child_id + recipe_id`), yang berarti kontradiksi dengan keputusan awal "bookmark diletakkan di master karena ingin lintas anak" — mohon dikonfirmasi ulang karena ini krusial dan mengubah skema secara signifikan.

---

## 4. Revisi Migration `000004` (Bagian yang Berubah)

Hanya bagian yang relevan dengan perbaikan (tidak mengubah seluruh file):

```sql
-- child_nutrition_reports <-> recipes (N:M)
CREATE TABLE child_nutrition_recipes (
    child_nutrition_report_id  UUID NOT NULL REFERENCES child_nutrition_reports(id) ON DELETE CASCADE,
    recipe_id                  UUID NOT NULL REFERENCES recipes(id), -- CATATAN: cascade dihapus, lihat §5
    meal_time                  meal_time_type NOT NULL, -- BARU: snapshot immutable dari recipes.meal_time
    portions_consumed          NUMERIC(3,2) NOT NULL DEFAULT 0
                                CHECK (portions_consumed IN (0, 0.33, 0.66, 1)),
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_nutrition_report_id, recipe_id)
);

CREATE INDEX idx_child_nutrition_recipes_report_id ON child_nutrition_recipes(child_nutrition_report_id);
CREATE INDEX idx_child_nutrition_recipes_recipe_id ON child_nutrition_recipes(recipe_id);

-- BARU: menjamin satu meal_time hanya muncul sekali per CNR
CREATE UNIQUE INDEX uq_cnr_report_mealtime
    ON child_nutrition_recipes (child_nutrition_report_id, meal_time);

CREATE TRIGGER trg_child_nutrition_recipes_updated_at BEFORE UPDATE ON child_nutrition_recipes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

**Ringkasan perubahan dari versi asli:**

| Perubahan                                    | Alasan                                                                                                      |
| -------------------------------------------- | ----------------------------------------------------------------------------------------------------------- |
| `ON DELETE CASCADE` pada `recipe_id` dihapus | Arahnya salah untuk kebutuhan "orphan cleanup" — lihat §5, harus ditangani eksplisit di transaksi aplikasi. |
| Tambah kolom `meal_time`                     | Sesuai keputusan A2, aman karena immutable (lihat §3 Q2).                                                   |
| Tambah `UNIQUE INDEX uq_cnr_report_mealtime` | Menjamin invariant "1 meal_time per CNR" di level DB, bukan hanya app-code.                                 |
| `portions_consumed FLOAT` → `NUMERIC(3,2)`   | Presisi exact decimal untuk match `CHECK` constraint.                                                       |

---

## 5. Contoh Query Pendukung

### a. Cleanup resep orphan saat CNR lama dihapus (dalam satu transaksi)

Sesuai A1 — **tidak** ada pengecualian untuk `is_bookmarked`:

```sql
BEGIN;

DELETE FROM child_nutrition_reports WHERE id = :old_cnr_id; -- junction ikut terhapus via cascade

DELETE FROM recipes
WHERE id = ANY(:recipe_ids_from_old_cnr)
  AND NOT EXISTS (
    SELECT 1 FROM child_nutrition_recipes cnr
    WHERE cnr.recipe_id = recipes.id
  );

COMMIT;
```

### b. Mengganti resep untuk `meal_time` tertentu pada suatu CNR

Memanfaatkan `uq_cnr_report_mealtime` sebagai pengaman:

```sql
BEGIN;

DELETE FROM child_nutrition_recipes
WHERE child_nutrition_report_id = :cnr_id
  AND meal_time = 'pagi';

INSERT INTO child_nutrition_recipes
    (child_nutrition_report_id, recipe_id, meal_time, portions_consumed)
VALUES
    (:cnr_id, :new_recipe_id, 'pagi', 0);

COMMIT;
```

### c. Ambil semua resep ter-bookmark (lintas anak, sesuai A4)

```sql
SELECT r.*
FROM recipes r
JOIN child_nutrition_recipes cnr ON cnr.recipe_id = r.id
JOIN child_nutrition_reports cn ON cn.id = cnr.child_nutrition_report_id
WHERE r.is_bookmarked = true
  AND cn.child_id = ANY(:children_ids_of_this_parent); -- filter akses tetap di backend
```

---

## 6. Daftar Follow-up yang Masih Terbuka

1. **(dari Q1)** Perlu pertimbangan produk/UX terkait bookmark yang bisa hilang tanpa peringatan — opsional, bukan isu skema.
2. **(dari Q2)** Pastikan tidak ada fitur "edit `meal_time` resep yang sudah pernah dipakai" di masa depan — kalau ada, asumsi immutable perlu ditinjau ulang.
3. **(dari Q3)** Konfirmasi apakah `UNIQUE (child_id, report_date)` (atau setara) sudah ada di migration tabel `child_nutrition_reports`, sebagai _safety net_ DB-level di luar locking backend.
4. **(dari Q4 — paling krusial)** Konfirmasi ulang: apakah `is_bookmarked` memang dimaksudkan **shared antar siblings** (satu keluarga), bukan per-anak? Ini menentukan apakah kolom tetap di `recipes` (desain saat ini) atau harus dipindah ke junction/tabel terpisah.

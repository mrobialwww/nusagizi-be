-- =========================================================================
-- SEEDER MEDICAL NOTES (JAN - DES) DENGAN PURE SQL CTE
-- =========================================================================

WITH child_ids(id) AS (
    -- 1. Daftar 3 ID Anak dari children.sql
    VALUES 
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::UUID),
    ('61fa1c7e-81f9-4960-b226-af43f57e8b62'::UUID),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::UUID)
),
months(m) AS (
    -- 2. Daftar Bulan (Januari s/d Desember)
    SELECT generate_series(1, 12)
),
note_multipliers(mult) AS (
    -- 3. Multiplier untuk mendapatkan 1 atau 2 notes
    VALUES (1), (2)
),
notes_base AS (
    -- 4. Merakit kerangka baris medical notes (Maret, Juni, Sept, Des punya 2 notes)
    SELECT 
        c.id AS child_id,
        m.m AS month,
        n.mult AS note_idx,
        gen_random_uuid() AS note_id
    FROM child_ids c
    CROSS JOIN months m
    CROSS JOIN note_multipliers n
    WHERE NOT (n.mult = 2 AND m.m NOT IN (3, 6, 9, 12)) 
),
inserted_notes AS (
    -- 5. EKSEKUSI INSERT tabel utama (medical_notes)
    INSERT INTO medical_notes (id, child_id, doctor_name, facility_name, recommendation, valid_until)
    SELECT 
        note_id,
        child_id,
        (ARRAY['dr. Tirta', 'dr. Mesty', 'dr. Budi', 'dr. Ani', 'dr. Richard'])[floor(random() * 5 + 1)],
        (ARRAY['RSIA Hermina', 'Klinik Tumbuh Kembang', 'RS Sardjito', 'Puskesmas Melati'])[floor(random() * 4 + 1)],
        'Pemberian nutrisi harus dijaga dengan baik sesuai anjuran dan pantangan gizi anak.',
        make_date(2026, month, (floor(random() * 28 + 1))::INT)
    FROM notes_base
    RETURNING id
),
inserted_allergies AS (
    -- 6. EKSEKUSI INSERT alergi (1 buah per medical note)
    INSERT INTO medical_restrictions (medical_note_id, type, item_name)
    SELECT 
        id,
        'allergy',
        (ARRAY['Kacang Tanah', 'Susu Sapi', 'Telur', 'Seafood', 'Udang', 'Kedelai'])[floor(random() * 6 + 1)]
    FROM inserted_notes
),
inserted_prohibitions AS (
    -- 7. EKSEKUSI INSERT pantangan (1 buah per medical note)
    INSERT INTO medical_restrictions (medical_note_id, type, item_name)
    SELECT 
        id,
        'prohibition',
        (ARRAY['Makanan Manis', 'Makanan Keras', 'Serat Tinggi', 'Garam Berlebih'])[floor(random() * 4 + 1)]
    FROM inserted_notes
)
-- 8. EKSEKUSI AKHIR INSERT 4 nutrisi per medical note
INSERT INTO daily_nutrition_targets (medical_note_id, nutrient, quantity)
SELECT 
    id,
    nutrient_type,
    floor(random() * 80 + 20)
FROM inserted_notes
CROSS JOIN unnest(ARRAY['calorie', 'protein', 'fat', 'carbohydrate']::nutrient_type[]) AS nutrient_type;

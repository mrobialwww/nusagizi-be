-- =====================================================================
-- SEED DATA: favorite_food_profiles, favorite_texture_profiles,
--            child_diet_profiles, child_chronic_disease_profiles,
--            child_allergy_profiles
-- =====================================================================

-- =============================================
-- 1. favorite_food_profiles
-- =============================================
INSERT INTO favorite_food_profiles (id, child_id, food_name, created_at, updated_at)
VALUES
    -- M. Razky (27 bulan, suka makanan beragam)
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'Nasi Tim Ayam', now(), now()),
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'Pisang Goreng', now(), now()),
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'Telur Rebus', now(), now()),
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'Kentang Panggang', now(), now()),

    -- Budi (17 bulan, masih suka makanan lembut)
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'Bubur Kentang', now(), now()),
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'Pisang Halus', now(), now()),
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'Wortel Rebus', now(), now()),
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'Alpukat', now(), now()),

    -- Ani (37 bulan, makanan keluarga)
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Nasi Goreng', now(), now()),
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Ikan Bakar', now(), now()),
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Sayur Sop', now(), now()),
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Apel Iris', now(), now()),
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Susu UHT', now(), now())
ON CONFLICT DO NOTHING;

-- =============================================
-- 2. favorite_texture_profiles
-- =============================================
INSERT INTO favorite_texture_profiles (id, child_id, texture_name, created_at, updated_at)
VALUES
    -- M. Razky (27 bulan, transisi ke makanan keluarga)
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'Finger food', now(), now()),
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'Potongan kecil', now(), now()),

    -- Budi (17 bulan, masih butuh makanan lembut)
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'Bubur halus', now(), now()),
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'Puree', now(), now()),
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'Lunak', now(), now()),

    -- Ani (37 bulan, makanan keluarga penuh)
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Makanan keluarga', now(), now()),
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Finger food', now(), now()),
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Camilan sehat', now(), now())
ON CONFLICT DO NOTHING;

-- =============================================
-- 3. child_diet_profiles
-- =============================================
INSERT INTO child_diet_profiles (id, child_id, diet_name, created_at, updated_at)
VALUES
    -- M. Razky
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'TETP', now(), now()),
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'MPASI', now(), now()),

    -- Budi
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'MPASI', now(), now()),

    -- Ani
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'TETP', now(), now()),
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Gizi Seimbang', now(), now())
ON CONFLICT DO NOTHING;

-- =============================================
-- 4. child_chronic_disease_profiles
-- =============================================
INSERT INTO child_chronic_disease_profiles (id, child_id, disease_name, created_at, updated_at)
VALUES
    -- M. Razky: ada riwayat asma ringan
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'Asma ringan', now(), now()),

    -- Budi: sehat, tidak ada penyakit kronis (tidak ada data)

    -- Ani: ada riwayat eksim
    (gen_random_uuid(), '88e88505-c619-4e95-9e7a-7749d9539468', 'Eksim', now(), now())
ON CONFLICT DO NOTHING;

-- =============================================
-- 5. child_allergy_profiles
-- =============================================
INSERT INTO child_allergy_profiles (id, child_id, category, allergen_name, created_at, updated_at)
VALUES
    -- M. Razky: alergi kacang (food), debu (other)
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'food', 'Kacang Tanah', now(), now()),
    (gen_random_uuid(), '0c4b4722-a02d-46d4-8d25-d420ba2b7f89', 'other', 'Debu', now(), now()),

    -- Budi: alergi susu sapi (food)
    (gen_random_uuid(), '61fa1c7e-81f9-4960-b226-af43f57e8b62', 'food', 'Susu Sapi', now(), now())
;

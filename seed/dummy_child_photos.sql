-- Seed script untuk child_photos dan photo_shares dummy
-- Harus dijalankan SETELAH seed/dummy_contacts.sql

DO $$ 
DECLARE 
    main_mother_id UUID; 
    main_child_id UUID;
    
    d1_mother UUID; d1_child UUID;
    d2_mother UUID; d2_child UUID;
    d3_mother UUID; d3_child UUID;
    
    i INT;
BEGIN
    SELECT id INTO main_mother_id FROM mother_profiles ORDER BY created_at ASC LIMIT 1;
    
    -- Bersihkan data foto anak percobaan sebelumnya agar rerrunable
    main_child_id := '0c4b4722-a02d-46d4-8d25-d420ba2b7f89';
    DELETE FROM child_photos WHERE child_id = main_child_id;

    -- Ambil ID mother profile dummy
    SELECT m.id INTO d1_mother FROM mother_profiles m JOIN users u ON m.user_id = u.id WHERE u.email = 'bundadummy1@nusagizi.com';
    SELECT m.id INTO d2_mother FROM mother_profiles m JOIN users u ON m.user_id = u.id WHERE u.email = 'bundadummy2@nusagizi.com';
    SELECT m.id INTO d3_mother FROM mother_profiles m JOIN users u ON m.user_id = u.id WHERE u.email = 'bundadummy3@nusagizi.com';
    
    -- Bersihkan anak dummy jika sebelumnya sudah ada
    DELETE FROM children WHERE mother_profile_id IN (d1_mother, d2_mother, d3_mother) AND full_name LIKE 'Anak Bunda Dummy%';


    -- ==========================================
    -- BUNDA DUMMY 1 (3 Foto 'all')
    -- ==========================================
    IF d1_mother IS NOT NULL THEN
        INSERT INTO children (mother_profile_id, full_name, gender, birth_date, streak_days) 
        VALUES (d1_mother, 'Anak Bunda Dummy 1', 'male', '2024-01-01', 0) RETURNING id INTO d1_child;
        
        FOR i IN 1..3 LOOP
            INSERT INTO child_photos (child_id, photo_url, caption, visibility, is_review_required)
            VALUES (d1_child, 'social/' || d1_mother || '/child-2.' || i || '.jpg', 'Foto anak 1 bagian ' || i, 'all', false);
        END LOOP;
    END IF;


    -- ==========================================
    -- BUNDA DUMMY 2 (3 Foto 'all')
    -- ==========================================
    IF d2_mother IS NOT NULL THEN
        INSERT INTO children (mother_profile_id, full_name, gender, birth_date, streak_days) 
        VALUES (d2_mother, 'Anak Bunda Dummy 2', 'female', '2024-01-02', 0) RETURNING id INTO d2_child;
        
        FOR i IN 1..3 LOOP
            INSERT INTO child_photos (child_id, photo_url, caption, visibility, is_review_required)
            VALUES (d2_child, 'social/' || d2_mother || '/child-3.' || i || '.jpg', 'Foto anak 2 bagian ' || i, 'all', false);
        END LOOP;
    END IF;


    -- ==========================================
    -- BUNDA DUMMY 3 (2 Foto 'all', 1 Foto 'private')
    -- ==========================================
    IF d3_mother IS NOT NULL THEN
        INSERT INTO children (mother_profile_id, full_name, gender, birth_date, streak_days) 
        VALUES (d3_mother, 'Anak Bunda Dummy 3', 'male', '2024-01-03', 0) RETURNING id INTO d3_child;
        
        FOR i IN 1..3 LOOP
            IF i <= 2 THEN
                INSERT INTO child_photos (child_id, photo_url, caption, visibility, is_review_required)
                VALUES (d3_child, 'social/' || d3_mother || '/child-4.' || i || '.jpg', 'Foto anak 3 (All) bagian ' || i, 'all', false);
            ELSE
                INSERT INTO child_photos (child_id, photo_url, caption, visibility, is_review_required)
                VALUES (d3_child, 'social/' || d3_mother || '/child-4.' || i || '.jpg', 'Foto anak 3 (Private) bagian ' || i, 'private', false);
            END IF;
        END LOOP;
    END IF;


    -- ==========================================
    -- AKUN UTAMA: Anda Sendiri (2 is_review_required = true, 1 = false)
    -- Menggunakan ID anak "M. Razky" (dari seed/children.sql)
    -- ==========================================
    
    FOR i IN 1..3 LOOP
        IF i <= 2 THEN
            INSERT INTO child_photos (child_id, photo_url, caption, visibility, is_review_required)
            VALUES (main_child_id, 'social/' || main_mother_id || '/child-1.' || i || '.jpg', 'Foto anak utamaku ' || i, 'all', true);
        ELSE
            INSERT INTO child_photos (child_id, photo_url, caption, visibility, is_review_required)
            VALUES (main_child_id, 'social/' || main_mother_id || '/child-1.' || i || '.jpg', 'Foto anak utamaku ' || i, 'all', false);
        END IF;
    END LOOP;

    -- Catatan: Tabel photo_shares sengaja tidak diisi pada dummy ini karena
    -- skenario yang Anda tentukan sekarang terbatas di 'all' dan 'private' saja.

END $$;

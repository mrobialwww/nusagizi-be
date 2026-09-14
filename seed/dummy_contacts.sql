-- Seed script: Membuat user dummy, profil ibu, dan entri khusus daftar contact
-- Script ini sifatnya idempotent (aman di-run berkali-kali)

DO $$ 
DECLARE 
    main_mother_id UUID;
    
    d1_user UUID; d1_mother UUID; 
    d2_user UUID; d2_mother UUID; 
    d3_user UUID; d3_mother UUID; 
BEGIN
    -- Pembersihan data dummy
    DELETE FROM users WHERE auth0_id LIKE 'auth0|dummy_bunda_%';
    
    SELECT id INTO main_mother_id FROM mother_profiles ORDER BY created_at ASC LIMIT 1;


    -- DUMMY BUNDA 1
    INSERT INTO users (auth0_id, email, full_name, gender, phone_number, photo_url)
    VALUES ('auth0|dummy_bunda_1', 'bundadummy1@nusagizi.com', 'Bunda Dummy Satu', 'female', '081234567801', NULL)
    RETURNING id INTO d1_user;
    
    d1_mother := 'c47c83e4-2df0-4418-b651-05448b34c245';
    INSERT INTO mother_profiles (id, user_id) VALUES (d1_mother, d1_user);
    INSERT INTO contacts (mother_profile_id, related_mother_profile_id) VALUES (main_mother_id, d1_mother);


    -- DUMMY BUNDA 2
    INSERT INTO users (auth0_id, email, full_name, gender, phone_number, photo_url)
    VALUES ('auth0|dummy_bunda_2', 'bundadummy2@nusagizi.com', 'Bunda Dummy Dua', 'female', '081234567802', NULL)
    RETURNING id INTO d2_user;
    
    d2_mother := 'ffac6c87-f091-4c46-ba09-72ace2e7cec4';
    INSERT INTO mother_profiles (id, user_id) VALUES (d2_mother, d2_user);
    INSERT INTO contacts (mother_profile_id, related_mother_profile_id) VALUES (main_mother_id, d2_mother);


    -- DUMMY BUNDA 3
    INSERT INTO users (auth0_id, email, full_name, gender, phone_number, photo_url)
    VALUES ('auth0|dummy_bunda_3', 'bundadummy3@nusagizi.com', 'Bunda Dummy Tiga', 'female', '081234567803', NULL)
    RETURNING id INTO d3_user;
    
    d3_mother := '3e93c2b6-f697-406a-bbc1-f77645c05b1c';
    INSERT INTO mother_profiles (id, user_id) VALUES (d3_mother, d3_user);
    INSERT INTO contacts (mother_profile_id, related_mother_profile_id) VALUES (main_mother_id, d3_mother);

END $$;

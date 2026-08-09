INSERT INTO notifications (id, user_id, title, message, notification_type, created_at, updated_at) 
VALUES 
(
    gen_random_uuid(), 
    '5e95e83e-ca95-45aa-83fb-dbe9783b4941', 
    'Waktunya Makan Siang!', 
    'Jangan lupa untuk menyiapkan dan mencatat menu makan siang bergizi untuk si kecil hari ini.', 
    'meal_reminder', 
    now() - interval '2 hours',
    now() - interval '2 hours'
),
(
    gen_random_uuid(), 
    '5e95e83e-ca95-45aa-83fb-dbe9783b4941', 
    'Pencapaian Luar Biasa!', 
    'Selamat! Anda telah berhasil mencatatkan rekor (streak) memantau nutrisi selama 5 hari berturut-turut.', 
    'achievement', 
    now() - interval '1 day',
    now() - interval '1 day'
),
(
    gen_random_uuid(), 
    '5e95e83e-ca95-45aa-83fb-dbe9783b4941', 
    'Pengingat Evaluasi Tumbuh Kembang', 
    'Sudah waktunya mengisi asesmen KPSP bulan ini untuk memastikan perkembangan anak sesuai usia.', 
    'growth_development_reminder', 
    now() - interval '3 days',
    now() - interval '3 days'
),
(
    gen_random_uuid(), 
    '5e95e83e-ca95-45aa-83fb-dbe9783b4941', 
    'Catatan Dokter Baru', 
    'Resep dan ringkasan pemeriksaan dari dr. Spesialis Anak telah berhasil ditambahkan.', 
    'doctor_activity', 
    now() - interval '4 days',
    now() - interval '4 days'
),
(
    gen_random_uuid(), 
    '5e95e83e-ca95-45aa-83fb-dbe9783b4941', 
    'Menu Baru Tersedia', 
    'Ada inspirasi menu baru untuk minggu ini! Cek menu rekomendasi AI yang disesuaikan dengan kondisi anak.', 
    'new_menu_reminder', 
    now() - interval '1 week',
    now() - interval '1 week'
);

INSERT INTO
    children (
        id,
        mother_profile_id,
        full_name,
        gender,
        birth_date,
        photo_url,
        streak_days,
        created_at,
        updated_at
    )
VALUES
    (
        '0c4b4722-a02d-46d4-8d25-d420ba2b7f89',
        '8bdee998-ad6e-46ba-bf57-2e69afaeef78',
        'M. Razky',
        'male',
        '2024-05-07',
        'https://images.unsplash.com/photo-1519238263530-99bdd11df2ea?q=80&w=601',
        2,
        now (),
        now ()
    ),
    (
        '61fa1c7e-81f9-4960-b226-af43f57e8b62',
        '8bdee998-ad6e-46ba-bf57-2e69afaeef78',
        'Budi',
        'male',
        '2025-03-07',
        'https://images.unsplash.com/photo-1620436226596-46a359a295bc?w=500',
        1,
        now (),
        now ()
    ),
    (
        '88e88505-c619-4e95-9e7a-7749d9539468',
        '8bdee998-ad6e-46ba-bf57-2e69afaeef78',
        'Ani',
        'female',
        '2023-07-07',
        'https://images.unsplash.com/photo-1503454537195-1dcabb73ffb9?auto=format&fit=crop&q=80&w=500',
        3,
        now (),
        now ()
    )
ON CONFLICT (id) DO UPDATE SET
    full_name = EXCLUDED.full_name,
    gender = EXCLUDED.gender,
    birth_date = EXCLUDED.birth_date,
    photo_url = EXCLUDED.photo_url,
    streak_days = EXCLUDED.streak_days,
    updated_at = now();
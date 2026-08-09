-- Dummy seed: child_nutrition_reports + turunannya, 1-9 Agustus 2026 (M. Razky, Budi, Ani)
-- Helper unit: "1.6 sdm (16g)", "8.2g", "1/2 sdt (2.5g)", "1 siung (5g)", "sejumput (0.5g)"
CREATE OR REPLACE FUNCTION pg_temp.rand_unit(min_g NUMERIC, max_g NUMERIC, is_spice BOOLEAN)
RETURNS TEXT AS $$
DECLARE
    m_names TEXT[] := ARRAY['sdm','potong','genggam','batang']; 
    m_factors NUMERIC[] := ARRAY[10,15,10,10];
    c_names TEXT[] := ARRAY['sdt','siung','ruas'];
    f_labels TEXT[] := ARRAY['1/4','1/2','3/4','1','1 1/2','2']; 
    f_values NUMERIC[] := ARRAY[0.25,0.5,0.75,1,1.5,2];
    r NUMERIC := random(); 
    style INT; 
    i INT; 
    fi INT; 
    g NUMERIC; 
    qty NUMERIC;
BEGIN
    IF is_spice THEN
        style := CASE 
            WHEN r < 0.2 THEN 3 
            WHEN r < 0.6 THEN 2 
            WHEN r < 0.85 THEN 1 
            ELSE 4 
        END;
    ELSE
        style := CASE 
            WHEN r < 0.35 THEN 1 
            ELSE 4 
        END;
    END IF;

    IF style = 3 THEN
        g := round((0.3 + random() * 0.5)::numeric, 1);
        RETURN 'sejumput (' || g || 'g)';
    ELSIF style = 2 THEN
        i := 1 + floor(random() * array_length(c_names, 1))::int; 
        fi := 1 + floor(random() * array_length(f_values, 1))::int;
        g := round((5 * f_values[fi])::numeric, 1);
        RETURN f_labels[fi] || ' ' || c_names[i] || ' (' || g || 'g)';
    ELSIF style = 4 THEN
        i := 1 + floor(random() * array_length(m_names, 1))::int; 
        qty := round((0.2 + random() * 2.3)::numeric, 1);
        g := round((qty * m_factors[i])::numeric, 1);
        RETURN qty || ' ' || m_names[i] || ' (' || g || 'g)';
    ELSE
        g := round((min_g + random() * (max_g - min_g))::numeric, 1);
        RETURN g || 'g';
    END IF;
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    mother_id UUID := 'da16603d-129d-49be-a311-7eb868b15984';
    d1 DATE := '2026-08-01'; 
    d2 DATE := '2026-08-09';
    meals TEXT[] := ARRAY[
        'breakfast:0.25:Bubur',
        'morning_snack:0.10:Camilan',
        'lunch:0.30:Nasi Tim',
        'afternoon_snack:0.10:Camilan Sore',
        'dinner:0.25:Sup'
    ];
    karbo TEXT[] := ARRAY[
        'AR001:Beras giling', 'BR013:Kentang', 'BR031:Ubi jalar putih',
        'BR028:Ubi jalar kuning', 'BR016:Singkong', 'AR005:Beras jagung kuning'
    ];
    protein TEXT[] := ARRAY[
        'FR005:Daging Ayam', 'GR046:Ikan mas', 'GR053:Ikan patin', 'HR002:Telur ayam ras',
        'FR024:Daging Sapi', 'GR048:Ikan mujahir', 'GR025:Ikan gabus', 'CP077:Tempe', 
        'CP061:Tahu', 'CR026:Kacang merah'
    ];
    sayur TEXT[] := ARRAY[
        'DR008:Bayam', 'DR166:Wortel', 'DR161:Tomat merah', 'DR097:Kacang panjang', 
        'DR013:Buncis', 'DR123:Labu siam', 'DR141:Sawi', 'DR044:Daun kubis', 
        'DR038:Daun kelor', 'DR113:Kool kembang'
    ];
    buah TEXT[] := ARRAY[
        'ER073:Pepaya', 'ER054:Mangga', 'ER087:Pisang mas', 'ER039:Jeruk manis', 
        'ER001:Alpukat', 'ER070:Nanas', 'ER097:Rambutan'
    ];
    lemak TEXT[] := ARRAY['KP003:Santan', 'JR006:Susu sapi', 'KR011:Minyak kelapa'];
    spices TEXT[] := ARRAY[
        'Bawang merah', 'Bawang putih', 'Garam', 'Gula pasir', 'Merica', 
        'Daun salam', 'Serai', 'Kunyit', 'Jahe', 'Ketumbar'
    ];
    steps TEXT[] := ARRAY[
        'Cuci bersih semua bahan.', 'Kukus/rebus bahan utama hingga matang.',
        'Haluskan/potong sesuai tekstur anak.', 'Tumis bumbu halus hingga harum.',
        'Masak dengan api kecil hingga matang.', 'Sajikan hangat sesuai tekstur anak.'
    ];
    
    c RECORD; 
    dt DATE; 
    rid UUID; 
    mid UUID; 
    recid UUID;
    tcal INT; 
    tprot INT; 
    tfat INT; 
    tcarb INT; 
    tex TEXT;
    i INT; 
    j INT; 
    k INT; 
    mt TEXT; 
    frac NUMERIC; 
    label TEXT; 
    snack BOOLEAN;
    cal INT; 
    prot INT; 
    carb INT; 
    fat INT;
    a TEXT; 
    b TEXT; 
    s TEXT; 
    nm TEXT; 
    ds TEXT; 
    port NUMERIC;
BEGIN
    FOR c IN
        SELECT id, full_name,
            CASE full_name WHEN 'M. Razky' THEN 1125 WHEN 'Budi' THEN 1000 ELSE 1250 END AS t1,
            CASE full_name WHEN 'M. Razky' THEN 23 WHEN 'Budi' THEN 20 ELSE 25 END AS t2,
            CASE full_name WHEN 'M. Razky' THEN 40 WHEN 'Budi' THEN 35 ELSE 45 END AS t3,
            CASE full_name WHEN 'M. Razky' THEN 170 WHEN 'Budi' THEN 150 ELSE 190 END AS t4,
            CASE full_name WHEN 'M. Razky' THEN 'potongan kecil' WHEN 'Budi' THEN 'cincang halus' ELSE 'tekstur normal' END AS t5
        FROM children WHERE mother_profile_id = mother_id
    LOOP
        tcal := c.t1; 
        tprot := c.t2; 
        tfat := c.t3; 
        tcarb := c.t4; 
        tex := c.t5;
        dt := d1;
        
        WHILE dt <= d2 LOOP
            rid := gen_random_uuid(); 
            mid := gen_random_uuid();
            
            INSERT INTO child_nutrition_reports (
                id, child_id, calories, target_calories, protein, target_protein, 
                fat, target_fat, carbohydrate, target_carbohydrate, created_at, updated_at
            ) VALUES (
                rid, c.id, round(tcal * (0.6 + random() * 0.45)), tcal, 
                round(tprot * (0.6 + random() * 0.45)), tprot,
                round(tfat * (0.6 + random() * 0.45)), tfat, 
                round(tcarb * (0.6 + random() * 0.45)), tcarb, 
                dt + TIME '08:00', dt + TIME '08:00'
            );
            
            INSERT INTO daily_menus (
                id, child_nutrition_report_id, created_at, updated_at
            ) VALUES (
                mid, rid, dt + TIME '08:00', dt + TIME '08:00'
            );

            FOR i IN 1..5 LOOP
                mt := split_part(meals[i], ':', 1); 
                frac := split_part(meals[i], ':', 2)::numeric; 
                label := split_part(meals[i], ':', 3);
                snack := mt IN ('morning_snack', 'afternoon_snack');
                
                cal := round(tcal * frac * (0.85 + random() * 0.3)); 
                prot := round(tprot * frac * (0.85 + random() * 0.3));
                carb := round(tcarb * frac * (0.85 + random() * 0.3)); 
                fat := round(tfat * frac * (0.85 + random() * 0.3));
                recid := gen_random_uuid();

                IF snack THEN
                    a := buah[1 + floor(random() * array_length(buah, 1))::int]; 
                    b := lemak[1 + floor(random() * array_length(lemak, 1))::int];
                    nm := label || ' ' || split_part(a, ':', 2);
                    ds := 'Camilan sehat berbahan ' || lower(split_part(a, ':', 2)) || '.';
                ELSE
                    a := karbo[1 + floor(random() * array_length(karbo, 1))::int]; 
                    b := protein[1 + floor(random() * array_length(protein, 1))::int];
                    s := sayur[1 + floor(random() * array_length(sayur, 1))::int];
                    nm := label || ' ' || split_part(a, ':', 2) || ' ' || split_part(b, ':', 2);
                    ds := 'Menu ' || replace(mt, '_', ' ') || ' dengan ' || lower(split_part(a, ':', 2)) || 
                          ', ' || lower(split_part(b, ':', 2)) || ', dan ' || lower(split_part(s, ':', 2)) || '.';
                END IF;

                port := (ARRAY[0, 0.33, 0.66, 1])[1 + floor(random() * 4)::int];
                
                INSERT INTO recipes (
                    id, daily_menu_id, name, meal_time, meal_texture, calories, 
                    protein, carbohydrate, fat, description, cooking_time, is_bookmarked, portions_consumed
                ) VALUES (
                    recid, mid, nm, mt::meal_time_type, tex, cal, prot, carb, fat, ds,
                    (10 + floor(random() * 4) * 10) || ' menit', random() < 0.2, port
                );

                INSERT INTO main_ingredients (
                    id, recipe_id, ingredient_id, unit, priority, slot
                ) VALUES (
                    gen_random_uuid(), recid, split_part(a, ':', 1), pg_temp.rand_unit(20, 100, false), 1,
                    CASE WHEN snack THEN 'buah' ELSE 'karbohidrat' END
                );
                
                INSERT INTO main_ingredients (
                    id, recipe_id, ingredient_id, unit, priority, slot
                ) VALUES (
                    gen_random_uuid(), recid, split_part(b, ':', 1), pg_temp.rand_unit(10, 50, false), 2,
                    CASE WHEN snack THEN 'lemak_susu' ELSE 'protein' END
                );
                
                IF NOT snack THEN
                    INSERT INTO main_ingredients (
                        id, recipe_id, ingredient_id, unit, priority, slot
                    ) VALUES (
                        gen_random_uuid(), recid, split_part(s, ':', 1), pg_temp.rand_unit(15, 60, false), 3, 'sayur'
                    );
                END IF;

                FOR j IN 1..(2 + floor(random() * 3)::int) LOOP
                    INSERT INTO cooking_steps (
                        id, recipe_id, step_number, instruction, image_url
                    ) VALUES (
                        gen_random_uuid(), recid, j, steps[1 + floor(random() * array_length(steps, 1))::int],
                        'https://placehold.co/600x400?text=Langkah+' || j
                    );
                END LOOP;

                FOR k IN 1..(1 + floor(random() * 3)::int) LOOP
                    j := 1 + floor(random() * array_length(spices, 1))::int;
                    INSERT INTO recipe_spices (
                        id, recipe_id, name, unit
                    ) VALUES (
                        gen_random_uuid(), recid, spices[j], pg_temp.rand_unit(1, 8, true)
                    );
                END LOOP;
            END LOOP;
            
            dt := dt + 1;
        END LOOP;
    END LOOP;
END $$;
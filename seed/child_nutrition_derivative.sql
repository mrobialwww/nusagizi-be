-- Dummy seed: child_nutrition_reports + turunannya, 1-13 Agustus 2026 (M. Razky, Budi, Ani)
-- Deterministic unit generator: each ingredient code has its own consistent measurement unit (e.g. AR001 is always sdm)
CREATE OR REPLACE FUNCTION pg_temp.fmt_unit(spec TEXT)
RETURNS TEXT AS $$
DECLARE
    u_name TEXT := split_part(spec, ':', 3);
    factor NUMERIC := split_part(spec, ':', 4)::numeric;
    qty NUMERIC;
    g NUMERIC;
BEGIN
    IF u_name = 'sejumput' THEN
        g := round((0.3 + random() * 0.5)::numeric, 1);
        RETURN 'sejumput (' || g || 'g)';
    END IF;

    -- Generate realistic multiplier between 0.4 and 2.0
    qty := round((0.4 + random() * 1.6)::numeric, 1);
    IF factor <= 0 THEN
        factor := 10;
    END IF;
    g := round((qty * factor)::numeric, 1);
    
    RETURN qty || ' ' || u_name || ' (' || g || 'g)';
END;
$$ LANGUAGE plpgsql;

DO $$
DECLARE
    mother_id UUID := 'f545662d-feca-453d-bfa2-a19aab224a16';
    d1 DATE := '2026-08-01'; 
    d2 DATE := '2026-08-20';
    meals TEXT[] := ARRAY[
        'breakfast:0.25:Bubur',
        'morning_snack:0.10:Camilan',
        'lunch:0.30:Nasi Tim',
        'afternoon_snack:0.10:Camilan Sore',
        'dinner:0.25:Sup'
    ];
    karbo TEXT[] := ARRAY[
        'AR001:Beras giling:sdm:10', 'BR013:Kentang:potong:50', 'BR031:Ubi jalar putih:potong:40',
        'BR028:Ubi jalar kuning:potong:40', 'BR016:Singkong:potong:45', 'AR005:Beras jagung kuning:sdm:10'
    ];
    protein TEXT[] := ARRAY[
        'FR005:Daging Ayam:potong:25', 'GR046:Ikan mas:potong:30', 'GR053:Ikan patin:potong:30', 'HR002:Telur ayam ras:butir:50',
        'FR024:Daging Sapi:potong:30', 'GR048:Ikan mujahir:potong:30', 'GR025:Ikan gabus:potong:30', 'CP077:Tempe:potong:20', 
        'CP061:Tahu:potong:25', 'CR026:Kacang merah:sdm:12'
    ];
    sayur TEXT[] := ARRAY[
        'DR008:Bayam:genggam:15', 'DR166:Wortel:batang:25', 'DR161:Tomat merah:potong:20', 'DR097:Kacang panjang:batang:15', 
        'DR013:Buncis:batang:10', 'DR123:Labu siam:potong:30', 'DR141:Sawi:genggam:15', 'DR044:Daun kubis:lembar:15', 
        'DR038:Daun kelor:genggam:10', 'DR113:Kool kembang:potong:25'
    ];
    buah TEXT[] := ARRAY[
        'ER073:Pepaya:potong:50', 'ER054:Mangga:potong:40', 'ER087:Pisang mas:buah:40', 'ER039:Jeruk manis:buah:50', 
        'ER001:Alpukat:potong:45', 'ER070:Nanas:potong:40', 'ER097:Rambutan:buah:20'
    ];
    lemak TEXT[] := ARRAY['KP003:Santan:sdm:10', 'JR006:Susu sapi:sdm:15', 'KR011:Minyak kelapa:sdm:5'];
    spices TEXT[] := ARRAY[
        'Bawang merah:Bawang merah:siung:5', 'Bawang putih:Bawang putih:siung:5', 'Garam:Garam:sdt:2', 
        'Gula pasir:Gula pasir:sdt:4', 'Merica:Merica:sejumput:0.5', 'Daun salam:Daun salam:lembar:1', 
        'Serai:Serai:batang:10', 'Kunyit:Kunyit:ruas:5', 'Jahe:Jahe:ruas:5', 'Ketumbar:Ketumbar:sdt:2'
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
    tcal NUMERIC(6,2); 
    tprot NUMERIC(6,2); 
    tfat NUMERIC(6,2); 
    tcarb NUMERIC(6,2); 
    tex TEXT;
    i INT; 
    j INT; 
    k INT; 
    p INT;
    mt TEXT; 
    frac NUMERIC; 
    label TEXT; 
    snack BOOLEAN;
    cal NUMERIC(6,2); 
    prot NUMERIC(6,2); 
    carb NUMERIC(6,2); 
    fat NUMERIC(6,2);
    a TEXT; 
    b TEXT; 
    s TEXT; 
    nm TEXT; 
    ds TEXT; 
    port NUMERIC;
    buah_chosen TEXT[];
    lemak_chosen TEXT[];
    karbo_chosen TEXT[];
    protein_chosen TEXT[];
    sayur_chosen TEXT[];
    img_url TEXT;
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
                id, child_id, report_date, calories, target_calories, protein, target_protein, fat, target_fat, carbohydrate, target_carbohydrate, created_at, updated_at
            ) VALUES (
                rid, c.id, dt,
                round((tcal * (0.6 + random() * 0.45))::numeric, 2), tcal,
                round((tprot * (0.6 + random() * 0.45))::numeric, 2), tprot,
                round((tfat * (0.6 + random() * 0.45))::numeric, 2), tfat,
                round((tcarb * (0.6 + random() * 0.45))::numeric, 2), tcarb,
                dt + TIME '08:00', dt + TIME '08:00'
            );

            FOR i IN 1..5 LOOP
                mt := split_part(meals[i], ':', 1); 
                frac := split_part(meals[i], ':', 2)::numeric; 
                label := split_part(meals[i], ':', 3);
                snack := mt IN ('morning_snack', 'afternoon_snack');
                
                cal := round((tcal * frac * (0.85 + random() * 0.3))::numeric, 2); 
                prot := round((tprot * frac * (0.85 + random() * 0.3))::numeric, 2);
                carb := round((tcarb * frac * (0.85 + random() * 0.3))::numeric, 2); 
                fat := round((tfat * frac * (0.85 + random() * 0.3))::numeric, 2);
                recid := gen_random_uuid();

                IF snack THEN
                    SELECT array_agg(item) INTO buah_chosen FROM (SELECT unnest(buah) AS item ORDER BY random() LIMIT 3) t;
                    SELECT array_agg(item) INTO lemak_chosen FROM (SELECT unnest(lemak) AS item ORDER BY random() LIMIT 3) t;
                    a := buah_chosen[1]; 
                    b := lemak_chosen[1];
                    nm := label || ' ' || split_part(a, ':', 2);
                    ds := 'Camilan sehat berbahan ' || lower(split_part(a, ':', 2)) || '.';
                ELSE
                    SELECT array_agg(item) INTO karbo_chosen FROM (SELECT unnest(karbo) AS item ORDER BY random() LIMIT 3) t;
                    SELECT array_agg(item) INTO protein_chosen FROM (SELECT unnest(protein) AS item ORDER BY random() LIMIT 3) t;
                    SELECT array_agg(item) INTO sayur_chosen FROM (SELECT unnest(sayur) AS item ORDER BY random() LIMIT 3) t;
                    a := karbo_chosen[1]; 
                    b := protein_chosen[1];
                    s := sayur_chosen[1];
                    nm := label || ' ' || split_part(a, ':', 2) || ' ' || split_part(b, ':', 2);
                    ds := 'Menu ' || replace(mt, '_', ' ') || ' dengan ' || lower(split_part(a, ':', 2)) || 
                          ', ' || lower(split_part(b, ':', 2)) || ', dan ' || lower(split_part(s, ':', 2)) || '.';
                END IF;

                port := (ARRAY[0, 0.33, 0.66, 1])[1 + floor(random() * 4)::int];
                
                IF nm ILIKE '%nasi%' THEN
                    img_url := 'recipes/nasi-gurih.png';
                ELSIF nm ILIKE '%camilan%' THEN
                    img_url := 'recipes/bola-bola.png';
                ELSIF nm ILIKE '%sup%' THEN
                    img_url := 'recipes/sup.png';
                ELSE
                    img_url := 'recipes/tumis.png';
                END IF;

                INSERT INTO recipes (
                    id, name, meal_time, meal_texture, calories, 
                    protein, carbohydrate, fat, description, cooking_time, is_bookmarked, image_url,
                    created_at, updated_at
                ) VALUES (
                    recid, nm, mt::meal_time_type, tex, cal, prot, carb, fat, ds,
                    (10 + floor(random() * 4) * 10) || ' menit', random() < 0.2, img_url,
                    dt + TIME '08:00', dt + TIME '08:00'
                );

                INSERT INTO child_nutrition_recipes (
                    child_nutrition_report_id, recipe_id, portions_consumed, created_at, updated_at
                ) VALUES (
                    rid, recid, port, dt + TIME '08:00', dt + TIME '08:00'
                );

                IF snack THEN
                    FOR p IN 1..3 LOOP
                        IF p <= array_length(buah_chosen, 1) THEN
                            INSERT INTO main_ingredients (
                                id, recipe_id, ingredient_id, unit, priority, slot, created_at, updated_at
                            ) VALUES (
                                gen_random_uuid(), recid, split_part(buah_chosen[p], ':', 1), pg_temp.fmt_unit(buah_chosen[p]), p, 'buah', dt + TIME '08:00', dt + TIME '08:00'
                            );
                        END IF;
                        IF p <= array_length(lemak_chosen, 1) THEN
                            INSERT INTO main_ingredients (
                                id, recipe_id, ingredient_id, unit, priority, slot, created_at, updated_at
                            ) VALUES (
                                gen_random_uuid(), recid, split_part(lemak_chosen[p], ':', 1), pg_temp.fmt_unit(lemak_chosen[p]), p, 'lemak_susu', dt + TIME '08:00', dt + TIME '08:00'
                            );
                        END IF;
                    END LOOP;
                ELSE
                    FOR p IN 1..3 LOOP
                        IF p <= array_length(karbo_chosen, 1) THEN
                            INSERT INTO main_ingredients (
                                id, recipe_id, ingredient_id, unit, priority, slot, created_at, updated_at
                            ) VALUES (
                                gen_random_uuid(), recid, split_part(karbo_chosen[p], ':', 1), pg_temp.fmt_unit(karbo_chosen[p]), p, 'karbohidrat', dt + TIME '08:00', dt + TIME '08:00'
                            );
                        END IF;
                        IF p <= array_length(protein_chosen, 1) THEN
                            INSERT INTO main_ingredients (
                                id, recipe_id, ingredient_id, unit, priority, slot, created_at, updated_at
                            ) VALUES (
                                gen_random_uuid(), recid, split_part(protein_chosen[p], ':', 1), pg_temp.fmt_unit(protein_chosen[p]), p, 'protein', dt + TIME '08:00', dt + TIME '08:00'
                            );
                        END IF;
                        IF p <= array_length(sayur_chosen, 1) THEN
                            INSERT INTO main_ingredients (
                                id, recipe_id, ingredient_id, unit, priority, slot, created_at, updated_at
                            ) VALUES (
                                gen_random_uuid(), recid, split_part(sayur_chosen[p], ':', 1), pg_temp.fmt_unit(sayur_chosen[p]), p, 'sayur', dt + TIME '08:00', dt + TIME '08:00'
                            );
                        END IF;
                    END LOOP;
                END IF;

                FOR j IN 1..(2 + floor(random() * 3)::int) LOOP
                    INSERT INTO cooking_steps (
                        id, recipe_id, step_number, instruction, image_url, created_at, updated_at
                    ) VALUES (
                        gen_random_uuid(), recid, j, steps[1 + floor(random() * array_length(steps, 1))::int],
                        'https://placehold.co/600x400?text=Langkah+' || j, dt + TIME '08:00', dt + TIME '08:00'
                    );
                END LOOP;

                FOR k IN 1..(1 + floor(random() * 3)::int) LOOP
                    j := 1 + floor(random() * array_length(spices, 1))::int;
                    INSERT INTO recipe_spices (
                        id, recipe_id, name, unit, created_at, updated_at
                    ) VALUES (
                        gen_random_uuid(), recid, split_part(spices[j], ':', 2), pg_temp.fmt_unit(spices[j]), dt + TIME '08:00', dt + TIME '08:00'
                    );
                END LOOP;
            END LOOP;
            
            dt := dt + 1;
        END LOOP;
    END LOOP;
END $$;
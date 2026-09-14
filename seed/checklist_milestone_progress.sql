-- =====================================================================
-- SEED DATA: checklist_milestone_progress
-- Variatif per anak, per bulan, per domain — mirip data real-case
-- =====================================================================

WITH
-- =============================================
-- Profil target: berapa pertanyaan yang di-check
-- per child, per month_target (dari 10 total)
-- =============================================
child_month_target (child_id, month_target, target_checked) AS (
    VALUES
    -- M. Razky (lahir 2024-05-07, ~27 bln)
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid,  3,  8),
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid,  6,  7),
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid,  9,  6),
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, 12,  8),
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, 15,  7),
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, 18,  5),
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, 21,  7),
    ('0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, 24,  6),

    -- Budi (lahir 2025-03-07, ~17 bln)
    ('61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid,  3,  9),
    ('61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid,  6,  8),
    ('61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid,  9,  7),
    ('61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid, 12,  6),
    ('61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid, 15,  8),

    -- Ani (lahir 2023-07-07, ~37 bln)
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid,  3, 10),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid,  6,  9),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid,  9,  8),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid, 12,  7),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid, 15,  8),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid, 18,  6),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid, 21,  7),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid, 24,  8),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid, 30,  7),
    ('88e88505-c619-4e95-9e7a-7749d9539468'::uuid, 36,  6)
),
-- =============================================
-- Ambil semua pertanyaan yang relevan
-- =============================================
questions AS (
    SELECT
        q.id AS question_id,
        q.month_target,
        q.developmental_domain,
        q.question_text
    FROM assessment_kpsp_questions q
    WHERE q.month_target IN (3,6,9,12,15,18,21,24,30,36)
),
-- =============================================
-- Distribusi checked per domain:
-- gross_motor: 3 soal, fine_motor: 3 soal
-- speech: 2 soal, social: 2 soal
-- =============================================
domain_check AS (
    -- Razky: gross_motor kuat, speech lemah di bulan 18
    SELECT
        cmt.child_id,
        cmt.month_target,
        CASE
            WHEN cmt.month_target IN (3,6,12,15,21) THEN 3
            WHEN cmt.month_target = 9 THEN 2
            WHEN cmt.month_target = 18 THEN 1
            ELSE 2
        END AS gross_checked,
        CASE
            WHEN cmt.month_target IN (3,12,21) THEN 3
            WHEN cmt.month_target IN (6,15,24) THEN 2
            ELSE 1
        END AS fine_checked,
        CASE
            WHEN cmt.month_target IN (3,12) THEN 2
            WHEN cmt.month_target IN (6,15,21) THEN 1
            ELSE 0
        END AS speech_checked,
        CASE
            WHEN cmt.month_target IN (3,6,9,12,15,18,21,24) THEN
                CASE cmt.target_checked
                    WHEN 8 THEN 1
                    WHEN 7 THEN 1
                    WHEN 6 THEN 2
                    WHEN 5 THEN 2
                    ELSE 1
                END
            ELSE 1
        END AS social_checked
    FROM child_month_target cmt
    WHERE cmt.child_id = '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'

    UNION ALL

    -- Budi: semua domain merata, speech lebih kuat
    SELECT
        cmt.child_id,
        cmt.month_target,
        CASE
            WHEN cmt.month_target IN (3,6,9) THEN 3
            ELSE 2
        END AS gross_checked,
        CASE
            WHEN cmt.month_target IN (3,15) THEN 3
            ELSE 2
        END AS fine_checked,
        CASE
            WHEN cmt.month_target IN (3,6,12,15) THEN 2
            ELSE 1
        END AS speech_checked,
        CASE
            WHEN cmt.month_target = 3 THEN 2
            WHEN cmt.month_target IN (6,9) THEN 1
            WHEN cmt.month_target = 12 THEN 1
            ELSE 3
        END AS social_checked
    FROM child_month_target cmt
    WHERE cmt.child_id = '61fa1c7e-81f9-4960-b226-af43f57e8b62'

    UNION ALL

    -- Ani: gross_motor konsisten kuat, fine_motor bervariasi
    SELECT
        cmt.child_id,
        cmt.month_target,
        CASE
            WHEN cmt.month_target <= 24 THEN 3
            WHEN cmt.month_target = 30 THEN 2
            ELSE 1
        END AS gross_checked,
        CASE
            WHEN cmt.month_target IN (3,6,9) THEN 3
            WHEN cmt.month_target IN (12,15) THEN 2
            WHEN cmt.month_target IN (18,21) THEN 1
            WHEN cmt.month_target = 24 THEN 2
            WHEN cmt.month_target = 30 THEN 3
            ELSE 1
        END AS fine_checked,
        CASE
            WHEN cmt.month_target IN (3,6) THEN 2
            WHEN cmt.month_target IN (9,12) THEN 1
            WHEN cmt.month_target IN (15,18) THEN 2
            WHEN cmt.month_target IN (21,24) THEN 1
            WHEN cmt.month_target = 30 THEN 2
            ELSE 1
        END AS speech_checked,
        CASE
            WHEN cmt.month_target IN (3,6,9) THEN 2
            WHEN cmt.month_target IN (12,15) THEN 2
            WHEN cmt.month_target IN (18,21) THEN 2
            WHEN cmt.month_target IN (24,30) THEN 1
            ELSE 2
        END AS social_checked
    FROM child_month_target cmt
    WHERE cmt.child_id = '88e88505-c619-4e95-9e7a-7749d9539468'
),
-- =============================================
-- Assign urutan per domain untuk menentukan
-- soal mana yang di-check (md5 untuk variasi)
-- =============================================
ranked_questions AS (
    SELECT
        q.question_id,
        q.month_target,
        q.developmental_domain,
        ROW_NUMBER() OVER (
            PARTITION BY q.month_target, q.developmental_domain
            ORDER BY md5(q.question_id::text)
        ) AS domain_row_num
    FROM questions q
),
-- =============================================
-- Pilih soal berdasarkan jumlah checked per domain
-- =============================================
selected_checks AS (
    SELECT
        dc.child_id,
        rq.question_id,
        rq.month_target,
        rq.developmental_domain
    FROM domain_check dc
    JOIN ranked_questions rq ON rq.month_target = dc.month_target
    WHERE
        (rq.developmental_domain = 'gross_motor_skills'
            AND rq.domain_row_num <= dc.gross_checked)
        OR
        (rq.developmental_domain = 'fine_motor_skills'
            AND rq.domain_row_num <= dc.fine_checked)
        OR
        (rq.developmental_domain = 'speech_and_language'
            AND rq.domain_row_num <= dc.speech_checked)
        OR
        (rq.developmental_domain = 'socialization'
            AND rq.domain_row_num <= dc.social_checked)
),
-- =============================================
-- Timestamp historis: milestone dicentang
-- saat usia anak mencapai bulan tersebut
-- =============================================
child_birth AS (
    SELECT id AS child_id, birth_date FROM children
)
INSERT INTO checklist_milestone_progress (child_id, assessment_kpsp_question_id, created_at, updated_at)
SELECT
    sc.child_id,
    sc.question_id,
    cb.birth_date + (sc.month_target || ' months')::interval,
    now()
FROM selected_checks sc
JOIN child_birth cb ON cb.child_id = sc.child_id
ON CONFLICT (child_id, assessment_kpsp_question_id) DO UPDATE SET
    updated_at = EXCLUDED.updated_at;

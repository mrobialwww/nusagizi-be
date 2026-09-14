-- =====================================================================
-- SEED DATA: child_development_reports, assessment_kpsp_answers,
--            development_report_recommendations
-- =====================================================================

WITH report_data (id, child_id, child_birth_date, month_target, next_milestone, target_score) AS (
    VALUES
    -- =======================================
    -- M. Razky (0c4b4722-a02d-46d4-8d25-d420ba2b7f89, lahir 2024-05-07)
    -- Sekarang ~27 bulan. Milestones: 3, 6, 9, 12, 15, 18, 21, 24
    -- 10 pertanyaan/bulan: 7 TRUE, 3 FALSE = 30% FALSE (~1/3)
    -- =======================================
    ('d4304859-9941-47c3-8fef-9a5c0becc1e0'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 3, 6, 7),
    ('39129532-6aed-4f27-90d2-bf56f2e21b0e'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 6, 9, 7),
    ('8cc29dca-3c58-450f-a492-f0459c5d013f'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 9, 12, 7),
    ('a4cdbf78-1a8c-4f7f-afbd-028a05f9919f'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 12, 15, 7),
    ('e485e927-4406-444a-bee6-d71e2ef64d1f'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 15, 18, 7),
    ('77f6b9bc-4b35-4dbb-a81e-128a1ea2221b'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 18, 21, 7),
    ('1eb0579e-4a6c-48be-bea1-6cbcd324eedd'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 21, 24, 7),
    ('54e60153-f726-47da-993d-d1be8606c134'::uuid, '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'::uuid, '2024-05-07'::date, 24, 30, 7),

    -- =======================================
    -- Budi (61fa1c7e-81f9-4960-b226-af43f57e8b62, lahir 2025-03-07)
    -- Sekarang ~17 bulan. Milestones: 3, 6, 9, 12, 15
    -- 10 pertanyaan/bulan: 7 TRUE, 3 FALSE = 30% FALSE (~1/3)
    -- =======================================
    ('b81781cb-63df-40cb-baec-4d51b7aff734'::uuid, '61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid, '2025-03-07'::date, 3, 6, 7),
    ('18534ec5-feab-4246-abe2-1fe701b22e70'::uuid, '61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid, '2025-03-07'::date, 6, 9, 7),
    ('f366113b-aecd-47bc-8cd3-492e59dfdcd7'::uuid, '61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid, '2025-03-07'::date, 9, 12, 7),
    ('65c929fa-33cc-4e89-8d89-9a2cf110a19e'::uuid, '61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid, '2025-03-07'::date, 12, 15, 7),
    ('4f697472-e1c9-4b67-ba30-d380b27b6c50'::uuid, '61fa1c7e-81f9-4960-b226-af43f57e8b62'::uuid, '2025-03-07'::date, 15, 18, 7),

    -- =======================================
    -- Ani (88e88505-c619-4e95-9e7a-7749d9539468, lahir 2023-07-07)
    -- Sekarang ~37 bulan. Milestones: 3, 6, 9, 12, 15, 18, 21, 24, 30, 36
    -- 10 pertanyaan/bulan: 7 TRUE, 3 FALSE = 30% FALSE (~1/3)
    -- =======================================
    ('c00ad8fe-ad04-45e0-afd9-3796bba46cda'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 3, 6, 7),
    ('b6f4e668-5f21-4f11-9a70-80a55255470d'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 6, 9, 7),
    ('538e1215-dcfa-4f21-ada6-baff2aed3c8c'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 9, 12, 7),
    ('a9de4f2c-e5b1-4fec-be10-096eb255b550'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 12, 15, 7),
    ('ed526017-d54e-4f0f-8c3b-74b8e0108341'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 15, 18, 7),
    ('1a8b1115-3eb7-40d1-93c6-e916ea861a7a'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 18, 21, 7),
    ('109ffb73-0eb4-4ddf-ba57-8ab05e830e9d'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 21, 24, 7),
    ('62c3e1b7-7e6d-4952-b883-75a6c3f3abdd'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 24, 30, 7),
    ('a8b8d4f4-500b-4f4c-83b1-f92ac06fb0c1'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 30, 36, 7),
    ('b4c7c880-9975-4518-a6b1-40be3cd5d487'::uuid, '88e88505-c619-4e95-9e7a-7749d9539468'::uuid, '2023-07-07'::date, 36, 42, 7)
),
inserted_reports AS (
    INSERT INTO child_development_reports (
        id, child_id, kpsp_score, month_target, next_check_date, created_at, updated_at
    )
    SELECT
        id,
        child_id,
        target_score,
        month_target,
        CASE
            WHEN next_milestone IS NOT NULL
            THEN child_birth_date + (next_milestone || ' months')::interval
            ELSE NULL
        END,
        child_birth_date + (month_target || ' months')::interval,
        now()
    FROM report_data
    ON CONFLICT (id) DO UPDATE SET
        kpsp_score = EXCLUDED.kpsp_score,
        next_check_date = EXCLUDED.next_check_date,
        updated_at = EXCLUDED.updated_at
    RETURNING id
),
answers_to_insert AS (
    SELECT
        r.id AS report_id,
        q.id AS question_id,
        CASE
            WHEN row_number() OVER (PARTITION BY r.id ORDER BY md5(r.id::text || q.id::text)) <= r.target_score THEN TRUE
            ELSE FALSE
        END AS answer,
        r.child_birth_date + (r.month_target || ' months')::interval AS historical_created_at
    FROM report_data r
    JOIN assessment_kpsp_questions q ON q.month_target = r.month_target
),
inserted_answers AS (
    INSERT INTO assessment_kpsp_answers (
        child_development_report_id, assessment_kpsp_question_id, answer, created_at, updated_at
    )
    SELECT report_id, question_id, answer, historical_created_at, now()
    FROM answers_to_insert
    ON CONFLICT (child_development_report_id, assessment_kpsp_question_id) DO UPDATE SET
        answer = EXCLUDED.answer,
        updated_at = EXCLUDED.updated_at
    RETURNING child_development_report_id, assessment_kpsp_question_id
)
INSERT INTO development_report_recommendations (
    child_development_report_id, recommended_action_id, created_at
)
SELECT
    a.report_id,
    ra.id,
    a.historical_created_at
FROM answers_to_insert a
JOIN recommended_actions ra ON ra.assessment_kpsp_question_id = a.question_id
WHERE a.answer = FALSE
ON CONFLICT (child_development_report_id, recommended_action_id) DO NOTHING;

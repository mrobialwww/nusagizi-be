-- Seed file for caregiver_engagements
-- Caregiver has access to 3 children

INSERT INTO caregiver_engagements (caregiver_profile_id, child_id)
VALUES
    ('eb9cd3a5-18de-4a79-b4ee-69f9ecf87504', '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'),
    ('eb9cd3a5-18de-4a79-b4ee-69f9ecf87504', '61fa1c7e-81f9-4960-b226-af43f57e8b62'),
    ('eb9cd3a5-18de-4a79-b4ee-69f9ecf87504', '88e88505-c619-4e95-9e7a-7749d9539468')
ON CONFLICT DO NOTHING;

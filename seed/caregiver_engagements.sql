-- Seed file for caregiver_engagements
-- Caregiver has access to 3 children

INSERT INTO caregiver_engagements (caregiver_profile_id, child_id)
VALUES
    ('debdd9b0-a7c6-4278-acd3-eb7d10e38537', '0c4b4722-a02d-46d4-8d25-d420ba2b7f89'),
    ('debdd9b0-a7c6-4278-acd3-eb7d10e38537', '61fa1c7e-81f9-4960-b226-af43f57e8b62'),
    ('debdd9b0-a7c6-4278-acd3-eb7d10e38537', '88e88505-c619-4e95-9e7a-7749d9539468')
ON CONFLICT DO NOTHING;

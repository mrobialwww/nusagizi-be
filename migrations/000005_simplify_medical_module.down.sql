-- ---------------------------------------------------------------------
-- revert simplify medical module
-- ---------------------------------------------------------------------

-- Drop the simplified medical notes and its dependencies
DROP TABLE IF EXISTS daily_nutrition_targets;
DROP TABLE IF EXISTS medical_restrictions;
DROP TABLE IF EXISTS medical_notes;

-- Recreate doctor_profiles
CREATE TABLE doctor_profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    hospital_name   VARCHAR(150),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_doctor_profiles_user_id ON doctor_profiles(user_id);
CREATE TRIGGER trg_doctor_profiles_updated_at BEFORE UPDATE ON doctor_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Recreate medical_relationships
CREATE TABLE medical_relationships (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doctor_profile_id   UUID NOT NULL REFERENCES doctor_profiles(id) ON DELETE CASCADE,
    child_id            UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    deleted_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_medical_relationships_doctor_id ON medical_relationships(doctor_profile_id);
CREATE INDEX idx_medical_relationships_child_id ON medical_relationships(child_id);
CREATE TRIGGER trg_medical_relationships_updated_at BEFORE UPDATE ON medical_relationships
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Recreate original medical_notes
CREATE TABLE medical_notes (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_relationship_id     UUID NOT NULL REFERENCES medical_relationships(id) ON DELETE CASCADE,
    recommendation              TEXT NOT NULL,
    valid_date                  DATE NOT NULL,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_medical_notes_relationship_id ON medical_notes(medical_relationship_id);
CREATE TRIGGER trg_medical_notes_updated_at BEFORE UPDATE ON medical_notes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Recreate original medical_restrictions
CREATE TABLE medical_restrictions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_note_id  UUID NOT NULL REFERENCES medical_notes(id) ON DELETE CASCADE,
    type             medical_restriction_type NOT NULL,
    substance_name   TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_medical_restrictions_note_id ON medical_restrictions(medical_note_id);
CREATE TRIGGER trg_medical_restrictions_updated_at BEFORE UPDATE ON medical_restrictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Recreate original daily_nutrition_targets
CREATE TABLE daily_nutrition_targets (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_note_id  UUID NOT NULL REFERENCES medical_notes(id) ON DELETE CASCADE,
    nutrient_name    VARCHAR(100) NOT NULL,
    quantity         NUMERIC(8,2) NOT NULL CHECK (quantity > 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_daily_nutrition_targets_note_id ON daily_nutrition_targets(medical_note_id);
CREATE TRIGGER trg_daily_nutrition_targets_updated_at BEFORE UPDATE ON daily_nutrition_targets
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
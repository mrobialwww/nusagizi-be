-- ---------------------------------------------------------------------
-- simplify medical module (Incremental changes from 000004)
-- ---------------------------------------------------------------------

-- Drop removed tables
DROP TABLE IF EXISTS medical_restrictions;
DROP TABLE IF EXISTS daily_nutrition_targets;
DROP TABLE IF EXISTS medical_notes;
DROP TABLE IF EXISTS medical_relationships CASCADE;
DROP TABLE IF EXISTS doctor_profiles CASCADE;

-- Recreate medical_notes linking to child directly
CREATE TABLE medical_notes (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id                    UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    doctor_name                 VARCHAR(150) NOT NULL,
    facility_name           VARCHAR(150) NOT NULL,
    recommendation              TEXT NOT NULL,
    valid_until                  DATE NOT NULL,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_medical_notes_child_id ON medical_notes(child_id);
CREATE TRIGGER trg_medical_notes_updated_at BEFORE UPDATE ON medical_notes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Recreate medical_restrictions
CREATE TABLE medical_restrictions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_note_id  UUID NOT NULL REFERENCES medical_notes(id) ON DELETE CASCADE,
    type             medical_restriction_type NOT NULL,
    item_name   TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_medical_restrictions_note_id ON medical_restrictions(medical_note_id);
CREATE TRIGGER trg_medical_restrictions_updated_at BEFORE UPDATE ON medical_restrictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Recreate daily_nutrition_targets
CREATE TABLE daily_nutrition_targets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_note_id UUID NOT NULL REFERENCES medical_notes(id) ON DELETE CASCADE,
    nutrient        nutrient_type NOT NULL,
    quantity        NUMERIC(8,2) NOT NULL CHECK (quantity > 0),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_daily_nutrition_targets_note_id ON daily_nutrition_targets(medical_note_id);
CREATE TRIGGER trg_daily_nutrition_targets_updated_at BEFORE UPDATE ON daily_nutrition_targets
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
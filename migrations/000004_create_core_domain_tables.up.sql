
CREATE EXTENSION IF NOT EXISTS pgcrypto;


CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ENUM: gender
CREATE TYPE gender_type AS ENUM ('male', 'female');

-- ENUM: child_allergy_profile.category
CREATE TYPE child_allergy_type AS ENUM ('food', 'medicine', 'animal', 'other');

-- ENUM: medical_notes.nutrient
CREATE TYPE nutrient_type AS ENUM ('calorie', 'protein', 'fat', 'carbohydrate');

-- ENUM: developmental_domain
CREATE TYPE developmental_domain AS ENUM (
    'gross_motor_skills',
    'fine_motor_skills',
    'speech_and_language',
    'socialization'
);

-- ENUM: recipes.meal_time
CREATE TYPE meal_time_type AS ENUM (
    'breakfast',
    'lunch',
    'dinner',
    'morning_snack',
    'afternoon_snack'
);

-- ENUM: child_photos.visibility
CREATE TYPE child_photo_visibility AS ENUM ('all', 'private', 'selected_only');

-- ENUM: medical_restrictions.type
CREATE TYPE medical_restriction_type AS ENUM ('prohibition', 'allergy');

-- ENUM: notification.type
CREATE TYPE notification_type AS ENUM ('meal_reminder', 'growth_development_reminder', 'doctor_activity', 'photo_activity', 'shop_activity', 'invitation_activity', 'new_menu_reminder', 'achievement');

-- users (already created in 000001)
ALTER TABLE users 
    DROP COLUMN role,
    ADD COLUMN full_name VARCHAR(150),
    ADD COLUMN gender gender_type,
    ADD COLUMN phone_number VARCHAR(20),
    ADD COLUMN photo_url TEXT,
    ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- user -> mother_profiles (1:1)
-- Stores the profile information of mothers (1:1 with users)
CREATE TABLE mother_profiles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_mother_profiles_user_id ON mother_profiles(user_id);
CREATE TRIGGER trg_mother_profiles_updated_at BEFORE UPDATE ON mother_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- user -> caregiver_profiles (1:1)
-- Stores the profile information of caregivers (1:1 with users)
CREATE TABLE caregiver_profiles (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_caregiver_profiles_user_id ON caregiver_profiles(user_id);
CREATE TRIGGER trg_caregiver_profiles_updated_at BEFORE UPDATE ON caregiver_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- user -> doctor_profiles (1:1)
-- Stores the profile information of doctors (1:1 with users)
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

-- ---------------------------------------------------------------------
-- create child core
-- ---------------------------------------------------------------------
-- mother_profiles -> children (1:N)
-- Master table for children profiles, linked to mother_profiles (1:N)
CREATE TABLE children (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mother_profile_id       UUID NOT NULL REFERENCES mother_profiles(id) ON DELETE CASCADE,
    full_name               VARCHAR(150) NOT NULL,
    gender                  gender_type NOT NULL,
    birth_date              DATE NOT NULL,
    photo_url               TEXT,
    notes_profile           TEXT,
    streak_days             INTEGER NOT NULL CHECK (streak_days >= 0),
    last_streak_date        DATE,
    deleted_at              TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_children_mother_profile_id ON children(mother_profile_id);
CREATE TRIGGER trg_children_updated_at BEFORE UPDATE ON children
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> favorite_food_profiles (1:N)
-- Stores favorite foods for each child (1:N)
CREATE TABLE favorite_food_profiles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id    UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    food_name   VARCHAR(150) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_favorite_food_profiles_child_id ON favorite_food_profiles(child_id);
CREATE TRIGGER trg_favorite_food_profiles_updated_at BEFORE UPDATE ON favorite_food_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> favorite_texture_profiles (1:N)
-- Stores favorite food textures for each child (1:N)
CREATE TABLE favorite_texture_profiles (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id      UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    texture_name  VARCHAR(100) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_favorite_texture_profiles_child_id ON favorite_texture_profiles(child_id);
CREATE TRIGGER trg_favorite_texture_profiles_updated_at BEFORE UPDATE ON favorite_texture_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> child_diet_profiles (1:N)
-- Stores dietary preferences or profiles for each child (1:N)
CREATE TABLE child_diet_profiles (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id     UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    diet_name    VARCHAR(150) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_diet_profiles_child_id ON child_diet_profiles(child_id);
CREATE TRIGGER trg_child_diet_profiles_updated_at BEFORE UPDATE ON child_diet_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> child_chronic_disease_profiles (1:N)
-- Stores chronic diseases affecting the child (1:N)
CREATE TABLE child_chronic_disease_profiles (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id       UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    disease_name   VARCHAR(150) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_chronic_disease_profiles_child_id ON child_chronic_disease_profiles(child_id);
CREATE TRIGGER trg_child_chronic_disease_profiles_updated_at BEFORE UPDATE ON child_chronic_disease_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> child_allergy_profiles (1:N), enum category
-- Stores allergy profiles including category and allergen name (1:N)
CREATE TABLE child_allergy_profiles (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id       UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    category       child_allergy_type NOT NULL,
    allergen_name  VARCHAR(150) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_allergy_profiles_child_id ON child_allergy_profiles(child_id);
CREATE TRIGGER trg_child_allergy_profiles_updated_at BEFORE UPDATE ON child_allergy_profiles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create child growth
-- ---------------------------------------------------------------------
-- child -> child_growth_reports (1:N)
-- Stores periodic growth measurements (weight, height, head circumference)
CREATE TABLE child_growth_reports (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id                UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    measured_at             DATE NOT NULL,
    weight_kg               NUMERIC(5,2) CHECK (weight_kg > 0),
    height_cm               NUMERIC(5,2) CHECK (height_cm > 0),
    head_circumference_cm   NUMERIC(5,2) CHECK (head_circumference_cm > 0),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_growth_reports_child_id ON child_growth_reports(child_id);
CREATE TRIGGER trg_child_growth_reports_updated_at BEFORE UPDATE ON child_growth_reports
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create child development
-- ---------------------------------------------------------------------
-- child -> child_development_reports (1:N)

-- Stores developmental assessment reports (KPSP scores) and targets
CREATE TABLE child_development_reports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id        UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    kpsp_score      INTEGER NOT NULL,
    month_target    INTEGER NOT NULL CHECK (month_target IN (3,6,9,12,15,18,21,24,30,36,42,48,54,60)),
    next_check_date DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_development_reports_child_id ON child_development_reports(child_id);
CREATE TRIGGER trg_child_development_reports_updated_at BEFORE UPDATE ON child_development_reports
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();


-- Master table for KPSP developmental assessment questions by age and domain
CREATE TABLE assessment_kpsp_questions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    month_target            INTEGER NOT NULL CHECK (month_target IN (3,6,9,12,15,18,21,24,30,36,42,48,54,60)),
    developmental_domain    developmental_domain NOT NULL,
    question_text           TEXT NOT NULL,
    image_url               TEXT,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_assessment_kpsp_questions_month_target ON assessment_kpsp_questions(month_target);
CREATE TRIGGER trg_assessment_kpsp_questions_updated_at BEFORE UPDATE ON assessment_kpsp_questions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child_development_reports -> assessment_kpsp_answers (1:N)
-- assessment_kpsp_questions -> assessment_kpsp_answers (1:N)
-- Stores the child's answers (Yes/No) to specific KPSP questions
CREATE TABLE assessment_kpsp_answers (
    id                            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_development_report_id   UUID NOT NULL REFERENCES child_development_reports(id) ON DELETE CASCADE,
    assessment_kpsp_question_id   UUID NOT NULL REFERENCES assessment_kpsp_questions(id) ON DELETE RESTRICT,
    answer                        BOOLEAN NOT NULL,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (child_development_report_id, assessment_kpsp_question_id)
);
CREATE INDEX idx_assessment_kpsp_answers_report_id ON assessment_kpsp_answers(child_development_report_id);
CREATE INDEX idx_assessment_kpsp_answers_question_id ON assessment_kpsp_answers(assessment_kpsp_question_id);
CREATE TRIGGER trg_assessment_kpsp_answers_updated_at BEFORE UPDATE ON assessment_kpsp_answers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- assessment_kpsp_questions -> recommended_actions (1:1)
-- Stores recommended actions linked to specific KPSP questions
CREATE TABLE recommended_actions (
    id                            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_kpsp_question_id   UUID NOT NULL UNIQUE REFERENCES assessment_kpsp_questions(id) ON DELETE CASCADE,
    title                         VARCHAR(50) NOT NULL,
    action_text                   TEXT NOT NULL,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_recommended_actions_question_id ON recommended_actions(assessment_kpsp_question_id);
CREATE TRIGGER trg_recommended_actions_updated_at BEFORE UPDATE ON recommended_actions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Junction N:M: child_development_reports <-> recommended_actions
-- Junction table linking child development reports to recommended actions
CREATE TABLE development_report_recommendations (
    child_development_report_id   UUID NOT NULL REFERENCES child_development_reports(id) ON DELETE CASCADE,
    recommended_action_id         UUID NOT NULL REFERENCES recommended_actions(id) ON DELETE CASCADE,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_development_report_id, recommended_action_id)
);
CREATE INDEX idx_attempt_recomendations_report_id ON development_report_recommendations(child_development_report_id);
CREATE INDEX idx_attempt_recomendations_action_id ON development_report_recommendations(recommended_action_id);


-- Junction N:M: child <-> assessment_kpsp_questions
-- Junction table tracking a child's progress on specific milestone tasks
CREATE TABLE checklist_milestone_progress (
    child_id                    UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    assessment_kpsp_question_id UUID NOT NULL REFERENCES assessment_kpsp_questions(id) ON DELETE CASCADE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_id, assessment_kpsp_question_id)
);
CREATE INDEX idx_checklist_milestone_progress_child_id ON checklist_milestone_progress(child_id);
CREATE INDEX idx_checklist_milestone_progress_question_id ON checklist_milestone_progress(assessment_kpsp_question_id);
CREATE TRIGGER trg_checklist_milestone_progress_updated_at BEFORE UPDATE ON checklist_milestone_progress
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create child nutrition
-- ---------------------------------------------------------------------
-- child -> child_nutrition_reports (1:N)
-- Stores daily nutrition targets and achieved macros for a child
CREATE TABLE child_nutrition_reports (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id                UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    report_date             DATE NOT NULL,
    calories                NUMERIC(6,2) NOT NULL CHECK (calories >= 0),
    target_calories         NUMERIC(6,2) NOT NULL CHECK (target_calories >= 0),
    protein                 NUMERIC(6,2) NOT NULL CHECK (protein >= 0),
    target_protein          NUMERIC(6,2) NOT NULL CHECK (target_protein >= 0),
    fat                     NUMERIC(6,2) NOT NULL CHECK (fat >= 0),
    target_fat              NUMERIC(6,2) NOT NULL CHECK (target_fat >= 0),
    carbohydrate            NUMERIC(6,2) NOT NULL CHECK (carbohydrate >= 0),
    target_carbohydrate     NUMERIC(6,2) NOT NULL CHECK (target_carbohydrate >= 0),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (child_id, report_date)
);
CREATE INDEX idx_child_nutrition_reports_child_id ON child_nutrition_reports(child_id);
CREATE UNIQUE INDEX unique_child_report_date ON child_nutrition_reports(child_id, report_date);
CREATE TRIGGER trg_child_nutrition_reports_updated_at BEFORE UPDATE ON child_nutrition_reports
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Stores individual meal recipes including macros, cooking time, and meal time
CREATE TABLE recipes (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(150) NOT NULL,
    meal_time           meal_time_type NOT NULL,
    meal_texture        VARCHAR(25) NOT NULL,
    calories            NUMERIC(6,2) NOT NULL CHECK (calories >= 0),
    protein             NUMERIC(6,2) NOT NULL CHECK (protein >= 0),
    carbohydrate        NUMERIC(6,2) NOT NULL CHECK (carbohydrate >= 0),
    fat                 NUMERIC(6,2) NOT NULL CHECK (fat >= 0),
    description         TEXT NOT NULL,
    image_url           TEXT,
    cooking_time        VARCHAR(100) NOT NULL,
    is_bookmarked       BOOLEAN NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TRIGGER trg_recipes_updated_at BEFORE UPDATE ON recipes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();


-- child_nutrition_reports <-> recipes (N:M)
-- Junction table linking daily nutrition reports to recipes (N:M)
CREATE TABLE child_nutrition_recipes (
    child_nutrition_report_id  UUID NOT NULL REFERENCES child_nutrition_reports(id) ON DELETE CASCADE,
    recipe_id                  UUID NOT NULL REFERENCES recipes(id),
    portions_consumed          NUMERIC(3,2) NOT NULL DEFAULT 0,
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_nutrition_report_id, recipe_id)
);
CREATE INDEX idx_child_nutrition_recipes_report_id ON child_nutrition_recipes(child_nutrition_report_id);
CREATE INDEX idx_child_nutrition_recipes_recipe_id ON child_nutrition_recipes(recipe_id);
CREATE TRIGGER trg_child_nutrition_recipes_updated_at BEFORE UPDATE ON child_nutrition_recipes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();


-- Master table for all available food ingredients
CREATE TABLE ingredients (
    id          VARCHAR(100) PRIMARY KEY,
    name        VARCHAR(150) NOT NULL UNIQUE,
    image_url   TEXT,
    category    VARCHAR(100),
    price       VARCHAR(10),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TRIGGER trg_ingredients_updated_at BEFORE UPDATE ON ingredients
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- recipes -> main_ingredients (1:N), ingredients -> main_ingredients (1:N)
-- Junction table linking recipes to ingredients, specifying unit and priority
CREATE TABLE main_ingredients (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id      UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id  VARCHAR(5) NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    unit           VARCHAR(50) NOT NULL,
    priority       INTEGER NOT NULL CHECK (priority > 0),
    slot           VARCHAR(25) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (recipe_id, ingredient_id)
);
CREATE INDEX idx_main_ingredients_recipe_id ON main_ingredients(recipe_id);
CREATE INDEX idx_main_ingredients_ingredient_id ON main_ingredients(ingredient_id);
CREATE TRIGGER trg_main_ingredients_updated_at BEFORE UPDATE ON main_ingredients
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- recipes -> cooking_steps (1:N)
-- Stores sequential cooking instructions for recipes
CREATE TABLE cooking_steps (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id    UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step_number  INTEGER NOT NULL CHECK (step_number > 0),
    instruction  TEXT NOT NULL,
    image_url   TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (recipe_id, step_number)
);
CREATE INDEX idx_cooking_steps_recipe_id ON cooking_steps(recipe_id);
CREATE TRIGGER trg_cooking_steps_updated_at BEFORE UPDATE ON cooking_steps
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- recipes -> recipe_spices (1:N)
-- Stores additional spices or seasonings needed for a recipe
CREATE TABLE recipe_spices (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id    UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    name         VARCHAR(150) NOT NULL,
    unit         VARCHAR(50) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_recipe_spices_recipe_id ON recipe_spices(recipe_id);
CREATE TRIGGER trg_recipe_spices_updated_at BEFORE UPDATE ON recipe_spices
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create photos and contacts
-- ---------------------------------------------------------------------
-- child -> child_photos (1:N), enum visibility
-- Stores photos of the child with visibility settings
CREATE TABLE child_photos (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id            UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    photo_url           TEXT NOT NULL,
    caption             TEXT,
    visibility          child_photo_visibility NOT NULL DEFAULT 'private',
    is_review_required  BOOLEAN NOT NULL DEFAULT false,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_photos_child_id ON child_photos(child_id);
CREATE TRIGGER trg_child_photos_updated_at BEFORE UPDATE ON child_photos
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- mother_profiles -> contacts (1:N)
-- Stores contact relationships between mother profiles
CREATE TABLE contacts (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mother_profile_id           UUID NOT NULL REFERENCES mother_profiles(id) ON DELETE CASCADE,
    related_mother_profile_id   UUID NOT NULL REFERENCES mother_profiles(id) ON DELETE CASCADE,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_contacts_mother_profile_id ON contacts(mother_profile_id);
CREATE TRIGGER trg_contacts_updated_at BEFORE UPDATE ON contacts
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Junction N:M: child_photos <-> contacts
-- Junction table for sharing child photos with specific contacts
CREATE TABLE photo_shares (
    child_photo_id  UUID NOT NULL REFERENCES child_photos(id) ON DELETE CASCADE,
    contact_id      UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_photo_id, contact_id)
);
CREATE INDEX idx_photo_shares_photo_id ON photo_shares(child_photo_id);
CREATE INDEX idx_photo_shares_contact_id ON photo_shares(contact_id);

-- ---------------------------------------------------------------------
-- create caregiver engagements
-- ---------------------------------------------------------------------
-- caregiver_profiles -> caregiver_engagements (1:N), child -> caregiver_engagements (1:N)
-- Stores active engagements mapping caregivers to children
CREATE TABLE caregiver_engagements (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caregiver_profile_id   UUID NOT NULL REFERENCES caregiver_profiles(id) ON DELETE CASCADE,
    child_id               UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    deleted_at             TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_caregiver_engagements_caregiver_id ON caregiver_engagements(caregiver_profile_id);
CREATE INDEX idx_caregiver_engagements_child_id ON caregiver_engagements(child_id);
CREATE UNIQUE INDEX unique_active_engagement ON caregiver_engagements (child_id, caregiver_profile_id) WHERE deleted_at IS NULL;
CREATE TRIGGER trg_caregiver_engagements_updated_at BEFORE UPDATE ON caregiver_engagements
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create medical
-- ---------------------------------------------------------------------
-- doctor_profiles -> medical_relationships (1:N), child -> medical_relationships (1:N)
-- Stores relationships connecting doctors to children
CREATE TABLE medical_relationships (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    doctor_profile_id   UUID NOT NULL REFERENCES doctor_profiles(id) ON DELETE CASCADE,
    child_id            UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    deleted_at          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_medical_relationships_doctor_id ON medical_relationships(doctor_profile_id);
CREATE INDEX idx_medical_relationships_child_id ON medical_relationships(child_id);
CREATE TRIGGER trg_medical_relationships_updated_at BEFORE UPDATE ON medical_relationships
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- medical_relationships -> medical_notes (1:N)
-- Stores medical notes or recommendations provided by doctors
CREATE TABLE medical_notes (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_relationship_id     UUID NOT NULL REFERENCES medical_relationships(id) ON DELETE CASCADE,
    recommendation              TEXT NOT NULL,
    valid_until                  DATE NOT NULL,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_medical_notes_relationship_id ON medical_notes(medical_relationship_id);
CREATE TRIGGER trg_medical_notes_updated_at BEFORE UPDATE ON medical_notes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- medical_notes -> medical_restrictions (1:N), enum type
-- Stores specific medical restrictions (e.g. prohibitions, allergies) from notes
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

-- medical_notes -> daily_nutrition_targets (1:N)
-- Stores specific daily macro/nutrient targets prescribed in medical notes
CREATE TABLE daily_nutrition_targets (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    medical_note_id  UUID NOT NULL REFERENCES medical_notes(id) ON DELETE CASCADE,
    nutrient         nutrient_type NOT NULL,
    quantity         NUMERIC(8,2) NOT NULL CHECK (quantity > 0),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_daily_nutrition_targets_note_id ON daily_nutrition_targets(medical_note_id);
CREATE TRIGGER trg_daily_nutrition_targets_updated_at BEFORE UPDATE ON daily_nutrition_targets
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create notification
-- ---------------------------------------------------------------------
-- user -> notifications (1:N)
-- Stores system notifications for users
CREATE TABLE notifications (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title               VARCHAR(150) NOT NULL,
    message             TEXT NOT NULL,
    notification_type   notification_type NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE TRIGGER trg_notifications_updated_at BEFORE UPDATE ON notifications
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create checkin QR module
-- ---------------------------------------------------------------------
-- Stores short-lived QR tokens for caregiver check-in
CREATE TABLE qr_tokens (
    token VARCHAR(255) PRIMARY KEY,
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Stores logs of caregiver check-ins using QR tokens
CREATE TABLE qr_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(255) NOT NULL,
    child_id UUID NOT NULL REFERENCES children(id) ON DELETE CASCADE,
    caregiver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    checked_in_at TIMESTAMPTZ DEFAULT NOW()
);

-- Extension untuk gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Fungsi generik untuk auto-update kolom updated_at di semua tabel
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
CREATE TYPE child_allergy_category AS ENUM ('food', 'medicine', 'animal', 'others');

-- ENUM: child_growth_analyses.analysis_type
CREATE TYPE growth_analysis_type AS ENUM (
    'weight_for_age',
    'height_for_age',
    'weight_for_height',
    'bmi_for_age',
    'head_circumference_for_age'
);

-- ENUM: nerve_name (aspek tumbuh kembang KPSP)
-- ASUMSI: dipakai di assessment_kpsp_questions.nerve_name dan checklist_milestone_tasks.nerve_name
CREATE TYPE nerve_name AS ENUM (
    'Gross motor skills',
    'Fine motor skills',
    'Speech and language',
    'Socialization'
);

-- ENUM: recipes.meal_time
CREATE TYPE recipe_meal_time AS ENUM (
    'sarapan',
    'makan siang',
    'makan malam',
    'selingan siang',
    'selingan sore'
);

-- ENUM: child_photos.visibility
CREATE TYPE child_photo_visibility AS ENUM ('all', 'private', 'only');

-- ENUM: medical_restrictions.type
CREATE TYPE medical_restriction_type AS ENUM ('prohibition', 'allergy');

-- ENUM: notification.type
CREATE TYPE notification_type AS ENUM ('meal_reminder', 'growth_development_reminder', 'doctor_activity', 'photos_activity', 'shop_activity', 'invitation_activity', 'new_menu_reminder', 'achievement');

-- users (already created in 000001, adding new columns)
ALTER TABLE users 
    ADD COLUMN full_name VARCHAR(150),
    ADD COLUMN gender gender_type,
    ADD COLUMN phone_number VARCHAR(20),
    ADD COLUMN photo_url TEXT,
    ADD COLUMN deleted_at TIMESTAMPTZ;

CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- user -> mother_profiles (1:1)
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
-- mother_profiles -> child (1:N)
CREATE TABLE child (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mother_profile_id       UUID NOT NULL REFERENCES mother_profiles(id) ON DELETE CASCADE,
    full_name               VARCHAR(150) NOT NULL,
    gender                  VARCHAR(10) NOT NULL CHECK (gender IN ('male', 'female')),
    birth_date              DATE NOT NULL,
    photo_url               TEXT,
    food_frequency_profile  INTEGER NOT NULL,
    food_goal_profile       VARCHAR(100) NOT NULL,
    notes_profile           TEXT,
    upload_streak_days      INTEGER NOT NULL CHECK (upload_streak_days >= 0),
    deleted_at              TIMESTAMPTZ,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_mother_profile_id ON child(mother_profile_id);
CREATE TRIGGER trg_child_updated_at BEFORE UPDATE ON child
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> favorite_food_profile (1:N)
CREATE TABLE favorite_food_profile (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id    UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    food_name   VARCHAR(150) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_favorite_food_profile_child_id ON favorite_food_profile(child_id);
CREATE TRIGGER trg_favorite_food_profile_updated_at BEFORE UPDATE ON favorite_food_profile
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> favorite_texture_profile (1:N)
CREATE TABLE favorite_texture_profile (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id      UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    texture_name  VARCHAR(100) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_favorite_texture_profile_child_id ON favorite_texture_profile(child_id);
CREATE TRIGGER trg_favorite_texture_profile_updated_at BEFORE UPDATE ON favorite_texture_profile
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> child_diet_profile (1:N)
CREATE TABLE child_diet_profile (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id     UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    diet_name    VARCHAR(150) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_diet_profile_child_id ON child_diet_profile(child_id);
CREATE TRIGGER trg_child_diet_profile_updated_at BEFORE UPDATE ON child_diet_profile
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> child_chronic_disease_profile (1:N)
CREATE TABLE child_chronic_disease_profile (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id       UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    disease_name   VARCHAR(150) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_chronic_disease_profile_child_id ON child_chronic_disease_profile(child_id);
CREATE TRIGGER trg_child_chronic_disease_profile_updated_at BEFORE UPDATE ON child_chronic_disease_profile
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child -> child_allergy_profile (1:N), enum category
CREATE TABLE child_allergy_profile (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id       UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    category       child_allergy_category NOT NULL,
    allergen_name  VARCHAR(150) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_allergy_profile_child_id ON child_allergy_profile(child_id);
CREATE TRIGGER trg_child_allergy_profile_updated_at BEFORE UPDATE ON child_allergy_profile
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create child growth
-- ---------------------------------------------------------------------
-- child -> child_growth_reports (1:N)
CREATE TABLE child_growth_reports (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id                UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
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

-- child_growth_reports -> child_growth_analyses (1:N), enum analysis_type
CREATE TABLE child_growth_analyses (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_growth_report_id   UUID NOT NULL REFERENCES child_growth_reports(id) ON DELETE CASCADE,
    analysis_type            growth_analysis_type NOT NULL,
    z_score                  NUMERIC(5,2),
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_growth_analyses_report_id ON child_growth_analyses(child_growth_report_id);
CREATE TRIGGER trg_child_growth_analyses_updated_at BEFORE UPDATE ON child_growth_analyses
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create child development
-- ---------------------------------------------------------------------
-- child -> child_development_reports (1:N)
-- ASUMSI: month_target di sini = usia (bulan) saat asesmen KPSP dilakukan
CREATE TABLE child_development_reports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id        UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    kpsp_score      INTEGER NOT NULL,
    month_target    INTEGER NOT NULL CHECK (month_target IN (3,6,9,12,15,18,21,24,30,36,42,48,54,60)),
    next_check_date DATE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_development_reports_child_id ON child_development_reports(child_id);
CREATE TRIGGER trg_child_development_reports_updated_at BEFORE UPDATE ON child_development_reports
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Master pertanyaan KPSP per kelompok umur & aspek (nerve_name)
CREATE TABLE assessment_kpsp_questions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    month_target   INTEGER NOT NULL CHECK (month_target IN (3,6,9,12,15,18,21,24,30,36,42,48,54,60)),
    nerve_name     nerve_name NOT NULL,
    question_text  TEXT NOT NULL,
    description    TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_assessment_kpsp_questions_month_target ON assessment_kpsp_questions(month_target);
CREATE TRIGGER trg_assessment_kpsp_questions_updated_at BEFORE UPDATE ON assessment_kpsp_questions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child_development_reports -> assessment_kpsp_answers (1:N)
-- assessment_kpsp_questions -> assessment_kpsp_answers (1:N)
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

-- assessment_kpsp_questions -> recommended_actions (1:N)
CREATE TABLE recommended_actions (
    id                            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_kpsp_question_id   UUID NOT NULL REFERENCES assessment_kpsp_questions(id) ON DELETE CASCADE,
    title                         VARCHAR(50) NOT NULL,
    action_text                   TEXT NOT NULL,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_recommended_actions_question_id ON recommended_actions(assessment_kpsp_question_id);
CREATE TRIGGER trg_recommended_actions_updated_at BEFORE UPDATE ON recommended_actions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Junction N:M: child_development_reports <-> recommended_actions
CREATE TABLE assessment_attempt_recomendations (
    child_development_report_id   UUID NOT NULL REFERENCES child_development_reports(id) ON DELETE CASCADE,
    recommended_action_id         UUID NOT NULL REFERENCES recommended_actions(id) ON DELETE CASCADE,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_development_report_id, recommended_action_id)
);
CREATE INDEX idx_attempt_recomendations_report_id ON assessment_attempt_recomendations(child_development_report_id);
CREATE INDEX idx_attempt_recomendations_action_id ON assessment_attempt_recomendations(recommended_action_id);

-- Master task milestone per kelompok umur & aspek (nerve_name)
CREATE TABLE checklist_milestone_tasks (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    month_target     INTEGER NOT NULL CHECK (month_target IN (3,6,9,12,15,18,21,24,30,36,42,48,54,60)),
    nerve_name       nerve_name NOT NULL,
    task_description TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_checklist_milestone_tasks_month_target ON checklist_milestone_tasks(month_target);
CREATE TRIGGER trg_checklist_milestone_tasks_updated_at BEFORE UPDATE ON checklist_milestone_tasks
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Junction N:M: child <-> checklist_milestone_tasks
CREATE TABLE checklist_milestone_progress (
    child_id                      UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    checklist_milestone_task_id   UUID NOT NULL REFERENCES checklist_milestone_tasks(id) ON DELETE CASCADE,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_id, checklist_milestone_task_id)
);
CREATE INDEX idx_checklist_milestone_progress_child_id ON checklist_milestone_progress(child_id);
CREATE INDEX idx_checklist_milestone_progress_task_id ON checklist_milestone_progress(checklist_milestone_task_id);
CREATE TRIGGER trg_checklist_milestone_progress_updated_at BEFORE UPDATE ON checklist_milestone_progress
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create child nutrition
-- ---------------------------------------------------------------------
-- child -> child_nutrition_reports (1:N)
CREATE TABLE child_nutrition_reports (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id                UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    calories                INTEGER NOT NULL CHECK (calories >= 0),
    target_calories         INTEGER NOT NULL CHECK (target_calories >= 0),
    protein                 INTEGER NOT NULL CHECK (protein >= 0),
    target_protein          INTEGER NOT NULL CHECK (target_protein >= 0),
    fat                     INTEGER NOT NULL CHECK (fat >= 0),
    target_fat              INTEGER NOT NULL CHECK (target_fat >= 0),
    carbohydrate            INTEGER NOT NULL CHECK (carbohydrate >= 0),
    target_carbohydrate     INTEGER NOT NULL CHECK (target_carbohydrate >= 0),
    created_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_child_nutrition_reports_child_id ON child_nutrition_reports(child_id);
CREATE TRIGGER trg_child_nutrition_reports_updated_at BEFORE UPDATE ON child_nutrition_reports
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- child_nutrition_reports -> daily_menus (1:N)
CREATE TABLE daily_menus (
    id                            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_nutrition_report_id     UUID NOT NULL REFERENCES child_nutrition_reports(id) ON DELETE CASCADE,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_daily_menus_report_id ON daily_menus(child_nutrition_report_id);
CREATE TRIGGER trg_daily_menus_updated_at BEFORE UPDATE ON daily_menus
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- daily_menus -> recipes (1:N), enum meal_time, alergen
CREATE TABLE recipes (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    daily_menu_id  UUID NOT NULL REFERENCES daily_menus(id) ON DELETE CASCADE,
    name           VARCHAR(150) NOT NULL,
    meal_time      recipe_meal_time NOT NULL,
    meal_texture   VARCHAR(25) NOT NULL,
    is_alergen     BOOLEAN NOT NULL,
    calories       INTEGER NOT NULL CHECK (calories >= 0),
    protein        INTEGER NOT NULL CHECK (protein >= 0),
    carbohydrate   INTEGER NOT NULL CHECK (carbohydrate >= 0),
    fat            INTEGER NOT NULL CHECK (fat >= 0),
    description    TEXT NOT NULL,
    cooking_time   INTEGER NOT NULL,
    is_bookmarked  BOOLEAN NOT NULL,
    is_completed   BOOLEAN NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_recipes_daily_menu_id ON recipes(daily_menu_id);
CREATE TRIGGER trg_recipes_updated_at BEFORE UPDATE ON recipes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Master ingredients (dipakai main_ingredients & ingredient_shopping_items)
CREATE TABLE ingredients (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(150) NOT NULL UNIQUE,
    image_url   TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TRIGGER trg_ingredients_updated_at BEFORE UPDATE ON ingredients
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- recipes -> main_ingredients (1:N), ingredients -> main_ingredients (1:N)
CREATE TABLE main_ingredients (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_id      UUID NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id  UUID NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity       NUMERIC(8,2) NOT NULL CHECK (quantity > 0),
    unit           VARCHAR(20) NOT NULL,
    priority       INTEGER NOT NULL CHECK (priority > 0),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_main_ingredients_recipe_id ON main_ingredients(recipe_id);
CREATE INDEX idx_main_ingredients_ingredient_id ON main_ingredients(ingredient_id);
CREATE TRIGGER trg_main_ingredients_updated_at BEFORE UPDATE ON main_ingredients
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- recipes -> cooking_steps (1:N)
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

-- child_nutrition_reports -> daily_shoppings (1:N)
CREATE TABLE daily_shoppings (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_nutrition_report_id   UUID NOT NULL REFERENCES child_nutrition_reports(id) ON DELETE CASCADE,
    is_completed                BOOLEAN NOT NULL DEFAULT false,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_daily_shoppings_report_id ON daily_shoppings(child_nutrition_report_id);
CREATE TRIGGER trg_daily_shoppings_updated_at BEFORE UPDATE ON daily_shoppings
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- daily_shoppings -> ingredient_shopping_items (1:N), ingredients -> ingredient_shopping_items (1:N)
CREATE TABLE ingredient_shopping_items (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    daily_shopping_id   UUID NOT NULL REFERENCES daily_shoppings(id) ON DELETE CASCADE,
    ingredient_id       UUID NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity            NUMERIC(8,2) NOT NULL CHECK (quantity > 0),
    unit                VARCHAR(20) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_ingredient_shopping_items_shopping_id ON ingredient_shopping_items(daily_shopping_id);
CREATE INDEX idx_ingredient_shopping_items_ingredient_id ON ingredient_shopping_items(ingredient_id);
CREATE TRIGGER trg_ingredient_shopping_items_updated_at BEFORE UPDATE ON ingredient_shopping_items
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create photos and contacts
-- ---------------------------------------------------------------------
-- child -> child_photos (1:N), enum visibility
CREATE TABLE child_photos (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    child_id            UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
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
CREATE TABLE photo_shared_with (
    child_photo_id  UUID NOT NULL REFERENCES child_photos(id) ON DELETE CASCADE,
    contact_id      UUID NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (child_photo_id, contact_id)
);
CREATE INDEX idx_photo_shared_with_photo_id ON photo_shared_with(child_photo_id);
CREATE INDEX idx_photo_shared_with_contact_id ON photo_shared_with(contact_id);

-- ---------------------------------------------------------------------
-- create caregiver engagements
-- ---------------------------------------------------------------------
-- caregiver_profiles -> caregiver_engagements (1:N), child -> caregiver_engagements (1:N)
CREATE TABLE caregiver_engagements (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    caregiver_profile_id   UUID NOT NULL REFERENCES caregiver_profiles(id) ON DELETE CASCADE,
    child_id               UUID NOT NULL REFERENCES child(id) ON DELETE CASCADE,
    deleted_at             TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_caregiver_engagements_caregiver_id ON caregiver_engagements(caregiver_profile_id);
CREATE INDEX idx_caregiver_engagements_child_id ON caregiver_engagements(child_id);
CREATE TRIGGER trg_caregiver_engagements_updated_at BEFORE UPDATE ON caregiver_engagements
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ---------------------------------------------------------------------
-- create medical
-- ---------------------------------------------------------------------
-- doctor_profiles -> medical_relationships (1:N), child -> medical_relationships (1:N)
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

-- medical_relationships -> medical_notes (1:N)
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

-- medical_notes -> medical_restrictions (1:N), enum type
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

-- medical_notes -> daily_nutrition_targets (1:N)
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

-- ---------------------------------------------------------------------
-- create notification
-- ---------------------------------------------------------------------
-- user -> notification (1:N)
CREATE TABLE notification (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title               VARCHAR(150) NOT NULL,
    message             TEXT NOT NULL,
    notification_type   notification_type NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notification_user_id ON notification(user_id);
CREATE TRIGGER trg_notification_updated_at BEFORE UPDATE ON notification
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();


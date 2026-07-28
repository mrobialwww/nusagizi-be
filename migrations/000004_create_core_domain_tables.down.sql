-- ---------------------------------------------------------------------
-- create notification
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS notification;

-- ---------------------------------------------------------------------
-- create medical
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS daily_nutrition_targets;
DROP TABLE IF EXISTS medical_restrictions;
DROP TABLE IF EXISTS medical_notes;
DROP TABLE IF EXISTS medical_relationships;

-- ---------------------------------------------------------------------
-- create caregiver engagements
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS caregiver_engagements;

-- ---------------------------------------------------------------------
-- create photos and contacts
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS photo_shared_with;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS child_photos;

-- ---------------------------------------------------------------------
-- create child nutrition
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS ingredient_shopping_items;
DROP TABLE IF EXISTS daily_shoppings;
DROP TABLE IF EXISTS cooking_steps;
DROP TABLE IF EXISTS main_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS recipes;
DROP TABLE IF EXISTS daily_menus;
DROP TABLE IF EXISTS child_nutrition_reports;

-- ---------------------------------------------------------------------
-- create child development
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS checklist_milestone_progress;
DROP TABLE IF EXISTS checklist_milestone_tasks;
DROP TABLE IF EXISTS assessment_attempt_recomendations;
DROP TABLE IF EXISTS recommended_actions;
DROP TABLE IF EXISTS assessment_kpsp_answers;
DROP TABLE IF EXISTS assessment_kpsp_questions;
DROP TABLE IF EXISTS child_development_reports;

-- ---------------------------------------------------------------------
-- create child growth
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS child_growth_reports;

-- ---------------------------------------------------------------------
-- create child core
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS child_allergy_profile;
DROP TABLE IF EXISTS child_chronic_disease_profile;
DROP TABLE IF EXISTS child_diet_profile;
DROP TABLE IF EXISTS favorite_texture_profile;
DROP TABLE IF EXISTS favorite_food_profile;
DROP TABLE IF EXISTS child;

-- ---------------------------------------------------------------------
-- create users and profiles
-- ---------------------------------------------------------------------
DROP TABLE IF EXISTS doctor_profiles;
DROP TABLE IF EXISTS caregiver_profiles;
DROP TABLE IF EXISTS mother_profiles;

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
ALTER TABLE users 
    ADD COLUMN role VARCHAR(50),
    DROP COLUMN IF EXISTS full_name,
    DROP COLUMN IF EXISTS gender,
    DROP COLUMN IF EXISTS phone_number,
    DROP COLUMN IF EXISTS photo_url,
    DROP COLUMN IF EXISTS deleted_at;

-- ---------------------------------------------------------------------
-- create enum types
-- ---------------------------------------------------------------------
DROP TYPE IF EXISTS medical_restriction_type;
DROP TYPE IF EXISTS child_photo_visibility;
DROP TYPE IF EXISTS recipe_meal_time;
DROP TYPE IF EXISTS nerve_name;
DROP TYPE IF EXISTS child_allergy_category;
DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS gender_type;


-- ---------------------------------------------------------------------
-- enable extensions and functions
-- ---------------------------------------------------------------------
DROP FUNCTION IF EXISTS set_updated_at();
DROP EXTENSION IF EXISTS pgcrypto;


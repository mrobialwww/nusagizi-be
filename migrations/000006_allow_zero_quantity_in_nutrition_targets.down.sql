-- Revert: restore quantity constraint to > 0
ALTER TABLE daily_nutrition_targets
    DROP CONSTRAINT daily_nutrition_targets_quantity_check;

ALTER TABLE daily_nutrition_targets
    ADD CONSTRAINT daily_nutrition_targets_quantity_check CHECK (quantity > 0);

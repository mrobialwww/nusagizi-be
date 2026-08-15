package repository

import (
	"context"
	"fmt"
	"time"

	"nusagizi_be/internal/models/menu"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuRepository struct {
	pool *pgxpool.Pool
}

func NewMenuRepository(pool *pgxpool.Pool) *MenuRepository {
	return &MenuRepository{pool: pool}
}

// ReplaceTodayMenu performs a complete "delete then insert" of a child's daily menu in one atomic database transaction.
// If the child already has a menu for today, it deletes the old menu and inserts the new AI-generated one.
// The transaction guarantees that if any insertion fails midway, the entire process rolls back, ensuring the old menu is not accidentally lost and no partial data is saved.
func (r *MenuRepository) ReplaceTodayMenu(
	ctx context.Context,
	childID uuid.UUID,
	today time.Time,
	targetCalories, targetProtein, targetFat, targetCarbo float64,
	menuData menu.FoodEngineMenu,
) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // no-op once committed

	// Get old CNR ID for orphan cleanup
	var existingCNR_ID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT id FROM child_nutrition_reports 
		WHERE child_id = $1 AND report_date = $2`,
		childID, today.Format("2006-01-02")).Scan(&existingCNR_ID)

	if err == nil && existingCNR_ID != uuid.Nil {
		// Collect old recipe IDs
		var oldRecipeIDs []uuid.UUID
		rows, _ := tx.Query(ctx, "SELECT recipe_id FROM child_nutrition_recipes WHERE child_nutrition_report_id = $1", existingCNR_ID)
		defer rows.Close()
		for rows.Next() {
			var id uuid.UUID
			if rows.Scan(&id) == nil {
				oldRecipeIDs = append(oldRecipeIDs, id)
			}
		}

		// Delete existing menu for today (CASCADE automatically removes relations in junction)
		const deleteQuery = `
			DELETE FROM child_nutrition_reports 
			WHERE id = $1`
		if _, err := tx.Exec(ctx, deleteQuery, existingCNR_ID); err != nil {
			return fmt.Errorf("failed to delete existing menu: %w", err)
		}

		// Delete Orphan Recipes
		if len(oldRecipeIDs) > 0 {
			var recipeIDStrs []string
			for _, id := range oldRecipeIDs {
				recipeIDStrs = append(recipeIDStrs, id.String())
			}
			const deleteOrphans = `
				DELETE FROM recipes 
				WHERE id = ANY($1::uuid[]) 
					AND NOT EXISTS (
						SELECT 1 
						FROM child_nutrition_recipes 
						WHERE recipe_id = recipes.id
					)`
			if _, err := tx.Exec(ctx, deleteOrphans, recipeIDStrs); err != nil {
				return fmt.Errorf("failed to delete orphan recipes: %w", err)
			}
		}
	} else if err != pgx.ErrNoRows {
		return fmt.Errorf("failed to check existing menu: %w", err)
	}

	// Insert new nutrition report
	const reportQuery = `
		INSERT INTO child_nutrition_reports
		(child_id, report_date, target_calories, target_protein, target_fat, target_carbohydrate, calories, protein, fat, carbohydrate)
		VALUES ($1, $2, $3, $4, $5, $6, 0, 0, 0, 0) RETURNING id`

	var reportID uuid.UUID
	if err := tx.QueryRow(ctx, reportQuery, childID, today.Format("2006-01-02"), targetCalories, targetProtein, targetFat, targetCarbo).Scan(&reportID); err != nil {
		return fmt.Errorf("failed to insert nutrition report: %w", err)
	}

	// Insert recipes (breakfast, lunch, dinner)
	for mealTime, sesi := range menuData.Sesi {
		// Map AI keys (breakfast, lunch, etc.) to DB ENUM if needed, assuming AI keys are updated to breakfast, lunch, dinner
		dbMealTime := mapMealTime(mealTime)
		waktuMasak := sesi.PanduanMasak.WaktuMasak
		cal, pro, fat, car := 0.0, 0.0, 0.0, 0.0
		if m, ok := sesi.Makro["energi"]; ok {
			cal = m.Dapat
		}
		if m, ok := sesi.Makro["protein"]; ok {
			pro = m.Dapat
		}
		if m, ok := sesi.Makro["lemak"]; ok {
			fat = m.Dapat
		}
		if m, ok := sesi.Makro["karbo"]; ok {
			car = m.Dapat
		}

		// Insert new recipe into the 'recipes' table.
		const recipeQuery = `
			INSERT INTO recipes
			(name, description, meal_time, meal_texture, cooking_time, calories, protein, fat, carbohydrate, is_bookmarked)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, false) RETURNING id`

		var recipeID uuid.UUID
		err = tx.QueryRow(ctx, recipeQuery,
			sesi.PanduanMasak.ResepNama, sesi.PanduanMasak.Catatan, dbMealTime, sesi.PanduanMasak.Tekstur, waktuMasak,
			cal, pro, fat, car,
		).Scan(&recipeID)
		if err != nil {
			return fmt.Errorf("failed to insert recipe %s: %w", mealTime, err)
		}

		// Insert junction
		if _, err := tx.Exec(ctx, `INSERT INTO child_nutrition_recipes (child_nutrition_report_id, recipe_id, portions_consumed) VALUES ($1, $2, 0)`, reportID, recipeID); err != nil {
			return fmt.Errorf("failed to insert junction for %s: %w", mealTime, err)
		}

		if err := insertRecipeDetails(ctx, tx, recipeID, sesi.Bahan, sesi.PanduanMasak); err != nil {
			return err
		}
	}

	// Insert recipes (morning_snack, afternoon_snack)
	for mealTime, selingan := range menuData.Selingan {
		dbMealTime := mapMealTime(mealTime)
		desc := ""
		if len(selingan.Pesan) > 0 {
			desc = selingan.Pesan[0]
		}

		// Insert new snack recipe into the 'recipes' table.
		const selinganRecipeQuery = `
			INSERT INTO recipes
			(name, description, meal_time, meal_texture, cooking_time, calories, protein, fat, carbohydrate, is_bookmarked)
			VALUES ($1, $2, $3, '', '', $4, $5, $6, $7, false) RETURNING id`

		var recipeID uuid.UUID
		err = tx.QueryRow(ctx, selinganRecipeQuery,
			selingan.Resep.Nama, desc, dbMealTime,
			selingan.Makro.Energi, selingan.Makro.Protein, selingan.Makro.Lemak, selingan.Makro.Karbo,
		).Scan(&recipeID)
		if err != nil {
			return fmt.Errorf("failed to insert recipe selingan %s: %w", mealTime, err)
		}

		// Insert junction
		if _, err := tx.Exec(ctx, `INSERT INTO child_nutrition_recipes (child_nutrition_report_id, recipe_id, portions_consumed) VALUES ($1, $2, 0)`, reportID, recipeID); err != nil {
			return fmt.Errorf("failed to insert junction selingan for %s: %w", mealTime, err)
		}

		if err := insertSelinganDetails(ctx, tx, recipeID, selingan.Bahan, selingan.Resep.Langkah); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// ================================== HELPER FUNCTION ===================================
// mapMealTime maps the AI response meal time keys to the database ENUM values.
func mapMealTime(aiKey string) string {
	switch aiKey {
	case "pagi", "sarapan":
		return "breakfast"
	case "siang":
		return "lunch"
	case "malam":
		return "dinner"
	case "selingan_1":
		return "morning_snack"
	case "selingan_2":
		return "afternoon_snack"
	default:
		return aiKey
	}
}

// insertRecipeDetails inserts the "sesi" ingredients and their substitutes for a recipe.
func insertRecipeDetails(ctx context.Context, tx pgx.Tx, recipeID uuid.UUID, bahan []menu.FoodEngineBahan, panduan menu.FoodEnginePanduanMasak) error {
	const queryPriority = `
		INSERT INTO main_ingredients (recipe_id, ingredient_id, unit, slot, priority)
		VALUES ($1, $2, $3, $4, $5)`

	for _, b := range bahan {
		// Priority 1 for main ingredient
		if _, err := tx.Exec(ctx, queryPriority, recipeID, b.Kode, b.Satuan, b.Slot, 1); err != nil {
			return fmt.Errorf("failed to insert main ingredient %s: %w", b.Kode, err)
		}
		// Priority 2, 3, 4... for substitutes
		for i, p := range b.Substitutes {
			if _, err := tx.Exec(ctx, queryPriority, recipeID, p.Kode, p.Satuan, b.Slot, i+2); err != nil {
				return fmt.Errorf("failed to insert substitute ingredient %s: %w", p.Kode, err)
			}
		}
	}
	const stepQuery = `
		INSERT INTO cooking_steps (recipe_id, step_number, instruction, image_url)
		VALUES ($1, $2, $3, '')`

	for i, langkah := range panduan.Langkah {
		if _, err := tx.Exec(ctx, stepQuery, recipeID, i+1, langkah); err != nil {
			return fmt.Errorf("failed to insert cooking step %d: %w", i+1, err)
		}
	}

	const spiceQuery = `
		INSERT INTO recipe_spices (recipe_id, name, unit)
		VALUES ($1, $2, $3)`

	for _, bumbu := range panduan.Bumbu {
		if _, err := tx.Exec(ctx, spiceQuery, recipeID, bumbu.Nama, bumbu.Satuan); err != nil {
			return fmt.Errorf("failed to insert recipe spice %s: %w", bumbu.Nama, err)
		}
	}
	return nil
}

// insertSelinganDetails inserts the "selingan" ingredients and their substitutes for a recipe.
func insertSelinganDetails(ctx context.Context, tx pgx.Tx, recipeID uuid.UUID, bahan []menu.FoodEngineBahan, langkah []string) error {
	const queryPriority = `
		INSERT INTO main_ingredients (recipe_id, ingredient_id, unit, slot, priority)
		VALUES ($1, $2, $3, $4, $5)`

	for _, b := range bahan {
		// Priority 1 for main ingredient
		if _, err := tx.Exec(ctx, queryPriority, recipeID, b.Kode, b.Satuan, b.Slot, 1); err != nil {
			return fmt.Errorf("failed to insert main ingredient %s: %w", b.Kode, err)
		}
		// Priority 2, 3, 4... for substitutes
		for i, p := range b.Substitutes {
			if _, err := tx.Exec(ctx, queryPriority, recipeID, p.Kode, p.Satuan, b.Slot, i+2); err != nil {
				return fmt.Errorf("failed to insert substitute ingredient %s: %w", p.Kode, err)
			}
		}
	}
	const stepQuery = `
		INSERT INTO cooking_steps (recipe_id, step_number, instruction, image_url)
		VALUES ($1, $2, $3, '')`

	for i, l := range langkah {
		if _, err := tx.Exec(ctx, stepQuery, recipeID, i+1, l); err != nil {
			return fmt.Errorf("failed to insert selingan cooking step %d: %w", i+1, err)
		}
	}
	return nil
}

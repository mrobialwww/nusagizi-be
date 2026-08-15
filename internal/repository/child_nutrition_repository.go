package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	child_nutri "nusagizi_be/internal/models/child_nutrition"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChildNutritionRepository struct {
	pool *pgxpool.Pool
}

func NewChildNutritionRepository(pool *pgxpool.Pool) *ChildNutritionRepository {
	return &ChildNutritionRepository{pool: pool}
}

// GetTodayNutritionReport fetches nutrition, menu, recipes, shopping, and items for a specific date.
func (r *ChildNutritionRepository) GetTodayNutritionReport(ctx context.Context, childID uuid.UUID, reportDate time.Time) (*child_nutri.TodayNutritionReportResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var report child_nutri.TodayNutritionReportResponse

	// Fetch Report Base
	queryReport := `
		SELECT id, report_date, calories, target_calories, protein, target_protein, fat, target_fat, carbohydrate, target_carbohydrate
		FROM child_nutrition_reports
		WHERE child_id = $1 AND report_date = $2
	`
	var rDate time.Time
	err := r.pool.QueryRow(ctx, queryReport, childID, reportDate.Format("2006-01-02")).Scan(
		&report.ID, &rDate, &report.Calories, &report.TargetCalories, &report.Protein, &report.TargetProtein,
		&report.Fat, &report.TargetFat, &report.Carbohydrate, &report.TargetCarbohydrate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}

	// Format date and map Report ID to Menu ID for backward compatibility
	report.ReportDate = rDate.Format("2006-01-02")
	report.Menu.ID = report.ID

	// Fetch Recipes from Junction
	queryRecipes := `
		SELECT r.id, r.name, r.meal_time, r.meal_texture, r.calories, r.protein, cnr.portions_consumed, r.is_bookmarked
		FROM recipes r
		JOIN child_nutrition_recipes cnr ON cnr.recipe_id = r.id
		WHERE cnr.child_nutrition_report_id = $1
	`
	rows, err := r.pool.Query(ctx, queryRecipes, report.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rec child_nutri.RecipeResponse
			if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.Calories, &rec.Protein, &rec.PortionsConsumed, &rec.IsBookmarked); err == nil {
				report.Menu.Recipes = append(report.Menu.Recipes, rec)
			}
		}

		// Fetch Shopping List
		queryShopping := `
			SELECT
				i.id   AS ingredient_id,
				i.name AS ingredient_name,
				mi.unit
			FROM main_ingredients mi
			JOIN ingredients i ON i.id = mi.ingredient_id
			JOIN recipes r ON r.id = mi.recipe_id
			JOIN child_nutrition_recipes cnr ON cnr.recipe_id = r.id
			WHERE cnr.child_nutrition_report_id = $1
				AND mi.priority = 1
			ORDER BY i.name ASC
		`
		shoppingRows, err := r.pool.Query(ctx, queryShopping, report.ID)
		if err == nil {
			defer shoppingRows.Close()
			var rawList []child_nutri.ShoppingListItem
			for shoppingRows.Next() {
				var sr child_nutri.ShoppingListItem
				if err := shoppingRows.Scan(&sr.IngredientID, &sr.Name, &sr.Unit); err == nil {
					rawList = append(rawList, sr)
				}
			}
			report.ShoppingList = aggregateShoppingList(rawList)
		} else {
			report.ShoppingList = []child_nutri.ShoppingListItem{}
		}
	} else {
		report.ShoppingList = []child_nutri.ShoppingListItem{}
	}

	if report.ShoppingList == nil {
		report.ShoppingList = []child_nutri.ShoppingListItem{}
	}

	return &report, nil
}

// GetReportMenuByID gets a specific report menu by its ID, ensuring it belongs to the child.
func (r *ChildNutritionRepository) GetReportMenuByID(ctx context.Context, childID, reportID uuid.UUID) (*child_nutri.MenuResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var menu child_nutri.MenuResponse
	queryMenu := `
		SELECT id, created_at
		FROM child_nutrition_reports
		WHERE id = $1 AND child_id = $2
	`
	err := r.pool.QueryRow(ctx, queryMenu, reportID, childID).Scan(&menu.ID, &menu.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Fetch Recipes
	queryRecipes := `
		SELECT r.id, r.name, r.meal_time, r.meal_texture, r.calories, r.protein, cnr.portions_consumed
		FROM recipes r
		JOIN child_nutrition_recipes cnr ON cnr.recipe_id = r.id
		WHERE cnr.child_nutrition_report_id = $1
	`
	rows, err := r.pool.Query(ctx, queryRecipes, menu.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec child_nutri.RecipeResponse
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.Calories, &rec.Protein, &rec.PortionsConsumed); err == nil {
			menu.Recipes = append(menu.Recipes, rec)
		}
	}

	return &menu, nil
}

// GetTodayMenuShopping gets the menu and aggregated shopping list for a specific date.
func (r *ChildNutritionRepository) GetTodayMenuShopping(ctx context.Context, childID uuid.UUID, targetDate time.Time) (*child_nutri.MenuShoppingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var res child_nutri.MenuShoppingResponse

	// Fetch Menu & Recipes using a single query to get menu
	queryMenu := `
		SELECT id, created_at
		FROM child_nutrition_reports
		WHERE child_id = $1 
			AND report_date = $2
	`
	err := r.pool.QueryRow(ctx, queryMenu, childID, targetDate).Scan(&res.Menu.ID, &res.Menu.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			res.ShoppingList = []child_nutri.ShoppingListItem{}
			return &res, nil
		}
		return nil, err
	}

	// Fetch Recipes
	queryRecipes := `
		SELECT r.id, r.name, r.meal_time, r.meal_texture, r.calories, r.protein, cnr.portions_consumed
		FROM recipes r
		JOIN child_nutrition_recipes cnr ON cnr.recipe_id = r.id
		WHERE cnr.child_nutrition_report_id = $1
	`
	rows, err := r.pool.Query(ctx, queryRecipes, res.Menu.ID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var rec child_nutri.RecipeResponse
			if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.Calories, &rec.Protein, &rec.PortionsConsumed); err == nil {
				res.Menu.Recipes = append(res.Menu.Recipes, rec)
			}
		}
	}

	// Fetch Shopping List
	queryShopping := `
		SELECT
			i.id   AS ingredient_id,
			i.name AS ingredient_name,
			mi.unit
		FROM main_ingredients mi
		JOIN ingredients i ON i.id = mi.ingredient_id
		JOIN child_nutrition_recipes cnr ON cnr.recipe_id = mi.recipe_id
		WHERE cnr.child_nutrition_report_id = $1
			AND mi.priority = 1
		ORDER BY i.name ASC
	`
	shoppingRows, err := r.pool.Query(ctx, queryShopping, res.Menu.ID)
	if err == nil {
		defer shoppingRows.Close()
		var rawList []child_nutri.ShoppingListItem
		for shoppingRows.Next() {
			var sr child_nutri.ShoppingListItem
			if err := shoppingRows.Scan(&sr.IngredientID, &sr.Name, &sr.Unit); err == nil {
				rawList = append(rawList, sr)
			}
		}
		res.ShoppingList = aggregateShoppingList(rawList)
	} else {
		res.ShoppingList = []child_nutri.ShoppingListItem{}
	}

	if res.ShoppingList == nil {
		res.ShoppingList = []child_nutri.ShoppingListItem{}
	}

	return &res, nil
}

// GetReportByChildAndDate gets the nutrition report ID and can check existence based on childID and reportDate.
func (r *ChildNutritionRepository) GetReportByChildAndDate(ctx context.Context, childID uuid.UUID, reportDate time.Time) (*child_nutri.TodayNutritionReportResponse, error) {
	var report child_nutri.TodayNutritionReportResponse
	query := `
		SELECT id, report_date 
		FROM child_nutrition_reports 
		WHERE child_id = $1 
			AND report_date = $2`

	var rDate time.Time
	err := r.pool.QueryRow(ctx, query, childID, reportDate.Format("2006-01-02")).Scan(&report.ID, &rDate)
	if err != nil {
		return nil, err
	}
	report.ReportDate = rDate.Format("2006-01-02")
	return &report, nil
}

// IsReportReused checks if any recipe in the given report ID is also linked to another report (meaning it's shared/reused).
func (r *ChildNutritionRepository) IsReportReused(ctx context.Context, reportID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM child_nutrition_recipes cnr
			WHERE cnr.recipe_id IN (
				SELECT recipe_id FROM child_nutrition_recipes WHERE child_nutrition_report_id = $1
			)
			AND cnr.child_nutrition_report_id != $1
		)
	`
	var isReused bool
	err := r.pool.QueryRow(ctx, query, reportID).Scan(&isReused)
	return isReused, err
}

// GetChildIDByRecipeID helper to find child_id given a recipe_id
func (r *ChildNutritionRepository) GetChildIDByRecipeID(ctx context.Context, recipeID uuid.UUID) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var childID uuid.UUID
	query := `
		SELECT cnr.child_id 
		FROM child_nutrition_reports cnr
		JOIN child_nutrition_recipes cj ON cj.child_nutrition_report_id = cnr.id
		WHERE cj.recipe_id = $1
		LIMIT 1
	`
	err := r.pool.QueryRow(ctx, query, recipeID).Scan(&childID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("recipe not found or not linked to any child")
		}
		return uuid.Nil, err
	}
	return childID, nil
}

// UpdateRecipeCompleteStatus updates a recipe's completion status through portions consumed in junction table,
// and automatically updates the total daily macronutrients in the parent report table.
func (r *ChildNutritionRepository) UpdateRecipeCompleteStatus(ctx context.Context, recipeID uuid.UUID, reportID uuid.UUID, portionsConsumed float64) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // no-op if committed

	// Update recipe consumption in junction table
	queryUpdateJunction := `
		UPDATE child_nutrition_recipes
		SET portions_consumed = $1
		WHERE recipe_id = $2 AND child_nutrition_report_id = $3
	`
	tag, err := tx.Exec(ctx, queryUpdateJunction, portionsConsumed, recipeID, reportID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("recipe junction not found")
	}

	// Recalculate daily totals and update parent report table atomically
	queryUpdateReport := `
		WITH recipe_totals AS (
			SELECT
				COALESCE(SUM(r.calories * cnr.portions_consumed), 0)     AS total_calories,
				COALESCE(SUM(r.protein * cnr.portions_consumed), 0)      AS total_protein,
				COALESCE(SUM(r.fat * cnr.portions_consumed), 0)          AS total_fat,
				COALESCE(SUM(r.carbohydrate * cnr.portions_consumed), 0) AS total_carbohydrate
			FROM child_nutrition_recipes cnr
			JOIN recipes r ON r.id = cnr.recipe_id
			WHERE cnr.child_nutrition_report_id = $1
		)
		UPDATE child_nutrition_reports AS report
		SET
			calories     = recipe_totals.total_calories,
			protein      = recipe_totals.total_protein,
			fat          = recipe_totals.total_fat,
			carbohydrate = recipe_totals.total_carbohydrate
		FROM recipe_totals
		WHERE report.id = $1
		`
	if _, err := tx.Exec(ctx, queryUpdateReport, reportID); err != nil {
		return fmt.Errorf("failed to update daily nutrition totals: %w", err)
	}

	return tx.Commit(ctx)
}

// UpdateRecipeBookmarkStatus updates a recipe's bookmark status.
func (r *ChildNutritionRepository) UpdateRecipeBookmarkStatus(ctx context.Context, recipeID uuid.UUID, isBookmarked bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	queryUpdate := `
		UPDATE recipes 
		SET is_bookmarked = $1 
		WHERE id = $2
	`
	res, err := r.pool.Exec(ctx, queryUpdate, isBookmarked, recipeID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

// GetRecipeDetail fetches full recipe details, prioritizing priority=1 main ingredients.
func (r *ChildNutritionRepository) GetRecipeDetail(ctx context.Context, recipeID uuid.UUID) (*child_nutri.RecipeDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var detail child_nutri.RecipeDetailResponse

	// Fetch Recipe
	queryRecipe := `
		SELECT id, name, meal_time, meal_texture, cooking_time, description, calories, protein, is_bookmarked
		FROM recipes
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, queryRecipe, recipeID).Scan(&detail.ID, &detail.Name, &detail.MealTime, &detail.MealTexture, &detail.CookingTime, &detail.Description, &detail.Calories, &detail.Protein, &detail.IsBookmarked)
	if err != nil {
		return nil, err
	}

	// Fetch main ingredients where priority = 1
	queryIngredients := `
		SELECT 
			m.id, 
			m.recipe_id,
			m.unit, 
			m.priority, 
			m.slot,
			i.id, 
			i.name, 
			i.image_url,
			i.category,
			i.price
		FROM main_ingredients m
		JOIN ingredients i ON m.ingredient_id = i.id
		WHERE m.recipe_id = $1 
			AND m.priority = 1
	`
	rows, err := r.pool.Query(ctx, queryIngredients, recipeID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var m child_nutri.MainIngredientResponse
			if err := rows.Scan(
				&m.ID, &m.RecipeID, &m.Unit, &m.Priority, &m.Slot,
				&m.Ingredient.ID, &m.Ingredient.Name, &m.Ingredient.ImageURL, &m.Ingredient.Category, &m.Ingredient.Price,
			); err == nil {
				detail.MainIngredients = append(detail.MainIngredients, m)
			}
		}
	}

	// Fetch cooking steps
	querySteps := `
		SELECT id, recipe_id, step_number, instruction, image_url
		FROM cooking_steps
		WHERE recipe_id = $1
		ORDER BY step_number ASC
	`
	rowsSteps, err := r.pool.Query(ctx, querySteps, recipeID)
	if err == nil {
		defer rowsSteps.Close()
		for rowsSteps.Next() {
			var s child_nutri.CookingStep
			if err := rowsSteps.Scan(&s.ID, &s.RecipeID, &s.StepNumber, &s.Instruction, &s.ImageURL); err == nil {
				detail.CookingSteps = append(detail.CookingSteps, s)
			}
		}
	}

	// Fetch recipe spices
	querySpices := `
		SELECT id, name, unit
		FROM recipe_spices
		WHERE recipe_id = $1
	`
	rowsSpices, err := r.pool.Query(ctx, querySpices, recipeID)
	if err == nil {
		defer rowsSpices.Close()
		for rowsSpices.Next() {
			var sp child_nutri.RecipeSpice
			if err := rowsSpices.Scan(&sp.ID, &sp.Name, &sp.Unit); err == nil {
				detail.RecipeSpices = append(detail.RecipeSpices, sp)
			}
		}
	}

	return &detail, nil
}

// SwapMainIngredientPriority sets target ingredient to priority 1, shifts others down for multiple requests atomically.
func (r *ChildNutritionRepository) SwapMainIngredientPriority(ctx context.Context, requests []child_nutri.SwapPriorityRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		WITH ranked AS (
			SELECT id,
					ROW_NUMBER() OVER (
						ORDER BY CASE WHEN priority = $3 THEN 0 ELSE priority END ASC
					) AS new_priority
			FROM main_ingredients
			WHERE recipe_id = $1 AND slot = $2 AND priority <= $3
		)
		UPDATE main_ingredients mi
		SET priority = r.new_priority
		FROM ranked r
		WHERE mi.id = r.id;
	`

	for _, req := range requests {
		tag, err := tx.Exec(ctx, query, req.RecipeID, req.Slot, req.Priority)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("record not found or no rows updated for recipe_id %s", req.RecipeID)
		}
	}

	return tx.Commit(ctx)
}

// GetBookmarkedRecipes returns all bookmarked recipes for a specific child.
func (r *ChildNutritionRepository) GetBookmarkedRecipes(ctx context.Context, childID uuid.UUID) ([]child_nutri.RecipeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Query across multiple days via the junction table for this specific child
	query := `
		SELECT DISTINCT
			r.id, 
			r.name, 
			r.meal_time, 
			r.meal_texture, 
			r.calories, 
			r.protein, 
			r.is_bookmarked
		FROM recipes r
		JOIN child_nutrition_recipes cj ON r.id = cj.recipe_id
		JOIN child_nutrition_reports cnr ON cj.child_nutrition_report_id = cnr.id
		WHERE cnr.child_id = $1 
			AND r.is_bookmarked = true
		ORDER BY r.name ASC
	`
	rows, err := r.pool.Query(ctx, query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []child_nutri.RecipeResponse
	for rows.Next() {
		var rec child_nutri.RecipeResponse
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.Calories, &rec.Protein, &rec.IsBookmarked); err == nil {
			recipes = append(recipes, rec)
		}
	}
	return recipes, nil
}

// GetNutritionReportsByMonth aggregates monthly data.
func (r *ChildNutritionRepository) GetNutritionReportsByMonth(ctx context.Context, childID uuid.UUID, month, year int) ([]child_nutri.NutritionReportSummaryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT 
			r.id, 
			r.created_at, 
			r.report_date,
			r.calories, 
			r.protein, 
			r.fat, 
			r.carbohydrate,
			COALESCE(
				ARRAY_AGG(rec.meal_time) FILTER (WHERE cnr.portions_consumed > 0), 
				'{}'
			) as meal_times
		FROM child_nutrition_reports r
		LEFT JOIN child_nutrition_recipes cnr ON cnr.child_nutrition_report_id = r.id
		LEFT JOIN recipes rec ON rec.id = cnr.recipe_id
		WHERE 
			r.child_id = $1 
			AND EXTRACT(MONTH FROM r.report_date) = $2 
			AND EXTRACT(YEAR FROM r.report_date) = $3
		GROUP BY r.id
		ORDER BY r.report_date ASC
	`
	rows, err := r.pool.Query(ctx, query, childID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []child_nutri.NutritionReportSummaryResponse
	for rows.Next() {
		var res child_nutri.NutritionReportSummaryResponse
		var rDate time.Time
		if err := rows.Scan(&res.ID, &res.CreatedAt, &rDate, &res.Calories, &res.Protein, &res.Fat, &res.Carbohydrate, &res.MealTimes); err == nil {
			res.ReportDate = rDate.Format("2006-01-02")
			results = append(results, res)
		}
	}
	return results, nil
}

// GetDailyShopByMother fetches daily shop ingredients for all children of a mother.
func (r *ChildNutritionRepository) GetDailyShopByMother(ctx context.Context, motherProfileID uuid.UUID, targetDate time.Time) ([]child_nutri.DailyShopIngredient, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT 
			c.full_name, 
			mi.recipe_id, 
			mi.slot, 
			mi.unit, 
			mi.priority,
			i.name AS ingredient_name
		FROM children c
		JOIN child_nutrition_reports cnr ON cnr.child_id = c.id
		JOIN child_nutrition_recipes cjr ON cjr.child_nutrition_report_id = cnr.id
		JOIN recipes r ON r.id = cjr.recipe_id
		JOIN main_ingredients mi ON mi.recipe_id = r.id
		JOIN ingredients i ON i.id = mi.ingredient_id
		WHERE c.mother_profile_id = $1
			AND cnr.report_date = $2
			AND c.deleted_at IS NULL
		ORDER BY c.full_name, mi.recipe_id, mi.slot, mi.priority ASC
	`
	rows, err := r.pool.Query(ctx, query, motherProfileID, targetDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flatRows []FlatIngredientRow
	for rows.Next() {
		var row FlatIngredientRow
		if err := rows.Scan(&row.ChildName, &row.RecipeID, &row.Slot, &row.Unit, &row.Priority, &row.IngredientName); err == nil {
			flatRows = append(flatRows, row)
		}
	}
	return groupDailyShopIngredients(flatRows), nil
}

// GetDailyShopByCaregiver fetches daily shop ingredients for children engaged by a caregiver.
func (r *ChildNutritionRepository) GetDailyShopByCaregiver(ctx context.Context, caregiverProfileID uuid.UUID, targetDate time.Time) ([]child_nutri.DailyShopIngredient, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT 
			c.full_name, 
			mi.recipe_id, 
			mi.slot, 
			mi.unit, 
			mi.priority,
			i.name AS ingredient_name
		FROM children c
		JOIN caregiver_engagements ce ON ce.child_id = c.id
		JOIN child_nutrition_reports cnr ON cnr.child_id = c.id
		JOIN child_nutrition_recipes cjr ON cjr.child_nutrition_report_id = cnr.id
		JOIN recipes r ON r.id = cjr.recipe_id
		JOIN main_ingredients mi ON mi.recipe_id = r.id
		JOIN ingredients i ON i.id = mi.ingredient_id
		WHERE ce.caregiver_profile_id = $1
			AND ce.deleted_at IS NULL
			AND cnr.report_date = $2
			AND c.deleted_at IS NULL
		ORDER BY c.full_name, mi.recipe_id, mi.slot, mi.priority ASC
	`
	rows, err := r.pool.Query(ctx, query, caregiverProfileID, targetDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flatRows []FlatIngredientRow
	for rows.Next() {
		var row FlatIngredientRow
		if err := rows.Scan(&row.ChildName, &row.RecipeID, &row.Slot, &row.Unit, &row.Priority, &row.IngredientName); err == nil {
			flatRows = append(flatRows, row)
		}
	}
	return groupDailyShopIngredients(flatRows), nil
}

// GetRecipeMealTime gets the meal time of a given recipe.
func (r *ChildNutritionRepository) GetRecipeMealTime(ctx context.Context, recipeID uuid.UUID) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const query = `
		SELECT meal_time 
		FROM recipes WHERE id = $1`
	var mealTime string
	err := r.pool.QueryRow(ctx, query, recipeID).Scan(&mealTime)
	return mealTime, err
}

// ReuseRecipeTx replaces a recipe in a target CNR for a specific meal time with a source recipe atomically.
func (r *ChildNutritionRepository) ReuseRecipeTx(ctx context.Context, targetCNR_ID uuid.UUID, sourceRecipeID uuid.UUID, targetMealTime string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	const getOldRecipeQuery = `
		SELECT r.id 
		FROM recipes r
		JOIN child_nutrition_recipes cnr ON cnr.recipe_id = r.id
		WHERE cnr.child_nutrition_report_id = $1 AND r.meal_time = $2`

	const deleteJunctionQuery = `
		DELETE FROM child_nutrition_recipes 
		WHERE child_nutrition_report_id = $1 AND recipe_id = $2`

	const deleteOrphanRecipeQuery = `
		DELETE FROM recipes 
		WHERE id = $1
			AND NOT EXISTS (
				SELECT 1 
				FROM child_nutrition_recipes 
				WHERE recipe_id = $1
			)`

	const insertJunctionQuery = `
		INSERT INTO child_nutrition_recipes (child_nutrition_report_id, recipe_id, portions_consumed)
		VALUES ($1, $2, 0)`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Step 1: Find the old recipe in the target slot with the same meal_time
	var oldRecipeIDInSlot uuid.UUID
	err = tx.QueryRow(ctx, getOldRecipeQuery, targetCNR_ID, targetMealTime).Scan(&oldRecipeIDInSlot)

	if err == nil && oldRecipeIDInSlot != uuid.Nil {
		// Remove relation from target junction since it will be replaced
		if _, err = tx.Exec(ctx, deleteJunctionQuery, targetCNR_ID, oldRecipeIDInSlot); err != nil {
			return err
		}

		// Delete the old recipe if it is now an orphan (not used in any CNR)
		if _, err = tx.Exec(ctx, deleteOrphanRecipeQuery, oldRecipeIDInSlot); err != nil {
			return err
		}
	} else if err != pgx.ErrNoRows {
		return err
	}

	// Step 2: Link sourceRecipeID to target CNR
	if _, err = tx.Exec(ctx, insertJunctionQuery, targetCNR_ID, sourceRecipeID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// ================================== HELPER FUNCTION ===================================
var unitRegex = regexp.MustCompile(`(?i)^([\d.]+)\s*([a-z]+(?:\s+[a-z]+)*)(?:\s*\(\s*([\d.]+)\s*([a-z]+)\s*\))?$`)

// Optional parenthesised part is returned as zeros/empty string when absent.
func parseComplexUnit(s string) (float64, string, float64, string, bool) {
	// parseComplexUnit parses "1.6 sdm (16g)" → (1.6, "sdm", 16, "g", true)
	// or "16g"                                → (16,  "g",   0,  "",  true).
	m := unitRegex.FindStringSubmatch(strings.TrimSpace(s))
	if len(m) == 0 {
		return 0, "", 0, "", false
	}
	v1, _ := strconv.ParseFloat(m[1], 64)
	v2, _ := strconv.ParseFloat(m[3], 64) // safe: ParseFloat("") returns 0
	return v1, m[2], v2, m[4], true
}

// aggregateShoppingList groups items by IngredientID and sums their units if they have matching measurements.
func aggregateShoppingList(items []child_nutri.ShoppingListItem) []child_nutri.ShoppingListItem {
	mainMap := make(map[string]*child_nutri.ShoppingListItem)
	var orderedKeys []string

	for _, item := range items {
		if ext, exists := mainMap[item.IngredientID]; !exists {
			// First time seeing this ingredient — register it and remember insertion order
			orderedKeys = append(orderedKeys, item.IngredientID)
			clone := item // copy value to avoid capturing loop variable by reference
			mainMap[item.IngredientID] = &clone
		} else {
			va1, ua1, vb1, ub1, ok1 := parseComplexUnit(ext.Unit)
			va2, ua2, vb2, ub2, ok2 := parseComplexUnit(item.Unit)

			// Only sum when both units are parseable and their labels match exactly
			if ok1 && ok2 && ua1 == ua2 && ub1 == ub2 {
				sumA := math.Round((va1+va2)*100) / 100
				sumB := math.Round((vb1+vb2)*100) / 100
				if ub1 != "" { // e.g. "1.6 sdm (16g)" + "1.6 sdm (16g)" → "3.2 sdm (32g)"
					ext.Unit = fmt.Sprintf("%g %s (%g%s)", sumA, ua1, sumB, ub1)
				} else { // e.g. "16g" + "16g" → "32 g"
					ext.Unit = fmt.Sprintf("%g %s", sumA, ua1)
				}
			}
		}
	}

	result := make([]child_nutri.ShoppingListItem, 0, len(orderedKeys))
	for _, k := range orderedKeys {
		result = append(result, *mainMap[k])
	}
	return result
}

// FlatIngredientRow is used internally to scan join results before grouping.
type FlatIngredientRow struct {
	ChildName      string
	RecipeID       uuid.UUID
	Slot           string
	Unit           string
	Priority       int
	IngredientName string
}

// groupDailyShopIngredients takes flat rows (pre-sorted by priority ASC) and groups them by (RecipeID, Slot).
func groupDailyShopIngredients(rows []FlatIngredientRow) []child_nutri.DailyShopIngredient {
	// key: "recipeID_slot"
	type key struct {
		recipeID uuid.UUID
		slot     string
	}

	mainMap := make(map[key]*child_nutri.DailyShopIngredient)
	var orderedKeys []key

	for _, row := range rows {
		k := key{recipeID: row.RecipeID, slot: row.Slot}

		if row.Priority == 1 {
			// Main ingredient — guaranteed to appear first due to ASC priority sorting
			if _, exists := mainMap[k]; !exists {
				orderedKeys = append(orderedKeys, k)
				mainMap[k] = &child_nutri.DailyShopIngredient{
					Name:        row.IngredientName,
					Unit:        row.Unit,
					Priority:    row.Priority,
					Slot:        row.Slot,
					RecipeID:    row.RecipeID,
					ChildName:   row.ChildName,
					Substitutes: []child_nutri.DailyShopSubstitute{},
				}
			}
		} else if main, exists := mainMap[k]; exists {
			// Substitute ingredient — its main ingredient parent is guaranteed to exist in the map
			main.Substitutes = append(main.Substitutes, child_nutri.DailyShopSubstitute{
				Name:     row.IngredientName,
				Unit:     row.Unit,
				Priority: row.Priority,
			})
		}
	}

	// Pre-allocate the result slice with a known capacity to optimize memory allocation
	result := make([]child_nutri.DailyShopIngredient, 0, len(orderedKeys))
	for _, k := range orderedKeys {
		result = append(result, *mainMap[k])
	}
	return result // The result is guaranteed to be non-nil due to make(), no need to check length
}

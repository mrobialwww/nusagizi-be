package repository

import (
	"context"
	"fmt"
	"time"
	child_nutri "nusagizi_be/internal/models/child_nutrition"

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

// GetTodayNutritionReport fetches nutrition, menu, recipes, shopping, and items for "today" (latest record).
func (r *ChildNutritionRepository) GetTodayNutritionReport(ctx context.Context, childID uuid.UUID) (*child_nutri.TodayNutritionReportResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var report child_nutri.TodayNutritionReportResponse

	// Fetch Report Base
	queryReport := `
		SELECT id, calories, target_calories, protein, target_protein, fat, target_fat, carbohydrate, target_carbohydrate
		FROM child_nutrition_reports
		WHERE child_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.pool.QueryRow(ctx, queryReport, childID).Scan(
		&report.ID, &report.Calories, &report.TargetCalories, &report.Protein, &report.TargetProtein,
		&report.Fat, &report.TargetFat, &report.Carbohydrate, &report.TargetCarbohydrate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}

	// Fetch Shopping
	queryShopping := `
		SELECT id, is_completed 
		FROM daily_shoppings 
		WHERE child_nutrition_report_id = $1`

	err = r.pool.QueryRow(ctx, queryShopping, report.ID).Scan(&report.Shopping.ID, &report.Shopping.IsCompleted)

	if err == nil {
		// Fetch Shopping Items
		queryItems := `
			SELECT i.id, ing.name, i.quantity, i.unit
			FROM ingredient_shopping_items i
			JOIN ingredients ing ON i.ingredient_id = ing.id
			WHERE i.daily_shopping_id = $1
		`
		rows, err := r.pool.Query(ctx, queryItems, report.Shopping.ID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var item child_nutri.ShoppingItemResponse
				if err := rows.Scan(&item.ID, &item.IngredientName, &item.Quantity, &item.Unit); err == nil {
					report.Shopping.Items = append(report.Shopping.Items, item)
				}
			}
		}
	}

	// Fetch Menu & Recipes
	queryMenu := `
		SELECT id 
		FROM daily_menus 
		WHERE child_nutrition_report_id = $1`

	err = r.pool.QueryRow(ctx, queryMenu, report.ID).Scan(&report.Menu.ID)
	if err == nil {
		queryRecipes := `
			SELECT id, name, meal_time, meal_texture, is_alergen, calories, protein, is_completed
			FROM recipes
			WHERE daily_menu_id = $1
		`
		rows, err := r.pool.Query(ctx, queryRecipes, report.Menu.ID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var rec child_nutri.RecipeResponse
				if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.IsAlergen, &rec.Calories, &rec.Protein, &rec.IsCompleted); err == nil {
					report.Menu.Recipes = append(report.Menu.Recipes, rec)
				}
			}
		}
	}

	return &report, nil
}

// GetTodayDailyMenu gets just the menu for "today".
func (r *ChildNutritionRepository) GetTodayDailyMenu(ctx context.Context, childID uuid.UUID) (*child_nutri.MenuResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// First get latest report ID
	var reportID uuid.UUID
	queryReport := `
		SELECT id 
		FROM child_nutrition_reports 
		WHERE child_id = $1 
		ORDER BY created_at DESC LIMIT 1
	`
	err := r.pool.QueryRow(ctx, queryReport, childID).Scan(&reportID)
	if err != nil {
		return nil, err
	}

	// Fetch Menu
	var menu child_nutri.MenuResponse
	queryMenu := `
		SELECT id 
		FROM daily_menus 
		WHERE child_nutrition_report_id = $1
	`
	err = r.pool.QueryRow(ctx, queryMenu, reportID).Scan(&menu.ID)
	if err != nil {
		return nil, err
	}

	// Fetch Recipes
	queryRecipes := `
		SELECT id, name, meal_time, meal_texture, is_alergen, calories, protein, is_completed
		FROM recipes
		WHERE daily_menu_id = $1
	`
	rows, err := r.pool.Query(ctx, queryRecipes, menu.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec child_nutri.RecipeResponse
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.IsAlergen, &rec.Calories, &rec.Protein, &rec.IsCompleted); err == nil {
			menu.Recipes = append(menu.Recipes, rec)
		}
	}

	return &menu, nil
}

// UpdateRecipeCompleteStatus updates a recipe's completion status.
func (r *ChildNutritionRepository) UpdateRecipeCompleteStatus(ctx context.Context, recipeID uuid.UUID, isCompleted bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	queryUpdate := `
		UPDATE recipes 
		SET is_completed = $1 
		WHERE id = $2
	`
	res, err := r.pool.Exec(ctx, queryUpdate, isCompleted, recipeID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
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
		SELECT id, name, meal_time, meal_texture, is_alergen, calories, protein
		FROM recipes
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, queryRecipe, recipeID).Scan(&detail.ID, &detail.Name, &detail.MealTime, &detail.MealTexture, &detail.IsAlergen, &detail.Calories, &detail.Protein)
	if err != nil {
		return nil, err
	}

	// Fetch main ingredients where priority = 1
	queryIngredients := `
		SELECT 
			m.id, 
			m.recipe_id,
			m.quantity, 
			m.unit, 
			m.priority, 
			m.slot,
			i.id, 
			i.name, 
			i.image_url
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
				&m.ID, &m.RecipeID, &m.Quantity, &m.Unit, &m.Priority, &m.Slot,
				&m.Ingredient.ID, &m.Ingredient.Name, &m.Ingredient.ImageURL,
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

	return &detail, nil
}

// SwapMainIngredientPriority sets target ingredient to priority 1, shifts others down.
func (r *ChildNutritionRepository) SwapMainIngredientPriority(ctx context.Context, recipeID uuid.UUID, mainIngredientID uuid.UUID, slot string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Shift others +1 in the same slot
	queryShiftOthers := `
		UPDATE main_ingredients 
		SET priority = priority + 1 
		WHERE recipe_id = $1 
			AND id != $2 
			AND slot = $3
	`
	_, err = tx.Exec(ctx, queryShiftOthers, recipeID, mainIngredientID, slot)
	if err != nil {
		return err
	}

	// Set target to 1
	querySetTarget := `
		UPDATE main_ingredients 
		SET priority = 1 
		WHERE id = $1 
			AND recipe_id = $2
			AND slot = $3
	`
	res, err := tx.Exec(ctx, querySetTarget, mainIngredientID, recipeID, slot)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}

	return tx.Commit(ctx)
}

// GetBookmarkedRecipes returns all bookmarked recipes for a child.
func (r *ChildNutritionRepository) GetBookmarkedRecipes(ctx context.Context, childID uuid.UUID) ([]child_nutri.RecipeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT 
			r.id, 
			r.name, 
			r.meal_time, 
			r.meal_texture, 
			r.is_alergen, 
			r.calories, 
			r.protein, 
			r.is_bookmarked
		FROM recipes r
		JOIN daily_menus dm ON r.daily_menu_id = dm.id
		JOIN child_nutrition_reports cnr ON dm.child_nutrition_report_id = cnr.id
		WHERE cnr.child_id = $1 
			AND r.is_bookmarked = true
	`
	rows, err := r.pool.Query(ctx, query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []child_nutri.RecipeResponse
	for rows.Next() {
		var rec child_nutri.RecipeResponse
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.IsAlergen, &rec.Calories, &rec.Protein, &rec.IsBookmarked); err == nil {
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
			r.calories, 
			r.protein, 
			r.fat, 
			r.carbohydrate,
			COALESCE(
				ARRAY_AGG(rec.meal_time) FILTER (WHERE rec.is_completed = true), 
				'{}'
			) as meal_times
		FROM child_nutrition_reports r
		LEFT JOIN daily_menus dm ON dm.child_nutrition_report_id = r.id
		LEFT JOIN recipes rec ON rec.daily_menu_id = dm.id
		WHERE 
			r.child_id = $1 
			AND EXTRACT(MONTH FROM r.created_at) = $2 
			AND EXTRACT(YEAR FROM r.created_at) = $3
		GROUP BY r.id
		ORDER BY r.created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, childID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []child_nutri.NutritionReportSummaryResponse
	for rows.Next() {
		var res child_nutri.NutritionReportSummaryResponse
		if err := rows.Scan(&res.ID, &res.CreatedAt, &res.Calories, &res.Protein, &res.Fat, &res.Carbohydrate, &res.MealTimes); err == nil {
			results = append(results, res)
		}
	}
	return results, nil
}

// UpdateDailyShoppingStatus updates shopping completion status.
func (r *ChildNutritionRepository) UpdateDailyShoppingStatus(ctx context.Context, shoppingID uuid.UUID, isCompleted bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	queryUpdate := `
		UPDATE daily_shoppings 
		SET is_completed = $1 
		WHERE id = $2
	`
	res, err := r.pool.Exec(ctx, queryUpdate, isCompleted, shoppingID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

// AddShoppingItem adds a manual item to the shopping list.
func (r *ChildNutritionRepository) AddShoppingItem(ctx context.Context, shoppingID uuid.UUID, item *child_nutri.ShoppingItemInput) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	newID := uuid.New()
	query := `
		INSERT INTO ingredient_shopping_items (id, daily_shopping_id, ingredient_id, quantity, unit)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, newID, shoppingID, item.IngredientID, item.Quantity, item.Unit)
	if err != nil {
		return uuid.Nil, err
	}
	return newID, nil
}

// UpdateShoppingItem patches quantity or unit.
func (r *ChildNutritionRepository) UpdateShoppingItem(ctx context.Context, itemID uuid.UUID, input *child_nutri.UpdateShoppingItemInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		UPDATE ingredient_shopping_items
		SET 
			quantity = COALESCE($1, quantity),
			unit = COALESCE($2, unit)
		WHERE id = $3
	`
	res, err := r.pool.Exec(ctx, query, input.Quantity, input.Unit, itemID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

// DeleteShoppingItem performs hard delete.
func (r *ChildNutritionRepository) DeleteShoppingItem(ctx context.Context, itemID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	queryDelete := `
		DELETE FROM ingredient_shopping_items 
		WHERE id = $1
	`
	res, err := r.pool.Exec(ctx, queryDelete, itemID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

// SearchIngredients provides autocomplete.
func (r *ChildNutritionRepository) SearchIngredients(ctx context.Context, nameQuery string) ([]child_nutri.Ingredient, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT id, name, image_url
		FROM ingredients
		WHERE name ILIKE $1
		ORDER BY name ASC
		LIMIT 20
	`
	// Add wildcard to query
	searchTerm := "%" + nameQuery + "%"
	rows, err := r.pool.Query(ctx, query, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []child_nutri.Ingredient
	for rows.Next() {
		var ing child_nutri.Ingredient
		if err := rows.Scan(&ing.ID, &ing.Name, &ing.ImageURL); err == nil {
			results = append(results, ing)
		}
	}
	return results, nil
}

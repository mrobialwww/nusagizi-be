package repository

import (
	"context"
	"fmt"
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

	// Fetch Menu & Recipes
	queryMenu := `
		SELECT id, created_at 
		FROM daily_menus 
		WHERE child_nutrition_report_id = $1`

	err = r.pool.QueryRow(ctx, queryMenu, report.ID).Scan(&report.Menu.ID, &report.Menu.CreatedAt)
	if err == nil {
		queryRecipes := `
			SELECT id, name, meal_time, meal_texture, calories, protein, portions_consumed
			FROM recipes
			WHERE daily_menu_id = $1
		`
		rows, err := r.pool.Query(ctx, queryRecipes, report.Menu.ID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var rec child_nutri.RecipeResponse
				if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.Calories, &rec.Protein, &rec.PortionsConsumed); err == nil {
					report.Menu.Recipes = append(report.Menu.Recipes, rec)
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
			JOIN recipes r ON r.id = mi.recipe_id
			WHERE r.daily_menu_id = $1
				AND mi.priority = 1
			ORDER BY i.name ASC
		`
		shoppingRows, err := r.pool.Query(ctx, queryShopping, report.Menu.ID)
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
		SELECT id, created_at
		FROM daily_menus 
		WHERE child_nutrition_report_id = $1
	`
	err = r.pool.QueryRow(ctx, queryMenu, reportID).Scan(&menu.ID, &menu.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch Recipes
	queryRecipes := `
		SELECT id, name, meal_time, meal_texture, calories, protein, portions_consumed
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
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.Calories, &rec.Protein, &rec.PortionsConsumed); err == nil {
			menu.Recipes = append(menu.Recipes, rec)
		}
	}

	return &menu, nil
}

// GetTodayMenuShopping gets the menu and aggregated shopping list for "today".
func (r *ChildNutritionRepository) GetTodayMenuShopping(ctx context.Context, childID uuid.UUID) (*child_nutri.MenuShoppingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var res child_nutri.MenuShoppingResponse

	// Fetch Menu & Recipes using a single query to get menu
	queryMenu := `
		SELECT dm.id, dm.created_at
		FROM daily_menus dm
		JOIN child_nutrition_reports cnr ON cnr.id = dm.child_nutrition_report_id
		WHERE cnr.child_id = $1
		ORDER BY cnr.created_at DESC
		LIMIT 1
	`
	err := r.pool.QueryRow(ctx, queryMenu, childID).Scan(&res.Menu.ID, &res.Menu.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			res.ShoppingList = []child_nutri.ShoppingListItem{}
			return &res, nil
		}
		return nil, err
	}

	// Fetch Recipes
	queryRecipes := `
		SELECT id, name, meal_time, meal_texture, calories, protein, portions_consumed
		FROM recipes
		WHERE daily_menu_id = $1
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
		JOIN recipes r ON r.id = mi.recipe_id
		WHERE r.daily_menu_id = $1
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

// GetChildIDByRecipeID helper to find child_id given a recipe_id
func (r *ChildNutritionRepository) GetChildIDByRecipeID(ctx context.Context, recipeID uuid.UUID) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var childID uuid.UUID
	query := `
		SELECT dm.child_nutrition_report_id
		FROM recipes r
		JOIN daily_menus dm ON dm.id = r.daily_menu_id
		WHERE r.id = $1
	`
	var reportID uuid.UUID
	err := r.pool.QueryRow(ctx, query, recipeID).Scan(&reportID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("recipe not found")
		}
		return uuid.Nil, err
	}

	query2 := `
		SELECT child_id 
		FROM child_nutrition_reports 
		WHERE id = $1`

	err = r.pool.QueryRow(ctx, query2, reportID).Scan(&childID)
	if err != nil {
		return uuid.Nil, err
	}

	return childID, nil
}

// UpdateRecipeCompleteStatus updates a recipe's completion status through portions consumed.
func (r *ChildNutritionRepository) UpdateRecipeCompleteStatus(ctx context.Context, recipeID uuid.UUID, portionsConsumed float64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	queryUpdate := `
		UPDATE recipes
		SET portions_consumed = $1
		WHERE id = $2
	`
	tag, err := r.pool.Exec(ctx, queryUpdate, portionsConsumed, recipeID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("recipe not found")
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
		SELECT id, name, meal_time, meal_texture, cooking_time, calories, protein
		FROM recipes
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, queryRecipe, recipeID).Scan(&detail.ID, &detail.Name, &detail.MealTime, &detail.MealTexture, &detail.CookingTime, &detail.Calories, &detail.Protein)
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

// SwapMainIngredientPriority sets target ingredient to priority 1, shifts others down.
func (r *ChildNutritionRepository) SwapMainIngredientPriority(ctx context.Context, recipeID uuid.UUID, slot string, priority int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	tag, err := r.pool.Exec(ctx, query, recipeID, slot, priority)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("record not found or no rows updated")
	}
	return nil
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
			r.calories, 
			r.protein, 
			r.portions_consumed,
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
		if err := rows.Scan(&rec.ID, &rec.Name, &rec.MealTime, &rec.MealTexture, &rec.Calories, &rec.Protein, &rec.PortionsConsumed, &rec.IsBookmarked); err == nil {
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
				ARRAY_AGG(rec.meal_time) FILTER (WHERE rec.portions_consumed > 0), 
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

// GetDailyShopByMother fetches daily shop ingredients for all children of a mother.
func (r *ChildNutritionRepository) GetDailyShopByMother(ctx context.Context, motherProfileID uuid.UUID) ([]child_nutri.DailyShopIngredient, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT 
			c.full_name, 
			mi.recipe_id, 
			mi.slot, 
			mi.unit, 
			mi.priority,
			i.id AS ingredient_id, 
			i.name AS ingredient_name
		FROM children c
		JOIN child_nutrition_reports cnr ON cnr.child_id = c.id
		JOIN daily_menus dm ON dm.child_nutrition_report_id = cnr.id
		JOIN recipes r ON r.daily_menu_id = dm.id
		JOIN main_ingredients mi ON mi.recipe_id = r.id
		JOIN ingredients i ON i.id = mi.ingredient_id
		WHERE c.mother_profile_id = $1
			AND DATE(cnr.created_at) = CURRENT_DATE
			AND c.deleted_at IS NULL
		ORDER BY c.full_name, mi.recipe_id, mi.slot, mi.priority ASC
	`
	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flatRows []FlatIngredientRow
	for rows.Next() {
		var row FlatIngredientRow
		if err := rows.Scan(&row.ChildName, &row.RecipeID, &row.Slot, &row.Unit, &row.Priority, &row.IngredientID, &row.IngredientName); err == nil {
			flatRows = append(flatRows, row)
		}
	}
	return groupDailyShopIngredients(flatRows), nil
}

// GetDailyShopByCaregiver fetches daily shop ingredients for children engaged by a caregiver.
func (r *ChildNutritionRepository) GetDailyShopByCaregiver(ctx context.Context, caregiverProfileID uuid.UUID) ([]child_nutri.DailyShopIngredient, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT 
			c.full_name, 
			mi.recipe_id, 
			mi.slot, 
			mi.unit, 
			mi.priority,
			i.id AS ingredient_id, 
			i.name AS ingredient_name
		FROM children c
		JOIN caregiver_engagements ce ON ce.child_id = c.id
		JOIN child_nutrition_reports cnr ON cnr.child_id = c.id
		JOIN daily_menus dm ON dm.child_nutrition_report_id = cnr.id
		JOIN recipes r ON r.daily_menu_id = dm.id
		JOIN main_ingredients mi ON mi.recipe_id = r.id
		JOIN ingredients i ON i.id = mi.ingredient_id
		WHERE ce.caregiver_profile_id = $1
			AND ce.deleted_at IS NULL
			AND DATE(cnr.created_at) = CURRENT_DATE
			AND c.deleted_at IS NULL
		ORDER BY c.full_name, mi.recipe_id, mi.slot, mi.priority ASC
	`
	rows, err := r.pool.Query(ctx, query, caregiverProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flatRows []FlatIngredientRow
	for rows.Next() {
		var row FlatIngredientRow
		if err := rows.Scan(&row.ChildName, &row.RecipeID, &row.Slot, &row.Unit, &row.Priority, &row.IngredientID, &row.IngredientName); err == nil {
			flatRows = append(flatRows, row)
		}
	}
	return groupDailyShopIngredients(flatRows), nil
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
				if ub1 != "" { // e.g. "1.6 sdm (16g)" + "1.6 sdm (16g)" → "3.2 sdm (32g)"
					ext.Unit = fmt.Sprintf("%g %s (%g%s)", va1+va2, ua1, vb1+vb2, ub1)
				} else { // e.g. "16g" + "16g" → "32 g"
					ext.Unit = fmt.Sprintf("%g %s", va1+va2, ua1)
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
	IngredientID   string
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
			// Bahan utama — karena urut ASC, ini selalu datang lebih dulu
			if _, exists := mainMap[k]; !exists {
				orderedKeys = append(orderedKeys, k)
				mainMap[k] = &child_nutri.DailyShopIngredient{
					Name:         row.IngredientName,
					IngredientID: row.IngredientID,
					Unit:         row.Unit,
					Priority:     row.Priority,
					Slot:         row.Slot,
					RecipeID:     row.RecipeID,
					ChildName:    row.ChildName,
					Pengganti:    []child_nutri.DailyShopSubstitute{},
				}
			}
		} else if main, exists := mainMap[k]; exists {
			// Bahan pengganti — parent-nya sudah pasti ada di map
			main.Pengganti = append(main.Pengganti, child_nutri.DailyShopSubstitute{
				Name:         row.IngredientName,
				IngredientID: row.IngredientID,
				Unit:         row.Unit,
				Priority:     row.Priority,
			})
		}
	}

	// Pre-allocate slice dengan kapasitas yang sudah diketahui
	result := make([]child_nutri.DailyShopIngredient, 0, len(orderedKeys))
	for _, k := range orderedKeys {
		result = append(result, *mainMap[k])
	}
	return result // make() sudah non-nil, tidak perlu cek len == 0
}

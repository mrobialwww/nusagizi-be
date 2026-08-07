package models

import (
	"time"

	"github.com/google/uuid"
)

type RecipeResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	MealTime     string    `json:"meal_time"`
	MealTexture  string    `json:"meal_texture"`
	Calories         int       `json:"calories"`
	Protein          int       `json:"protein"`
	PortionsConsumed float64   `json:"portions_consumed"`
	IsBookmarked     bool      `json:"is_bookmarked,omitempty"`
}

type MenuResponse struct {
	ID        uuid.UUID        `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Recipes   []RecipeResponse `json:"recipes"`
}

type ShoppingListItem struct {
	Name         string `json:"name"`
	IngredientID string `json:"ingredient_id"`
	Unit         string `json:"unit"`
}

type MenuShoppingResponse struct {
	Menu         MenuResponse       `json:"menu"`
	ShoppingList []ShoppingListItem `json:"shopping_list"`
}

type TodayNutritionReportResponse struct {
	ID                 uuid.UUID          `json:"id"`
	Calories           int                `json:"calories"`
	TargetCalories     int                `json:"target_calories"`
	Protein            int                `json:"protein"`
	TargetProtein      int                `json:"target_protein"`
	Fat                int                `json:"fat"`
	TargetFat          int                `json:"target_fat"`
	Carbohydrate       int                `json:"carbohydrate"`
	TargetCarbohydrate int                `json:"target_carbohydrate"`
	Status             string             `json:"status"`
	Menu               MenuResponse       `json:"menu"`
	ShoppingList       []ShoppingListItem `json:"shopping_list"`
}

type MainIngredientResponse struct {
	ID         uuid.UUID  `json:"id"`
	RecipeID   uuid.UUID  `json:"recipe_id"`
	Ingredient Ingredient `json:"ingredient"`
	Unit       string     `json:"unit"`
	Priority   int        `json:"priority"`
	Slot       string     `json:"slot"`
}

type RecipeSpice struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Unit string    `json:"unit"`
}

type RecipeDetailResponse struct {
	ID              uuid.UUID                `json:"id"`
	Name            string                   `json:"name"`
	MealTime        string                   `json:"meal_time"`
	MealTexture     string                   `json:"meal_texture"`
	CookingTime     string                   `json:"cooking_time"`
	Calories        int                      `json:"calories"`
	Protein         int                      `json:"protein"`
	MainIngredients []MainIngredientResponse `json:"main_ingredients"`
	RecipeSpices    []RecipeSpice            `json:"recipe_spices"`
	CookingSteps    []CookingStep            `json:"cooking_steps"`
}

type NutritionReportSummaryResponse struct {
	ID           uuid.UUID `json:"id"`
	Calories     int       `json:"calories"`
	Protein      int       `json:"protein"`
	Fat          int       `json:"fat"`
	Carbohydrate int       `json:"carbohydrate"`
	MealTimes    []string  `json:"meal_times"`
	CreatedAt    time.Time `json:"created_at"`
}

type SwapIngredientPriorityInput struct {
	Slot string `json:"slot" binding:"required,max=25"`
}

type DailyShopSubstitute struct {
	Name         string `json:"name"`
	IngredientID string `json:"ingredient_id"`
	Unit         string `json:"unit"`
	Priority     int    `json:"priority"`
}

type DailyShopIngredient struct {
	Name         string                `json:"name"`
	IngredientID string                `json:"ingredient_id"`
	Unit         string                `json:"unit"`
	Priority     int                   `json:"priority"`
	Slot         string                `json:"slot"`
	RecipeID     uuid.UUID             `json:"recipe_id"`
	ChildName    string                `json:"child_name"`
	Pengganti    []DailyShopSubstitute `json:"pengganti"`
}

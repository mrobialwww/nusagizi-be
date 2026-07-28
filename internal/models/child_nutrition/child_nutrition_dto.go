package models

import (
	"time"

	"github.com/google/uuid"
)

type ShoppingItemResponse struct {
	ID             uuid.UUID `json:"id"`
	IngredientName string    `json:"ingredient_name"`
	Quantity       float64   `json:"quantity"`
	Unit           string    `json:"unit"`
}

type ShoppingResponse struct {
	ID          uuid.UUID              `json:"id"`
	IsCompleted bool                   `json:"is_completed"`
	Items       []ShoppingItemResponse `json:"items"`
}

type RecipeResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	MealTime     string    `json:"meal_time"`
	MealTexture  string    `json:"meal_texture"`
	IsAlergen    bool      `json:"is_alergen"`
	Calories     int       `json:"calories"`
	Protein      int       `json:"protein"`
	IsCompleted  bool      `json:"is_completed,omitempty"`
	IsBookmarked bool      `json:"is_bookmarked,omitempty"`
}

type MenuResponse struct {
	ID      uuid.UUID        `json:"id"`
	Recipes []RecipeResponse `json:"recipes"`
}

type TodayNutritionReportResponse struct {
	ID                 uuid.UUID        `json:"id"`
	Calories           int              `json:"calories"`
	TargetCalories     int              `json:"target_calories"`
	Protein            int              `json:"protein"`
	TargetProtein      int              `json:"target_protein"`
	Fat                int              `json:"fat"`
	TargetFat          int              `json:"target_fat"`
	Carbohydrate       int              `json:"carbohydrate"`
	TargetCarbohydrate int              `json:"target_carbohydrate"`
	Shopping           ShoppingResponse `json:"shopping"`
	Menu               MenuResponse     `json:"menu"`
}

type MainIngredientResponse struct {
	ID         uuid.UUID  `json:"id"`
	RecipeID   uuid.UUID  `json:"recipe_id"`
	Ingredient Ingredient `json:"ingredient"`
	Quantity   float64    `json:"quantity"`
	Unit       string     `json:"unit"`
	Priority   int        `json:"priority"`
	Slot       string     `json:"slot"`
}

type RecipeDetailResponse struct {
	ID              uuid.UUID                `json:"id"`
	Name            string                   `json:"name"`
	MealTime        string                   `json:"meal_time"`
	MealTexture     string                   `json:"meal_texture"`
	IsAlergen       bool                     `json:"is_alergen"`
	Calories        int                      `json:"calories"`
	Protein         int                      `json:"protein"`
	MainIngredients []MainIngredientResponse `json:"main_ingredients"`
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

type ShoppingItemInput struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
}

type UpdateShoppingItemInput struct {
	Quantity *float64 `json:"quantity"`
	Unit     *string  `json:"unit"`
}

type SwapIngredientPriorityInput struct {
	Slot string `json:"slot" binding:"required,max=25"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type RecipeResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	MealTime         string    `json:"meal_time"`
	MealTexture      string    `json:"meal_texture"`
	Calories         float64   `json:"calories"`
	Protein          float64   `json:"protein"`
	PortionsConsumed float64   `json:"portions_consumed"`
	IsBookmarked     bool      `json:"is_bookmarked,omitempty"`
	ImageURL         *string   `json:"image_url"`
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
	ReportDate         string             `json:"report_date"`
	CanRegenerate      bool               `json:"can_regenerate"`
	Calories           float64            `json:"calories"`
	TargetCalories     float64            `json:"target_calories"`
	Protein            float64            `json:"protein"`
	TargetProtein      float64            `json:"target_protein"`
	Fat                float64            `json:"fat"`
	TargetFat          float64            `json:"target_fat"`
	Carbohydrate       float64            `json:"carbohydrate"`
	TargetCarbohydrate float64            `json:"target_carbohydrate"`
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
	Description     string                   `json:"description"`
	Calories        float64                  `json:"calories"`
	Protein         float64                  `json:"protein"`
	Fat             float64                  `json:"fat"`
	Carbohydrate    float64                  `json:"carbohydrate"`
	ImageURL        *string                  `json:"image_url"`
	IsBookmarked    bool                     `json:"is_bookmarked"`
	MainIngredients []MainIngredientResponse `json:"main_ingredients"`
	RecipeSpices    []RecipeSpice            `json:"recipe_spices"`
	CookingSteps    []CookingStep            `json:"cooking_steps"`
}

type NutritionReportSummaryResponse struct {
	ID           uuid.UUID `json:"id"`
	Calories     float64   `json:"calories"`
	Protein      float64   `json:"protein"`
	Fat          float64   `json:"fat"`
	Carbohydrate float64   `json:"carbohydrate"`
	MealTimes    []string  `json:"meal_times"`
	CreatedAt    time.Time `json:"created_at"`
	ReportDate   string    `json:"report_date"`
}

type SwapIngredientPriorityInput struct {
	Slot string `json:"slot" binding:"required,max=25"`
}

type ReuseRecipeInputDTO struct {
	SourceRecipeID uuid.UUID `json:"source_recipe_id" binding:"required"`
	ReportDate     string    `json:"report_date" binding:"required"`
}

type DailyShopSubstitute struct {
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Priority int    `json:"priority"`
}

type DailyShopIngredient struct {
	Name        string                `json:"name"`
	Unit        string                `json:"unit"`
	Priority    int                   `json:"priority"`
	Slot        string                `json:"slot"`
	RecipeID    uuid.UUID             `json:"recipe_id"`
	ChildName   string                `json:"child_name"`
	Substitutes []DailyShopSubstitute `json:"substitutes"`
}

type SwapPriorityRequest struct {
	RecipeID uuid.UUID `json:"recipe_id"`
	Slot     string    `json:"slot"`
	Priority int       `json:"priority"`
}

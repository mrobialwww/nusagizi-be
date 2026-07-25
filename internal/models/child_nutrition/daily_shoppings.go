package models

import (
	"time"

	"github.com/google/uuid"
)

type DailyShopping struct {
	ID                     uuid.UUID `json:"daily_shopping_id" db:"daily_shopping_id"`
	ChildNutritionReportID uuid.UUID `json:"child_nutrition_report_id" db:"child_nutrition_report_id"`
	IsDone                 bool      `json:"is_done" db:"is_done"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`

	// Relation 1:N to IngredientShoppingItem
	IngredientShoppingItems []IngredientShoppingItem `json:"ingredient_shopping_items,omitempty" db:"-"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type IngredientShoppingItem struct {
	ID              uuid.UUID `json:"ingredient_shopping_item_id" db:"ingredient_shopping_item_id"`
	IngredientID    uuid.UUID `json:"ingredient_id" db:"ingredient_id"`
	DailyShoppingID uuid.UUID `json:"daily_shopping_id" db:"daily_shopping_id"`
	Quantity        int       `json:"quantity" db:"quantity"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

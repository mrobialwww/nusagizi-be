package models

import (
	"time"

	"github.com/google/uuid"
)

type MainIngredient struct {
	ID           uuid.UUID `json:"main_ingredients_id" db:"main_ingredients_id"`
	FoodID       uuid.UUID `json:"food_id" db:"food_id"`
	IngredientID uuid.UUID `json:"ingredient_id" db:"ingredient_id"`
	Quantity     float64   `json:"quantity" db:"quantity"`
	Unit         string    `json:"unit" db:"unit"`
	Priority     string    `json:"priority" db:"priority"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

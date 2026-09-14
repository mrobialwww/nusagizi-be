package models

import (
	"time"

	"github.com/google/uuid"
)

type CookingStep struct {
	ID          uuid.UUID `json:"id" db:"id"`
	RecipeID    uuid.UUID `json:"recipe_id" db:"recipe_id"`
	StepNumber  int       `json:"step_number" db:"step_number"`
	Instruction string    `json:"instruction" db:"instruction"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	ID        uuid.UUID `json:"ingredient_id" db:"ingredient_id"`
	Name      string    `json:"name" db:"name"`
	ImageURL  string    `json:"image_url" db:"image_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

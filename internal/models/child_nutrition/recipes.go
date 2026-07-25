package models

import (
	"time"

	"github.com/google/uuid"
)

type MealTimeEnum string

const (
	MealTimeSarapan       MealTimeEnum = "SARAPAN"
	MealTimeMakanSiang    MealTimeEnum = "MAKAN_SIANG"
	MealTimeMakanMalam    MealTimeEnum = "MAKAN_MALAM"
	MealTimeSelinganSiang MealTimeEnum = "SELINGAN_SIANG"
	MealTimeSelinganSore  MealTimeEnum = "SELINGAN_SORE"
)

// IsValid memvalidasi bahwa nilai MealTimeEnum sesuai dengan ENUM yang didefinisikan di database.
func (r MealTimeEnum) IsValid() bool {
	switch r {
	case MealTimeSarapan, MealTimeMakanSiang, MealTimeMakanMalam, MealTimeSelinganSiang, MealTimeSelinganSore:
		return true
	}
	return false
}

type Recipes struct {
	ID           uuid.UUID    `json:"recipe_id" db:"recipe_id"`
	DailyMenuID  uuid.UUID    `json:"daily_menu_id" db:"daily_menu_id"`
	Name         string       `json:"name" db:"name"`
	MealTime     MealTimeEnum `json:"meal_time" db:"meal_time"`
	MealTexture  string       `json:"meal_texture" db:"meal_texture"`
	Alergen      *string      `json:"alergen" db:"alergen"` // Enum (masih tidak jelas), dibuat pointer string agar bisa null
	Calories     int          `json:"callories" db:"callories"`
	Protein      int          `json:"protein" db:"protein"`
	Carbohydrate int          `json:"carbohydrate" db:"carbohydrate"`
	Fat          int          `json:"fat" db:"fat"`
	Description  string       `json:"description" db:"description"`
	CookingTime  int          `json:"cooking_time" db:"cooking_time"`
	IsBookmarked bool         `json:"is_bookmarked" db:"is_bookmarked"`
	IsDone       bool         `json:"is_done" db:"is_done"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`

	// Relation 1:N to MainIngredient and CookingStep
	MainIngredients []MainIngredient `json:"main_ingredients,omitempty" db:"-"`
	CookingSteps    []CookingStep    `json:"cooking_steps,omitempty" db:"-"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type ChildNutritionReport struct {
	ID                 uuid.UUID `json:"child_nutrition_report_id" db:"child_nutrition_report_id"`
	ChildID            uuid.UUID `json:"child_id" db:"child_id"`
	Calories           int       `json:"callories" db:"callories"`
	TargetCalories     int       `json:"target_callories" db:"target_callories"`
	Protein            int       `json:"protein" db:"protein"`
	TargetProtein      int       `json:"target_protein" db:"target_protein"`
	Fat                int       `json:"fat" db:"fat"`
	TargetFat          int       `json:"target_fat" db:"target_fat"`
	Carbohydrate       int       `json:"carbohydrate" db:"carbohydrate"`
	TargetCarbohydrate int       `json:"target_carbohydrate" db:"target_carbohydrate"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

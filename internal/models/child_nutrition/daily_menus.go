package models

import (
	"time"

	"github.com/google/uuid"
)

type DailyMenu struct {
	ID                     uuid.UUID `json:"daily_menu_id" db:"daily_menu_id"`
	ChildNutritionReportID uuid.UUID `json:"child_nutrition_report_id" db:"child_nutrition_report_id"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`

	// Relation 1:N to DailyMenuItem
	DailyMenuItems []Recipes `json:"daily_menu_items,omitempty" db:"-"`
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type ChildGrowthReport struct {
	ID                uuid.UUID `json:"child_growth_report_id" db:"child_growth_report_id"`
	ChildID           uuid.UUID `json:"child_id" db:"child_id"`
	Weight            float64   `json:"weight" db:"weight"`                         // (Kg)
	Height            float64   `json:"height" db:"height"`                         // (Cm)
	HeadCircumference float64   `json:"head_circumference" db:"head_circumference"` // (Cm)
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

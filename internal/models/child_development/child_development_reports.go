package models

import (
	"time"

	"github.com/google/uuid"
)

type ChildDevelopmentReport struct {
	ID            uuid.UUID       `json:"child_development_report_id" db:"child_development_report_id"`
	ChildID       uuid.UUID       `json:"child_id" db:"child_id"`
	KPSPScore     int             `json:"kpsp_score" db:"kpsp_score"`
	MonthTarget   MonthTargetEnum `json:"month_target" db:"month_target"`
	NextCheckDate time.Time       `json:"next_check_date" db:"next_check_date"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`

	// Relation 1:N
	Recommendations []RecommendedAction `json:"recommendations,omitempty" db:"-"`
}

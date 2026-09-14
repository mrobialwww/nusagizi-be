package models

import (
	"time"

	"github.com/google/uuid"
)

type CaregiverProfile struct {
	ID        uuid.UUID `json:"caregiver_profile_id" db:"caregiver_profile_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Relation 1:N
	Engagements []CaregiverEngagement `json:"caregiver_engagements,omitempty" db:"-"`
}

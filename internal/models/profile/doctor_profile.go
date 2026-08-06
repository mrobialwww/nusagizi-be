package models

import (
	"time"

	"github.com/google/uuid"
)

type DoctorProfile struct {
	ID           uuid.UUID `json:"doctor_profile_id" db:"doctor_profile_id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	FacilityName string    `json:"facility_name" db:"facility_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

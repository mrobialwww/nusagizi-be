package models

import (
	"time"

	"github.com/google/uuid"
)

type MotherProfile struct {
	ID        uuid.UUID `json:"mother_profile_id" db:"mother_profile_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Contact struct {
	ID              uuid.UUID `json:"contact_id" db:"contact_id"`
	MotherProfileID uuid.UUID `json:"mother_profile_id" db:"mother_profile_id"`
	ContactMotherID uuid.UUID `json:"contact_mother_id" db:"contact_mother_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type MedicalRelationship struct {
	ID              uuid.UUID  `json:"medical_relationship_id" db:"medical_relationship_id"`
	MotherProfileID uuid.UUID  `json:"mother_profile_id" db:"mother_profile_id"`
	DoctorProfileID uuid.UUID  `json:"doctor_profile_id" db:"doctor_profile_id"`
	DeletedAt       *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type CaregiverEngagement struct {
	ID                 uuid.UUID  `json:"caregiver_engagement_id" db:"caregiver_engagement_id"`
	CaregiverProfileID uuid.UUID  `json:"caregiver_profile_id" db:"caregiver_profile_id"`
	ChildID            uuid.UUID  `json:"child_id" db:"child_id"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

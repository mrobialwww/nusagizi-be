package models

import (
	"time"

	"github.com/google/uuid"
)

// Child mirrors the child table in the database.
type Child struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	MotherProfileID      uuid.UUID  `json:"mother_profile_id" db:"mother_profile_id"`
	FullName             string     `json:"full_name" db:"full_name"`
	Gender               string     `json:"gender" db:"gender"`
	BirthDate            time.Time  `json:"-" db:"birth_date"`
	PhotoURL             *string    `json:"photo_url" db:"photo_url"`
	NotesProfile         *string    `json:"notes" db:"notes_profile"`
	UploadStreakDays     int        `json:"upload_streak_days" db:"upload_streak_days"`
	DeletedAt            *time.Time `json:"-" db:"deleted_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}

// ChildInfo is a lightweight struct used for generating AI menu payloads without loading the entire Child model.
type ChildInfo struct {
	ID        uuid.UUID
	Nama      string
	Sex       string
	BirthDate time.Time
}

// GrowthPoint represents a single measurement point for a child's growth.
type GrowthPoint struct {
	MeasurementDate time.Time
	WeightKg        float64
	HeightCm        float64
	HeadCircCm      float64
}

// MedicalNote represents an active medical note for a child.
type MedicalNote struct {
	ID             uuid.UUID
	DoctorName     string
	Recommendation string
}

// NutritionTargets represents the daily targets for energy and protein from a medical note.
type NutritionTargets struct {
	EnergiKkal float64
	ProteinG   float64
}

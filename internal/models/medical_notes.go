package models

import (
	"time"

	"github.com/google/uuid"
)

type RestrictionType string

const (
	RestrictionTypePantangan RestrictionType = "PANTANGAN"
	RestrictionTypeAlergi    RestrictionType = "ALERGI"
)

// IsValid memastikan nilai RestrictionType valid sesuai dengan database
func (r RestrictionType) IsValid() bool {
	switch r {
	case RestrictionTypePantangan, RestrictionTypeAlergi:
		return true
	}
	return false
}

type MedicalNote struct {
	ID                    uuid.UUID `json:"medical_note_id" db:"medical_note_id"`
	MedicalRelationshipID uuid.UUID `json:"medical_relationship_id" db:"medical_relationship_id"`
	MedicalRecommendation string    `json:"medical_recommendation" db:"medical_recommendation"`
	ValidDate             time.Time `json:"valid_date" db:"valid_date"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`

	// Relation 1:N to MedicalAdvice, DailyNutritionTarget, MedicalRestriction
	Advices          []MedicalAdvice        `json:"medical_advices,omitempty" db:"-"`
	NutritionTargets []DailyNutritionTarget `json:"daily_nutrition_targets,omitempty" db:"-"`
	Restrictions     []MedicalRestriction   `json:"medical_restrictions,omitempty" db:"-"`
}

type MedicalAdvice struct {
	ID            uuid.UUID `json:"medical_advice_id" db:"medical_advice_id"`
	MedicalNoteID uuid.UUID `json:"medical_note_id" db:"medical_note_id"`
	Advice        string    `json:"advice" db:"advice"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type DailyNutritionTarget struct {
	ID            uuid.UUID `json:"daily_nutrition_target_id" db:"daily_nutrition_target_id"`
	MedicalNoteID uuid.UUID `json:"medical_note_id" db:"medical_note_id"`
	NutrientName  string    `json:"nutrient_name" db:"nutrient_name"`
	Quantity      int       `json:"quantity" db:"quantity"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type MedicalRestriction struct {
	ID            uuid.UUID       `json:"medical_restriction_id" db:"medical_restriction_id"`
	MedicalNoteID uuid.UUID       `json:"medical_note_id" db:"medical_note_id"`
	SubstanceName string          `json:"substance_name" db:"substance_name"`
	Type          RestrictionType `json:"type" db:"type"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

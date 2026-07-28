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
	FoodFrequencyProfile int        `json:"food_frequency" db:"food_frequency_profile"`
	FoodGoalProfile      string     `json:"food_goal" db:"food_goal_profile"`
	NotesProfile         *string    `json:"notes" db:"notes_profile"`
	UploadStreakDays     int        `json:"upload_streak_days" db:"upload_streak_days"`
	DeletedAt            *time.Time `json:"-" db:"deleted_at"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}


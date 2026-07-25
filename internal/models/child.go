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

// FavoriteFood mirrors the favorite_food_profile table.
type FavoriteFood struct {
	ID       uuid.UUID `db:"id"`
	ChildID  uuid.UUID `db:"child_id"`
	FoodName string    `db:"food_name"`
}

// FavoriteTexture mirrors the favorite_texture_profile table.
type FavoriteTexture struct {
	ID          uuid.UUID `db:"id"`
	ChildID     uuid.UUID `db:"child_id"`
	TextureName string    `db:"texture_name"`
}

// ChildDiet mirrors the child_diet_profile table.
type ChildDiet struct {
	ID       uuid.UUID `db:"id"`
	ChildID  uuid.UUID `db:"child_id"`
	DietName string    `db:"diet_name"`
}

// ChildChronicDisease mirrors the child_chronic_disease_profile table.
type ChildChronicDisease struct {
	ID          uuid.UUID `db:"id"`
	ChildID     uuid.UUID `db:"child_id"`
	DiseaseName string    `db:"disease_name"`
}

// ChildAllergy mirrors the child_allergy_profile table.
type ChildAllergy struct {
	ID           uuid.UUID `db:"id"`
	ChildID      uuid.UUID `db:"child_id"`
	Category     string    `db:"category"` // ENUM: food, medicine, animal, others
	AllergenName string    `db:"allergen_name"`
}

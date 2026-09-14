package models

import "time"

// User mirrors the users table in the database.
type User struct {
	ID          string     `json:"id" db:"id"`
	Auth0ID     string     `json:"-" db:"auth0_id"`
	Email       string     `json:"email" db:"email"`
	FullName    *string    `json:"full_name" db:"full_name"`
	Gender      *string    `json:"gender" db:"gender"`
	PhoneNumber *string    `json:"phone_number" db:"phone_number"`
	PhotoURL    *string    `json:"photo_url" db:"photo_url"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// UserProfileResponse is the response body
type UserProfileResponse struct {
	ID                  string  `json:"id"`
	FullName            *string `json:"full_name"`
	Gender              *string `json:"gender"`
	Email               string  `json:"email"`
	PhoneNumber         *string `json:"phone_number"`
	PhotoURL            *string `json:"photo_url"`
	HasMotherProfile    bool    `json:"has_mother_profile"`
	HasCaregiverProfile bool    `json:"has_caregiver_profile"`
}

// UpdateUserInput is the request body
// All fields are optional — only send fields that need to be changed.
type UpdateUserInput struct {
	FullName    *string `json:"full_name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Gender      *string `json:"gender"`
	PhotoURL    *string `json:"photo_url"`
}

// MotherProfileResponse is the response body
type MotherProfileResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	FullName    *string `json:"full_name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	PhotoURL    *string `json:"photo_url"`
}

// CaregiverProfileResponse is the response body
type CaregiverProfileResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	FullName    *string `json:"full_name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	PhotoURL    *string `json:"photo_url"`
}

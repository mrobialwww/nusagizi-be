package models

import "time"

// DateLayout is the date format used throughout the API (DD-MM-YYYY).
const DateLayout = "02-01-2006"

// Allergies represents the nested allergies structure in child profile requests/responses.
type Allergies struct {
	Food     []string `json:"food"`
	Medicine []string `json:"medicine"`
	Animal   []string `json:"animal"`
	Others   []string `json:"others"`
}

// CreateChildInput is the request body
type CreateChildInput struct {
	FullName         string     `json:"full_name" binding:"required,max=150"`
	BirthDate        string     `json:"birth_date" binding:"required"` // DD-MM-YYYY
	Gender           string     `json:"gender" binding:"required,oneof=male female"`
	PhotoURL         *string    `json:"photo_url"`
	Allergies        *Allergies `json:"allergies"`
	ChronicDiseases  []string   `json:"chronic_diseases"`
	Diets            []string   `json:"diets"`
	FavoriteFoods    []string   `json:"favorite_foods"`
	FavoriteTextures []string   `json:"favorite_textures"`
	Notes            *string    `json:"notes"`
}

// UpdateChildInput is the request body
// All fields are optional. Pointer slices distinguish "not sent" (nil) from "empty / clear all" ([]).
type UpdateChildInput struct {
	FullName         *string    `json:"full_name"`
	BirthDate        *string    `json:"birth_date"` // DD-MM-YYYY
	Gender           *string    `json:"gender"`
	PhotoURL         *string    `json:"photo_url"`
	Allergies        *Allergies `json:"allergies"`
	ChronicDiseases  *[]string  `json:"chronic_diseases"`
	Diets            *[]string  `json:"diets"`
	FavoriteFoods    *[]string  `json:"favorite_foods"`
	FavoriteTextures *[]string  `json:"favorite_textures"`
	Notes            *string    `json:"notes"`
}

// ChildSimpleResponse is the response body for
type ChildSimpleResponse struct {
	ID         string  `json:"id"`
	FullName   string  `json:"full_name"`
	BirthDate  string  `json:"birth_date"`
	Gender     string  `json:"gender"`
	PhotoURL   *string `json:"photo_url"`
	StreakDays int     `json:"streak_days"`
	MotherName string  `json:"mother_name"`
}

// ChildDetailResponse is the response body
type ChildDetailResponse struct {
	ID               string    `json:"id"`
	FullName         string    `json:"full_name"`
	BirthDate        string    `json:"birth_date"` // DD-MM-YYYY
	Gender           string    `json:"gender"`
	PhotoURL         *string   `json:"photo_url"`
	StreakDays       int       `json:"streak_days"`
	Allergies        Allergies `json:"allergies"`
	ChronicDiseases  []string  `json:"chronic_diseases"`
	Diets            []string  `json:"diets"`
	FavoriteFoods    []string  `json:"favorite_foods"`
	FavoriteTextures []string  `json:"favorite_textures"`
	Notes            *string   `json:"notes"`
}

// ChildListItem is one entry in the list response
type ChildListItem struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	BirthDate time.Time `json:"-"`
	Age       string    `json:"age"`
	Gender    string  `json:"gender"`
	PhotoURL  *string `json:"photo_url"`
}

package models

// DateLayout is the date format used throughout the API (DD-MM-YYYY).
const DateLayout = "02-01-2006"

// Allergies represents the nested allergies structure in child profile requests/responses.
type Allergies struct {
	Food     []string `json:"food"`
	Medicine []string `json:"medicine"`
	Animal   []string `json:"animal"`
	Others   []string `json:"others"`
}

// CreateChildInput is the request body for POST /children (endpoint 22).
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
	FoodFrequency    int        `json:"food_frequency" binding:"required,min=1"`
	FoodGoal         string     `json:"food_goal" binding:"required"`
	Notes            *string    `json:"notes"`
}

// UpdateChildInput is the request body for PATCH /children/{child_id} (endpoint 24).
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
	FoodFrequency    *int       `json:"food_frequency"`
	FoodGoal         *string    `json:"food_goal"`
	Notes            *string    `json:"notes"`
}

// ChildDetailResponse is the response body for GET /children/{child_id} (endpoint 50).
type ChildDetailResponse struct {
	ID               string    `json:"id"`
	FullName         string    `json:"full_name"`
	BirthDate        string    `json:"birth_date"` // DD-MM-YYYY
	Gender           string    `json:"gender"`
	PhotoURL         *string   `json:"photo_url"`
	UploadStreakDays int       `json:"upload_streak_days"`
	Allergies        Allergies `json:"allergies"`
	ChronicDiseases  []string  `json:"chronic_diseases"`
	Diets            []string  `json:"diets"`
	FavoriteFoods    []string  `json:"favorite_foods"`
	FavoriteTextures []string  `json:"favorite_textures"`
	FoodFrequency    int       `json:"food_frequency"`
	FoodGoal         string    `json:"food_goal"`
	Notes            *string   `json:"notes"`
}

// ChildListItem is one entry in the list response for GET /mother-profiles/{id}/children (endpoint 51).
type ChildListItem struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	BirthDate string  `json:"birth_date"` // DD-MM-YYYY
	Gender    string  `json:"gender"`
	PhotoURL  *string `json:"photo_url"`
}

package medical

import (
	"time"

	"github.com/google/uuid"
)

type NutrientType string

const (
	NutrientCalories     NutrientType = "calorie"
	NutrientProtein      NutrientType = "protein"
	NutrientFat          NutrientType = "fat"
	NutrientCarbohydrate NutrientType = "carbohydrate"
)

type MedicalNoteResponse struct {
	ID               uuid.UUID `json:"id"`
	ChildName        string    `json:"child_name"`
	Recommendation   string    `json:"recommendation"`
	ValidUntil       string    `json:"valid_until"`
	CreatedAt        time.Time `json:"created_at"`
	ProhibitionCount int       `json:"prohibition_count"`
	AllergyCount     int       `json:"allergy_count"`
}

type DailyNutritionTarget struct {
	Nutrient NutrientType `json:"nutrient"`
	Quantity float64      `json:"quantity"`
}

type MedicalNoteDetailResponse struct {
	ID                    uuid.UUID              `json:"id"`
	ChildName             string                 `json:"child_name"`
	DoctorName            string                 `json:"doctor_name"`
	FacilityName          *string                `json:"facility_name"`
	Recommendation        string                 `json:"recommendation"`
	ValidUntil            string                 `json:"valid_until"`
	CreatedAt             time.Time              `json:"created_at"`
	Prohibitions          []string               `json:"prohibitions"`
	Allergies             []string               `json:"allergies"`
	DailyNutritionTargets []DailyNutritionTarget `json:"daily_nutrition_targets"`
}

type CreateMedicalNoteInput struct {
	ChildID               uuid.UUID              `json:"child_id"`
	DoctorName            string                 `json:"doctor_name"`
	FacilityName          *string                `json:"facility_name"`
	Recommendation        string                 `json:"recommendation"`
	DailyNutritionTargets []DailyNutritionTarget `json:"daily_nutrition_targets"`
	Prohibitions          []string               `json:"prohibitions"`
	Allergies             []string               `json:"allergies"`
	ValidUntil            string                 `json:"valid_until"`
}

type UpdateMedicalNoteInput struct {
	DoctorName            *string                 `json:"doctor_name"`
	FacilityName          *string                 `json:"facility_name"`
	Recommendation        *string                 `json:"recommendation"`
	DailyNutritionTargets *[]DailyNutritionTarget `json:"daily_nutrition_targets"`
	Prohibitions          *[]string               `json:"prohibitions"`
	Allergies             *[]string               `json:"allergies"`
	ValidUntil            *string                 `json:"valid_until"`
}

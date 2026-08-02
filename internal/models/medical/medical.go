package medical

import (
	"time"

	"github.com/google/uuid"
)

type NutrientType string

const (
	NutrientCalories     NutrientType = "calories"
	NutrientProtein      NutrientType = "protein"
	NutrientFat          NutrientType = "fat"
	NutrientCarbohydrate NutrientType = "carbohydrate"
)

type MedicalNoteResponse struct {
	ID               uuid.UUID `json:"id"`
	ChildName        string    `json:"child_name"`
	Recommendation   string    `json:"recommendation"`
	ValidDate        string    `json:"valid_date"`
	CreatedAt        time.Time `json:"created_at"`
	ProhibitionCount int       `json:"prohibition_count"`
	AllergyCount     int       `json:"allergy_count"`
}

type DailyNutritionTarget struct {
	NutrientName NutrientType `json:"nutrient_name"`
	Quantity     float64      `json:"quantity"`
}

type MedicalNoteDetailResponse struct {
	ID                    uuid.UUID              `json:"id"`
	ChildName             string                 `json:"child_name"`
	DoctorName            string                 `json:"doctor_name"`
	FacilityLocation      *string                `json:"facility_location"`
	Recommendation        string                 `json:"recommendation"`
	ValidDate             string                 `json:"valid_date"`
	CreatedAt             time.Time              `json:"created_at"`
	Prohibitions          []string               `json:"prohibitions"`
	Allergies             []string               `json:"allergies"`
	DailyNutritionTargets []DailyNutritionTarget `json:"daily_nutrition_targets"`
}

type CreateMedicalNoteInput struct {
	ChildID               uuid.UUID              `json:"child_id"`
	DoctorName            string                 `json:"doctor_name"`
	FacilityLocation      *string                `json:"facility_location"`
	Recommendation        string                 `json:"recommendation"`
	DailyNutritionTargets []DailyNutritionTarget `json:"daily_nutrition_targets"`
	Prohibitions          []string               `json:"prohibitions"`
	Allergies             []string               `json:"allergies"`
	ValidDate             string                 `json:"valid_date"`
}

type UpdateMedicalNoteInput struct {
	DoctorName            *string                 `json:"doctor_name"`
	FacilityLocation      *string                 `json:"facility_location"`
	Recommendation        *string                 `json:"recommendation"`
	DailyNutritionTargets *[]DailyNutritionTarget `json:"daily_nutrition_targets"`
	Prohibitions          *[]string               `json:"prohibitions"`
	Allergies             *[]string               `json:"allergies"`
	ValidDate             *string                 `json:"valid_date"`
}

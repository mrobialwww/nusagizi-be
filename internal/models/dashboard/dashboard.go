package dashboard

import (
	"time"

	"github.com/google/uuid"
)

type ChildSummaryResponse struct {
	ID                  uuid.UUID `json:"id"`
	FullName            string    `json:"full_name"`
	PhotoUrl            *string   `json:"photo_url"`
	Age                 string    `json:"age"`
	BirthDateRaw        time.Time `json:"-"`
	Gender              string    `json:"-"`
	HeightCm            *float64  `json:"height_cm"`
	WeightKg            *float64  `json:"weight_kg"`
	HeadCircumferenceCm *float64  `json:"-"`
	GrowthMeasuredAt    *time.Time `json:"-"`
	KpspScore           *int      `json:"kpsp_score"`
	KpspAnswersCount    *int      `json:"kpsp_answers_count"`
	Calories            *int      `json:"-"`
	TargetCalories      *int      `json:"-"`
	Protein             *int      `json:"protein"`
	TargetProtein       *int      `json:"target_protein"`
	Fat                 *int      `json:"-"`
	TargetFat           *int      `json:"-"`
	Carbohydrate        *int      `json:"-"`
	TargetCarbohydrate  *int      `json:"-"`
	Streak              int       `json:"streak"`
	StatusGrowth        string    `json:"status_growth"`
	StatusDevelopment   string    `json:"status_development"`
	StatusNutrition     string    `json:"status_nutrition"`
	Status              string    `json:"status"`
}

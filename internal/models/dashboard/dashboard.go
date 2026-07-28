package dashboard

import (
	"github.com/google/uuid"
)

type ChildSummaryResponse struct {
	ID            uuid.UUID `json:"id"`
	FullName      string    `json:"full_name"`
	BirthDate     string    `json:"birth_date"`
	HeightCm      *float64  `json:"height_cm"`
	WeightKg      *float64  `json:"weight_kg"`
	KpspScore     *int      `json:"kpsp_score"`
	Protein       *int      `json:"protein"`
	TargetProtein *int      `json:"target_protein"`
	Streak        int       `json:"streak"`
}

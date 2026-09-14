package models

import (
	"time"

	"github.com/google/uuid"
)

type ChildGrowthReport struct {
	ID                  uuid.UUID `json:"id" db:"id"`
	MeasuredAt          time.Time `json:"measured_at" db:"measured_at"` // As time.Time, formatted to string in response
	WeightKg            *float64  `json:"weight_kg" db:"weight_kg"`
	HeightCm            *float64  `json:"height_cm" db:"height_cm"`
	HeadCircumferenceCm *float64  `json:"head_circumference_cm" db:"head_circumference_cm"`
}

type CreateGrowthReportInput struct {
	MeasuredAt          string   `json:"measured_at"` // Format DD-MM-YYYY
	WeightKg            *float64 `json:"weight_kg"`
	HeightCm            *float64 `json:"height_cm"`
	HeadCircumferenceCm *float64 `json:"head_circumference_cm"`
}

type UpdateGrowthReportInput struct {
	MeasuredAt          *string  `json:"measured_at"` // Format DD-MM-YYYY
	WeightKg            *float64 `json:"weight_kg"`
	HeightCm            *float64 `json:"height_cm"`
	HeadCircumferenceCm *float64 `json:"head_circumference_cm"`
}

type GrowthDataPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type GrowthMessage struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GrowthAnalysisResult struct {
	DataPoints []GrowthDataPoint `json:"data_points"`
	Message    *GrowthMessage    `json:"message"`
}

type RawGrowthMeasurement struct {
	WeightKg            float64
	HeightCm            float64
	HeadCircumferenceCm float64
	AgeMonths           int
	Gender              string
}

type RawGrowthReportFull struct {
	ID                  uuid.UUID
	MeasuredAt          time.Time
	WeightKg            float64
	HeightCm            float64
	HeadCircumferenceCm float64
	AgeMonths           int
	Gender              string
}

type GrowthReportWithStatus struct {
	ID                  uuid.UUID `json:"id"`
	MeasuredAt          time.Time `json:"measured_at"`
	WeightKg            float64   `json:"weight_kg"`
	HeightCm            float64   `json:"height_cm"`
	HeadCircumferenceCm float64   `json:"head_circumference_cm"`
	Status              string    `json:"status"`
	Description         string    `json:"description"`
}

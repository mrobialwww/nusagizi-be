package models

import (
	"time"

	"github.com/google/uuid"
)

// RecommendedAction: master data rekomendasi yang terkait dengan pertanyaan KPSP tertentu
type RecommendedAction struct {
	ID         uuid.UUID `json:"recommended_action_id" db:"recommended_action_id"`
	Title      string    `json:"title" db:"title"`
	ActionText string    `json:"action_text" db:"action_text"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type AssessmentKPSPQuestion struct {
	ID                  uuid.UUID               `json:"id" db:"id"`
	DevelopmentalDomain DevelopmentalDomainEnum `json:"developmental_domain" db:"developmental_domain"`
	MonthTarget         int                     `json:"month_target" db:"month_target"`
	QuestionText        string                  `json:"question_text" db:"question_text"`
	ImageURL            *string                 `json:"image_url" db:"image_url"`
	CreatedAt           time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at" db:"updated_at"`
}

type RecommendationItem struct {
	DevelopmentalDomain string `json:"developmental_domain"`
	ActionText          string `json:"action_text"`
}

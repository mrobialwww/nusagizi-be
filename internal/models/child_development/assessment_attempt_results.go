package models

import (
	"time"

	"github.com/google/uuid"
)

// RecommendedAction: master data rekomendasi yang terkait dengan pertanyaan KPSP tertentu
type RecommendedAction struct {
	ID                       uuid.UUID `json:"recommended_action_id" db:"recommended_action_id"`
	AssessmentKPSPQuestionID uuid.UUID `json:"assessment_kpsp_question_id" db:"assessment_kpsp_question_id"`
	Title                    string    `json:"title" db:"title"`
	ActionText               string    `json:"action_text" db:"action_text"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
}

type AssessmentKPSPQuestion struct {
	ID                  uuid.UUID               `json:"assessment_kpsp_question_id" db:"assessment_kpsp_question_id"`
	Order               int                     `json:"order" db:"order"`
	DevelopmentalDomain DevelopmentalDomainEnum `json:"developmental_domain" db:"developmental_domain"`
	MonthTarget         MonthTargetEnum         `json:"month_target" db:"month_target"`
	Question            string                  `json:"question" db:"question"`
	Description         string                  `json:"description" db:"description"`
}

type RecommendationItem struct {
	DevelopmentalDomain string `json:"developmental_domain"`
	ActionText          string `json:"action_text"`
}

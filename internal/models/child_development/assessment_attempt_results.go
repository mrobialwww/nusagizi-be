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
	ID          uuid.UUID       `json:"assessment_kpsp_question_id" db:"assessment_kpsp_question_id"`
	Order       int             `json:"order" db:"order"`
	NerveName   NerveNameEnum   `json:"nerve_name" db:"nerve_name"`
	MonthTarget MonthTargetEnum `json:"month_target" db:"month_target"`
	Question    string          `json:"question" db:"question"`
	Description string          `json:"description" db:"description"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

type RecommendationItem struct {
	NerveName  string `json:"nerve_name"`
	ActionText string `json:"action_text"`
}

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
	Description              string    `json:"description" db:"description"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
}

type AssessmentKPSPAnswer struct {
	ID                       uuid.UUID `json:"assessment_kpsp_answer_id" db:"assessment_kpsp_answer_id"`
	ChildDevelopmentReportID uuid.UUID `json:"child_development_report_id" db:"child_development_report_id"`
	AssessmentKPSPQuestionID uuid.UUID `json:"assessment_kpsp_question_id" db:"assessment_kpsp_question_id"`
	Answer                   bool      `json:"assessment_kpsp_answer" db:"assessment_kpsp_answer"`
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

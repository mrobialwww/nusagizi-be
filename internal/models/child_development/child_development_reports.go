package models

import (
	"time"

	"github.com/google/uuid"
)

type DomainAggregate struct {
	NerveName     string `json:"nerve_name"`
	TotalQuestion int    `json:"total_question"`
	TrueAnswer    int    `json:"true_answer"`
}


type DevelopmentReportResponse struct {
	ID                 uuid.UUID           `json:"id"`
	KPSPScore          int                 `json:"kpsp_score"`
	MonthTarget        *int                `json:"month_target,omitempty"`
	NextCheckDate      *string             `json:"next_check_date,omitempty"`
	Status             *string             `json:"status,omitempty"`
	CreatedAt          *time.Time          `json:"created_at,omitempty"`
	Domains            []DomainAggregate   `json:"domains,omitempty"`
	RecommendedActions []RecommendedAction `json:"recommended_actions,omitempty"`
}

type KPSPAnswerInput struct {
	AssessmentKPSPQuestionID uuid.UUID `json:"assessment_kpsp_question_id"`
	Answer                   bool      `json:"answer"` // Also matches "assessment_kpsp_answer" for backwards compatibility if needed
}

type CreateDevelopmentReportInput struct {
	MonthTarget int               `json:"month_target"`
	ListAnswer  []KPSPAnswerInput `json:"list_answer"`
}

type UpdateDevelopmentReportInput struct {
	ListAnswer []KPSPAnswerInput `json:"list_answer"`
}

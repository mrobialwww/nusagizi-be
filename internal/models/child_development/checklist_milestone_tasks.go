package models

import (
	"github.com/google/uuid"
)

type ChecklistMilestoneResponse struct {
	ID                  uuid.UUID `json:"id"`
	DevelopmentalDomain string    `json:"developmental_domain"`
	QuestionText        string    `json:"question_text"`
	IsChecked           bool      `json:"is_checked"`
}

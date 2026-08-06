package models

import (
	"github.com/google/uuid"
)

type ChecklistMilestoneTaskResponse struct {
	ID                  uuid.UUID `json:"id"`
	DevelopmentalDomain string    `json:"developmental_domain"`
	TaskDescription     string    `json:"task_description"`
	IsChecked           bool      `json:"is_checked"`
}

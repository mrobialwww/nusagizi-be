package models

import (


	"github.com/google/uuid"
)

type ChecklistMilestoneTaskResponse struct {
	ID                 uuid.UUID `json:"id"`
	NerveName          string    `json:"nerve_name"`
	TaskDescription    string    `json:"task_description"`
	IsChecked          bool      `json:"is_checked"`
}


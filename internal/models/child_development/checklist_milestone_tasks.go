package models

import (
	"time"

	"github.com/google/uuid"
)

type ChecklistMilestoneTask struct {
	ID                 uuid.UUID       `json:"checklist_milestone_task_id" db:"checklist_milestone_task_id"`
	MonthTarget        MonthTargetEnum `json:"month_target" db:"month_target"`
	NerveName          NerveNameEnum   `json:"nerve_name" db:"nerve_name"`
	CheckedDescription string          `json:"checked_description" db:"checked_description"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at" db:"updated_at"`
}

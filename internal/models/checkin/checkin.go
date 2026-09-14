package checkin

import (
	"time"

	"github.com/google/uuid"
)

type TokenData struct {
	Token     string    `json:"token"`
	ChildID   uuid.UUID `json:"childId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type CheckinLog struct {
	ID          uuid.UUID `json:"-"`
	Token       string    `json:"-"`
	ChildID     uuid.UUID `json:"-"`
	CaregiverID uuid.UUID `json:"caregiverId"`
	CheckedInAt time.Time `json:"checkedInAt"`
}

type GenerateReq struct {
	ChildID uuid.UUID `json:"childId" binding:"required"`
}

type ValidateReq struct {
	Token       string    `json:"token" binding:"required"`
	CaregiverID uuid.UUID `json:"caregiverId" binding:"required"`
}

type ValidateResult struct {
	Valid        bool      `json:"valid"`
	EngagementID string    `json:"engagementId"`
	CheckedInAt  time.Time `json:"checkedInAt"`
	ChildName    string    `json:"childName"`
	ChildAge     string    `json:"childAge"`
	MotherName   string    `json:"motherName"`
}

type ChildCheckinInfo struct {
	ChildName  string
	BirthDate  time.Time
	MotherName string
}

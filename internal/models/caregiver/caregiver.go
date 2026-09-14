package caregiver

import (
	"github.com/google/uuid"
)

type CaregiverEngagementResponse struct {
	CaregiverEngagementID uuid.UUID `json:"caregiver_engagement_id"`
	ChildID               uuid.UUID `json:"child_id"`
	ChildName             string    `json:"child_name"`
	CaregiverProfileID    uuid.UUID `json:"caregiver_profile_id"`
	CaregiverName         string    `json:"caregiver_name"`
	PhoneNumber           *string   `json:"phone_number"`
}

type CaregiverProfileResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber *string   `json:"phone_number"`
	PhotoURL    *string   `json:"photo_url"`
}

type ChildSimpleResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	PhotoURL *string   `json:"photo_url"`
}

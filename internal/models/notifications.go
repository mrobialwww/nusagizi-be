package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationTypeEnum string

const (
	MealReminder       NotificationTypeEnum = "MEAL_REMINDER"
	GrowthDevelopment  NotificationTypeEnum = "GROWTH_DEVELOPMENT_REMINDER"
	DoctorActivity     NotificationTypeEnum = "DOCTOR_ACTIVITY"
	PhotosActivity     NotificationTypeEnum = "PHOTOS_ACTIVITY"
	ShopActivity       NotificationTypeEnum = "SHOP_ACTIVITY"
	InvitationActivity NotificationTypeEnum = "INVITATION_ACTIVITY"
	NewMenuActivity    NotificationTypeEnum = "NEW_MENU_REMINDER"
	Achievement        NotificationTypeEnum = "ACHIEVEMENT"
)

type Notification struct {
	ID        uuid.UUID            `json:"notification_id" db:"notification_id"`
	UserID    uuid.UUID            `json:"user_id" db:"user_id"`
	Type      NotificationTypeEnum `json:"type" db:"type"`
	Title     string               `json:"title" db:"title"`
	Message   string               `json:"message" db:"message"`
	CreatedAt time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" db:"updated_at"`
}

package notification

import (
	"time"
	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Message          string    `json:"message"`
	NotificationType string    `json:"notification_type"`
	CreatedAt        time.Time `json:"created_at"`
}

package services

import (
	"context"
	"fmt"

	"nusagizi_be/internal/models/notification"
	"nusagizi_be/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// GetNotifications (Endpoint: 64)
func (s *NotificationService) GetNotifications(ctx context.Context, authUserID string, notificationType *string) ([]notification.NotificationResponse, error) {
	if notificationType != nil {
		validTypes := map[string]bool{
			"meal_reminder":               true,
			"growth_development_reminder": true,
			"doctor_activity":             true,
			"photos_activity":             true,
			"shop_activity":               true,
			"invitation_activity":         true,
			"new_menu_reminder":           true,
			"achievement":                 true,
		}
		if !validTypes[*notificationType] {
			return nil, fmt.Errorf("invalid notification_type")
		}
	}

	return s.repo.GetNotifications(ctx, authUserID, notificationType)
}

package services

import (
	"context"

	"nusagizi_be/internal/models/notification"
	"nusagizi_be/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// GetNotifications (Endpoint: 66)
func (s *NotificationService) GetNotifications(ctx context.Context, authUserID string) ([]notification.NotificationResponse, error) {
	return s.repo.GetNotifications(ctx, authUserID)
}

// GetLatestNotification returns the most recent notification for a user. (Endpoint: 67)
func (s *NotificationService) GetLatestNotification(ctx context.Context, userID string) (*notification.NotificationResponse, error) {
	return s.repo.GetLatestNotification(ctx, userID)
}

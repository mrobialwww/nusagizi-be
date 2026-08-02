package repository

import (
	"context"
	"errors"
	"time"

	"nusagizi_be/internal/models/notification"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

// GetNotifications returns all notifications for a user.
func (r *NotificationRepository) GetNotifications(ctx context.Context, userID string) ([]notification.NotificationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT id, title, message, notification_type, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []notification.NotificationResponse
	for rows.Next() {
		var res notification.NotificationResponse
		if err := rows.Scan(&res.ID, &res.Title, &res.Message, &res.NotificationType, &res.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, rows.Err()
}

// GetLatestNotification returns the most recent notification for a user.
func (r *NotificationRepository) GetLatestNotification(ctx context.Context, userID string) (*notification.NotificationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT id, title, message, notification_type, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var res notification.NotificationResponse
	err := r.pool.QueryRow(ctx, query, userID).Scan(&res.ID, &res.Title, &res.Message, &res.NotificationType, &res.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &res, nil
}

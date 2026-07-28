package repository

import (
	"context"
	"time"

	"nusagizi_be/internal/models/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

// GetNotifications returns notifications for a user, optionally filtered by type.
func (r *NotificationRepository) GetNotifications(ctx context.Context, userID string, notificationType *string) ([]notification.NotificationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT id, title, message, notification_type, created_at
		FROM notification
		WHERE user_id = $1
	`
	args := []interface{}{userID}

	if notificationType != nil {
		query += ` AND notification_type = $2`
		args = append(args, *notificationType)
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query, args...)
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

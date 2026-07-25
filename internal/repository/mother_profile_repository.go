package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetMotherProfileByUserID returns the mother_profile ID for a given users.id.
// Used internally by the child service to obtain the mother's profile ID from a user ID.
func GetMotherProfileByUserID(pool *pgxpool.Pool, userID string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var motherProfileID uuid.UUID
	var query string = `
		SELECT id 
		FROM mother_profiles 
		WHERE user_id = $1`

	err := pool.QueryRow(ctx, query, userID).Scan(&motherProfileID)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrNotFound
	}
	return motherProfileID, err
}

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MotherProfileRepository struct {
	pool *pgxpool.Pool
}

func NewMotherProfileRepository(pool *pgxpool.Pool) *MotherProfileRepository {
	return &MotherProfileRepository{pool: pool}
}

// GetByUserID returns the mother_profile ID for a given users.id.
// Used internally by the child service to obtain the mother's profile ID from a user ID.
func (r *MotherProfileRepository) GetByUserID(ctx context.Context, userID string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var motherProfileID uuid.UUID
	var query string = `
		SELECT id 
		FROM mother_profiles 
		WHERE user_id = $1`

	err := r.pool.QueryRow(ctx, query, userID).Scan(&motherProfileID)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrNotFound
	}
	return motherProfileID, err
}

// Create inserts a new mother_profile for the given user ID.
func (r *MotherProfileRepository) Create(ctx context.Context, userID string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	newID := uuid.New()
	var query string = `
		INSERT INTO mother_profiles (id, user_id) 
		VALUES ($1, $2)
	`
	_, err := r.pool.Exec(ctx, query, newID, userID)
	if err != nil {
		return uuid.Nil, err
	}
	return newID, nil
}

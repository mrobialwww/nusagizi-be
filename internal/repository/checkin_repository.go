package repository

import (
	"context"
	"fmt"
	"nusagizi_be/internal/models/checkin"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckinRepository struct {
	pool *pgxpool.Pool
}

func NewCheckinRepository(pool *pgxpool.Pool) *CheckinRepository {
	return &CheckinRepository{pool: pool}
}

func (r *CheckinRepository) SaveToken(ctx context.Context, token string, childID uuid.UUID, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO checkin_tokens (token, child_id, expires_at) 
		VALUES ($1, $2, $3)`

	_, err := r.pool.Exec(ctx, query, token, childID, expiresAt)
	return err
}

// DeleteExpiredTokens menghapus semua token yang sudah melewati batas waktu kedaluwarsa.
func (r *CheckinRepository) DeleteExpiredTokens(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		DELETE 
		FROM checkin_tokens 
		WHERE expires_at < NOW()`
	_, err := r.pool.Exec(ctx, query)
	return err
}

func (r *CheckinRepository) GetToken(ctx context.Context, token string) (*checkin.TokenData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT token, child_id, expires_at 
		FROM checkin_tokens 
		WHERE token = $1`

	var data checkin.TokenData
	err := r.pool.QueryRow(ctx, query, token).Scan(
		&data.Token,
		&data.ChildID,
		&data.ExpiresAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &data, nil
}

func (r *CheckinRepository) SaveLog(ctx context.Context, token string, childID, caregiverID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO checkin_logs (token, child_id, caregiver_id) 
		VALUES ($1, $2, $3)`

	_, err := r.pool.Exec(ctx, query, token, childID, caregiverID)
	return err
}

// CreateCaregiverEngagement creates a new active engagement.
func (r *CheckinRepository) CreateCaregiverEngagement(ctx context.Context, childID, caregiverProfileID uuid.UUID) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if active engagement already exists
	var exists bool
	queryCheck := `
		SELECT EXISTS(
			SELECT 1 
			FROM caregiver_engagements 
			WHERE child_id = $1 
				AND caregiver_profile_id = $2 
				AND deleted_at IS NULL
		)
	`
	err := r.pool.QueryRow(ctx, queryCheck, childID, caregiverProfileID).Scan(&exists)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, fmt.Errorf("active engagement already exists")
	}

	// Create engagement
	newID := uuid.New()
	queryInsert := `
		INSERT INTO caregiver_engagements (id, child_id, caregiver_profile_id)
		VALUES ($1, $2, $3)
	`
	_, err = r.pool.Exec(ctx, queryInsert, newID, childID, caregiverProfileID)
	if err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}

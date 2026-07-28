package repository

import (
	"context"
	"fmt"
	"time"

	"nusagizi_be/internal/models/caregiver"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CaregiverRepository struct {
	pool *pgxpool.Pool
}

func NewCaregiverRepository(pool *pgxpool.Pool) *CaregiverRepository {
	return &CaregiverRepository{pool: pool}
}

// Create inserts a new caregiver_profiles for the given user ID.
func (r *CaregiverRepository) Create(ctx context.Context, userID string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	newID := uuid.New()
	var query string = `
		INSERT INTO caregiver_profiles (id, user_id) 
		VALUES ($1, $2)
	`
	_, err := r.pool.Exec(ctx, query, newID, userID)
	if err != nil {
		return uuid.Nil, err
	}
	return newID, nil
}

// GetByUserID returns the caregiver_profile ID for a given user ID.
func (r *CaregiverRepository) GetByUserID(ctx context.Context, userID string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var profileID uuid.UUID
	query := `
		SELECT id 
		FROM caregiver_profiles
		WHERE user_id = $1`

	err := r.pool.QueryRow(ctx, query, userID).Scan(&profileID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("record not found") // Maps to ErrNotFound if we had one here, but we'll return pgx.ErrNoRows or standard error
		}
		return uuid.Nil, err
	}
	return profileID, nil
}

// CreateCaregiverEngagement creates a new active engagement.
func (r *CaregiverRepository) CreateCaregiverEngagement(ctx context.Context, childID, caregiverProfileID uuid.UUID) (uuid.UUID, error) {
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

// GetCaregiverEngagements gets a list of engagements (active or revoked).
func (r *CaregiverRepository) GetCaregiverEngagements(ctx context.Context, motherProfileID uuid.UUID, isRevoked bool) ([]caregiver.CaregiverEngagementResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			c.id AS child_id,
			c.full_name AS child_name,
			ce.id AS caregiver_engagement_id,
			ce.caregiver_profile_id,
			u.full_name AS caregiver_name,
			u.phone_number AS phone_number
		FROM child c
		JOIN caregiver_engagements ce ON ce.child_id = c.id
		JOIN caregiver_profiles cp ON cp.id = ce.caregiver_profile_id
		JOIN users u ON u.id = cp.user_id
		WHERE c.mother_profile_id = $1
	`

	if isRevoked {
		query += ` AND ce.deleted_at IS NOT NULL`
	} else {
		query += ` AND ce.deleted_at IS NULL`
	}

	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []caregiver.CaregiverEngagementResponse
	for rows.Next() {
		var res caregiver.CaregiverEngagementResponse
		if err := rows.Scan(&res.ChildID, &res.ChildName, &res.CaregiverEngagementID, &res.CaregiverProfileID, &res.CaregiverName, &res.PhoneNumber); err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, rows.Err()
}

// DeleteCaregiverEngagement performs a soft delete.
func (r *CaregiverRepository) DeleteCaregiverEngagement(ctx context.Context, engagementID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE caregiver_engagements 
		SET deleted_at = now() 
		WHERE id = $1 
			AND deleted_at IS NULL`
	res, err := r.pool.Exec(ctx, query, engagementID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found or already deleted")
	}
	return nil
}

// GetCaregiverProfile returns profile details including user table join.
func (r *CaregiverRepository) GetCaregiverProfile(ctx context.Context, caregiverProfileID uuid.UUID) (*caregiver.CaregiverProfileResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var profile caregiver.CaregiverProfileResponse
	query := `
		SELECT 
			cp.id, cp.user_id, u.full_name, u.email, u.phone_number, u.photo_url
		FROM caregiver_profiles cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.id = $1
	`
	err := r.pool.QueryRow(ctx, query, caregiverProfileID).Scan(
		&profile.ID, &profile.UserID, &profile.FullName, &profile.Email, &profile.PhoneNumber, &profile.PhotoURL,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}
	return &profile, nil
}

// GetCaregiverChildren gets children associated with a caregiver engagement.
func (r *CaregiverRepository) GetCaregiverChildren(ctx context.Context, caregiverProfileID uuid.UUID) ([]caregiver.ChildSimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT c.id, c.full_name, c.photo_url
		FROM caregiver_engagements ce
		JOIN child c ON c.id = ce.child_id
		WHERE ce.caregiver_profile_id = $1 
			AND ce.deleted_at IS NULL
	`
	rows, err := r.pool.Query(ctx, query, caregiverProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []caregiver.ChildSimpleResponse
	for rows.Next() {
		var c caregiver.ChildSimpleResponse
		if err := rows.Scan(&c.ID, &c.FullName, &c.PhotoURL); err != nil {
			return nil, err
		}
		children = append(children, c)
	}
	return children, rows.Err()
}

// CheckCaregiverAccess returns true if the child has an active engagement with the caregiver.
func (r *CaregiverRepository) CheckCaregiverAccess(ctx context.Context, childID, caregiverProfileID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool
	var query string = `
		SELECT EXISTS(
			SELECT 1 FROM caregiver_engagements
			WHERE child_id = $1 
				AND caregiver_profile_id = $2 
				AND deleted_at IS NULL
		)
	`
	err := r.pool.QueryRow(ctx, query, childID, caregiverProfileID).Scan(&exists)
	return exists, err
}

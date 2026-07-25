package repository

import (
	"context"
	"errors"
	"fmt"
	"nusagizi_be/internal/models"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetUserBySub returns a user by Auth0 sub (auth0_id).
// Used by the auth middleware for lazy provisioning.
func GetUserBySub(pool *pgxpool.Pool, sub string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u models.User
	var query string = `SELECT id, auth0_id, email, full_name, gender, phone_number, photo_url, deleted_at, created_at, updated_at 
		FROM users 
		WHERE auth0_id = $1 AND deleted_at IS NULL`

	err := pool.QueryRow(ctx, query, sub).Scan(
		&u.ID,
		&u.Auth0ID,
		&u.Email,
		&u.FullName,
		&u.Gender,
		&u.PhoneNumber,
		&u.PhotoURL,
		&u.DeletedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// CreateUserFromAuth0 inserts a new user on first Auth0 login (lazy provisioning).
// The username parameter is accepted for backward compatibility but is not stored
// (there is no username column in the users table).
func CreateUserFromAuth0(pool *pgxpool.Pool, sub, email, username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u models.User
	var query string = `
		INSERT INTO users (auth0_id, email)
		VALUES ($1, $2)
		ON CONFLICT (auth0_id) DO UPDATE SET updated_at = now()
		RETURNING id, auth0_id, email, full_name, gender, phone_number, photo_url, deleted_at, created_at, updated_at`
	err := pool.QueryRow(ctx, query, sub, email).Scan(
		&u.ID,
		&u.Auth0ID,
		&u.Email,
		&u.FullName,
		&u.Gender,
		&u.PhoneNumber,
		&u.PhotoURL,
		&u.DeletedAt,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserProfileByID returns a user with has_mother_profile and has_caregiver_profile flags.
// Used by endpoint 45 (GET /users/{user_id}).
func GetUserProfileByID(pool *pgxpool.Pool, userID string) (*models.UserProfileResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var resp models.UserProfileResponse
	var query string = `
		SELECT
			u.id,
			u.email,
			u.full_name,
			u.gender,
			u.phone_number,
			u.photo_url,
			EXISTS(SELECT 1 FROM mother_profiles    WHERE user_id = u.id) AS has_mother_profile,
			EXISTS(SELECT 1 FROM caregiver_profiles WHERE user_id = u.id) AS has_caregiver_profile
		FROM users u
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`
	err := pool.QueryRow(ctx, query, userID).Scan(
		&resp.ID,
		&resp.Email,
		&resp.FullName,
		&resp.Gender,
		&resp.PhoneNumber,
		&resp.PhotoURL,
		&resp.HasMotherProfile,
		&resp.HasCaregiverProfile,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &resp, err
}

// UpdateUser performs a partial update on the users table.
// Only non-nil fields in input are updated.
// Used by endpoint 23 (PATCH /users/{user_id}).
func UpdateUser(pool *pgxpool.Pool, userID string, input models.UpdateUserInput) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	setParts := []string{}
	args := []any{}
	idx := 1

	if input.FullName != nil {
		setParts = append(setParts, fmt.Sprintf("full_name = $%d", idx))
		args = append(args, *input.FullName)
		idx++
	}
	if input.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", idx))
		args = append(args, *input.Email)
		idx++
	}
	if input.PhoneNumber != nil {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", idx))
		args = append(args, *input.PhoneNumber)
		idx++
	}
	if input.Gender != nil {
		setParts = append(setParts, fmt.Sprintf("gender = $%d::gender_type", idx))
		args = append(args, *input.Gender)
		idx++
	}

	if len(setParts) == 0 {
		return nil // nothing to update
	}

	args = append(args, userID)

	var query string = fmt.Sprintf(`
		UPDATE users 
		SET %s 
		WHERE id = $%d AND deleted_at IS NULL`,
		strings.Join(setParts, ", "),
		idx,
	)

	cmdTag, err := pool.Exec(ctx, query, args...)

	if err != nil {
		if isUniqueViolation(err) {
			return ErrConflict
		}
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SoftDeleteUser sets deleted_at = now() for the given user.
// Used by endpoint 25 (DELETE /users/{user_id}).
func SoftDeleteUser(pool *pgxpool.Pool, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE users 
		SET deleted_at = now() 
		WHERE id = $1 AND deleted_at IS NULL`

	cmdTag, err := pool.Exec(ctx, query, userID)

	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateUserOnboarding sets up the correct profile for a user based on their chosen role.
func UpdateUserOnboarding(pool *pgxpool.Pool, auth0ID, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tableName string
	switch role {
	case "mother":
		tableName = "mother_profiles"
	case "caregiver":
		tableName = "caregiver_profiles"
	case "doctor":
		tableName = "doctor_profiles"
	default:
		return errors.New("invalid role: must be mother, caregiver, or doctor")
	}

	var query string = fmt.Sprintf(`
		INSERT INTO %s (user_id) 
		SELECT id FROM users WHERE auth0_id = $1
		ON CONFLICT (user_id) DO NOTHING`, tableName)

	_, err := pool.Exec(ctx, query, auth0ID)
	if err != nil {
		return err
	}
	return nil
}

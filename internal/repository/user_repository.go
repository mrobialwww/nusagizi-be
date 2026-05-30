package repository

import (
	"context"
	"errors"
	"nusagizi_be/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Cari user berdasarkan Auth0 sub
func GetUserBySub(pool *pgxpool.Pool, sub string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT id, auth0_id, email, created_at, updated_at
        FROM users
        WHERE auth0_id = $1
    `
	var user models.User

	err := pool.QueryRow(ctx, query, sub).Scan(
		&user.ID,
		&user.Sub,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Insert user dari Auth0 (tanpa password)
func CreateUserFromAuth0(pool *pgxpool.Pool, sub, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO users (auth0_id, email)
        VALUES ($1, $2)
        ON CONFLICT (auth0_id) DO NOTHING
        RETURNING id, auth0_id, email, created_at, updated_at
    `
	var user models.User

	err := pool.QueryRow(ctx, query, sub, email).Scan(
		&user.ID,
		&user.Sub,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, email, created_at, updated_at
	`

	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User

	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(pool *pgxpool.Pool, id string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User

	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUserOnboarding(pool *pgxpool.Pool, auth0ID, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET role = $1, updated_at = NOW()
		WHERE auth0_id = $2
	`

	cmdTag, err := pool.Exec(ctx, query, role, auth0ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func RollbackUserOnboarding(pool *pgxpool.Pool, auth0ID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET role = NULL, updated_at = NOW()
		WHERE auth0_id = $1
	`

	_, err := pool.Exec(ctx, query, auth0ID)
	return err
}

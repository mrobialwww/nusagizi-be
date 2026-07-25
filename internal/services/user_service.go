package services

import (
	"errors"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrForbidden is returned when the requester does not own the resource.
var ErrForbidden = errors.New("forbidden")

// GetUser returns a user's profile (endpoint 45: GET /users/{user_id}).
// Returns ErrForbidden if the requester is not the target user.
func GetUser(pool *pgxpool.Pool, requesterUserID, targetUserID string) (*models.UserProfileResponse, error) {
	if requesterUserID != targetUserID {
		return nil, ErrForbidden
	}
	return repository.GetUserProfileByID(pool, targetUserID)
}

// UpdateUser partially updates a user's profile (endpoint 23: PATCH /users/{user_id}).
// Returns ErrForbidden if the requester is not the target user.
func UpdateUser(pool *pgxpool.Pool, requesterUserID, targetUserID string, input models.UpdateUserInput) error {
	if requesterUserID != targetUserID {
		return ErrForbidden
	}
	return repository.UpdateUser(pool, targetUserID, input)
}

// SoftDeleteUser marks a user as deleted (endpoint 25: DELETE /users/{user_id}).
// Returns ErrForbidden if the requester is not the target user.
func SoftDeleteUser(pool *pgxpool.Pool, requesterUserID, targetUserID string) error {
	if requesterUserID != targetUserID {
		return ErrForbidden
	}
	return repository.SoftDeleteUser(pool, targetUserID)
}

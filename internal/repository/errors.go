package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// ErrNotFound is returned when a resource is not found in the database.
var ErrNotFound = errors.New("not found")

// ErrConflict is returned when a unique constraint is violated.
var ErrConflict = errors.New("conflict")

// ErrForbidden is returned when a user does not have permission.
var ErrForbidden = errors.New("forbidden")

// ErrInvalidRole is returned when a user tries to assign an invalid role.
var ErrInvalidRole = errors.New("invalid role: must be mother, caregiver, or doctor")

const uniqueViolationCode = "23505"

// isUniqueViolation returns true if the error is a PostgreSQL unique constraint violation.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}

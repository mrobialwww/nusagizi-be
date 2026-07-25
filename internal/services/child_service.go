package services

import (
	"context"
	"errors"
	"fmt"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateChild creates a new child profile within a DB transaction (endpoint 22: POST /children).
// userID is the DB UUID (users.id) of the authenticated mother.
// Returns the new child's UUID on success.
func CreateChild(pool *pgxpool.Pool, userID string, input *models.CreateChildInput) (uuid.UUID, error) {
	// 1. Repo: resolve mother_profile_id from user_id
	motherProfileID, err := repository.GetMotherProfileByUserID(pool, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return uuid.Nil, fmt.Errorf("%w: no mother profile associated with this account", ErrForbidden)
	}
	if err != nil {
		return uuid.Nil, err
	}

	// 2. Begin transaction
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx) // no-op if already committed

	// 3. Repo: insert child main row
	childID, err := repository.CreateChild(tx, motherProfileID, input)
	if err != nil {
		return uuid.Nil, err
	}

	// 4. Repo: insert sub-profiles
	if err := repository.BulkInsertAllergyProfiles(tx, childID, input.Allergies); err != nil {
		return uuid.Nil, err
	}
	if err := repository.BulkInsertFavoriteFoods(tx, childID, input.FavoriteFoods); err != nil {
		return uuid.Nil, err
	}
	if err := repository.BulkInsertFavoriteTextures(tx, childID, input.FavoriteTextures); err != nil {
		return uuid.Nil, err
	}
	if err := repository.BulkInsertDiets(tx, childID, input.Diets); err != nil {
		return uuid.Nil, err
	}
	if err := repository.BulkInsertChronicDiseases(tx, childID, input.ChronicDiseases); err != nil {
		return uuid.Nil, err
	}

	// 5. Commit
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return childID, nil
}

// UpdateChild updates a child profile (endpoint 24: PATCH /children/{child_id}).
// Uses partial update for the main row and replace-all for sub-profiles that are present in input.
func UpdateChild(pool *pgxpool.Pool, userID string, childID uuid.UUID, input *models.UpdateChildInput) error {
	// 1. Repo: resolve mother_profile_id
	motherProfileID, err := repository.GetMotherProfileByUserID(pool, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("%w: no mother profile associated with this account", ErrForbidden)
	}
	if err != nil {
		return err
	}

	// 2. Repo: validate ownership
	owns, err := repository.CheckChildOwnership(pool, childID, motherProfileID)
	if err != nil {
		return err
	}
	if !owns {
		return ErrForbidden
	}

	// 3. Begin transaction
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 4. Repo: partial update of main child row
	if err := repository.UpdateChild(tx, childID, input); err != nil {
		return err
	}

	// 5. Repo: replace-all for each sub-profile type that was sent in the request
	// A nil pointer means the field was not sent → keep existing data.
	if input.Allergies != nil {
		if err := repository.DeleteAllergyProfiles(tx, childID); err != nil {
			return err
		}
		if err := repository.BulkInsertAllergyProfiles(tx, childID, input.Allergies); err != nil {
			return err
		}
	}
	if input.FavoriteFoods != nil {
		if err := repository.DeleteFavoriteFoods(tx, childID); err != nil {
			return err
		}
		if err := repository.BulkInsertFavoriteFoods(tx, childID, *input.FavoriteFoods); err != nil {
			return err
		}
	}
	if input.FavoriteTextures != nil {
		if err := repository.DeleteFavoriteTextures(tx, childID); err != nil {
			return err
		}
		if err := repository.BulkInsertFavoriteTextures(tx, childID, *input.FavoriteTextures); err != nil {
			return err
		}
	}
	if input.Diets != nil {
		if err := repository.DeleteDiets(tx, childID); err != nil {
			return err
		}
		if err := repository.BulkInsertDiets(tx, childID, *input.Diets); err != nil {
			return err
		}
	}
	if input.ChronicDiseases != nil {
		if err := repository.DeleteChronicDiseases(tx, childID); err != nil {
			return err
		}
		if err := repository.BulkInsertChronicDiseases(tx, childID, *input.ChronicDiseases); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// DeleteChild soft-deletes a child (endpoint 46: DELETE /children/{child_id}).
func DeleteChild(pool *pgxpool.Pool, userID string, childID uuid.UUID) error {
	motherProfileID, err := repository.GetMotherProfileByUserID(pool, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("%w: no mother profile associated with this account", ErrForbidden)
	}
	if err != nil {
		return err
	}

	owns, err := repository.CheckChildOwnership(pool, childID, motherProfileID)
	if err != nil {
		return err
	}
	if !owns {
		return ErrForbidden
	}

	return repository.SoftDeleteChild(pool, childID)
}

// GetChildDetail returns the full child detail (endpoint 50: GET /children/{child_id}).
func GetChildDetail(pool *pgxpool.Pool, userID string, childID uuid.UUID) (*models.ChildDetailResponse, error) {
	motherProfileID, err := repository.GetMotherProfileByUserID(pool, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("%w: no mother profile associated with this account", ErrForbidden)
	}
	if err != nil {
		return nil, err
	}

	owns, err := repository.CheckChildOwnership(pool, childID, motherProfileID)
	if err != nil {
		return nil, err
	}
	if !owns {
		return nil, ErrForbidden
	}

	return repository.GetChildByID(pool, childID)
}

// GetChildrenByMother returns the lightweight children list (endpoint 51: GET /mother-profiles/{id}/children).
// Also validates that the requested mother_profile_id belongs to the calling user.
func GetChildrenByMother(pool *pgxpool.Pool, userID string, motherProfileID uuid.UUID) ([]models.ChildListItem, error) {
	actualMotherProfileID, err := repository.GetMotherProfileByUserID(pool, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("%w: no mother profile associated with this account", ErrForbidden)
	}
	if err != nil {
		return nil, err
	}

	// Verify the path param matches the user's own profile
	if actualMotherProfileID != motherProfileID {
		return nil, ErrForbidden
	}

	return repository.GetChildrenByMotherProfileID(pool, motherProfileID)
}

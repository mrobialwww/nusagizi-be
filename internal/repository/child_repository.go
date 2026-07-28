package repository

import (
	"context"
	"errors"
	"fmt"
	"nusagizi_be/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChildRepository struct {
	pool *pgxpool.Pool
}

func NewChildRepository(pool *pgxpool.Pool) *ChildRepository {
	return &ChildRepository{pool: pool}
}

// Create inserts a new child and all sub-profiles within an internal transaction.
func (r *ChildRepository) Create(ctx context.Context, motherProfileID uuid.UUID, input *models.CreateChildInput) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx) // no-op if already committed

	childID, err := r.insertChild(ctx, tx, motherProfileID, input)
	if err != nil {
		return uuid.Nil, err
	}

	// Insert allergy sub-profiles
	if input.Allergies != nil {
		var query string = `
			INSERT INTO child_allergy_profile (child_id, category, allergen_name)
			VALUES ($1, $2::child_allergy_category, $3)`

		for _, name := range input.Allergies.Food {
			if _, err := tx.Exec(ctx, query, childID, "food", name); err != nil {
				return uuid.Nil, err
			}
		}
		for _, name := range input.Allergies.Medicine {
			if _, err := tx.Exec(ctx, query, childID, "medicine", name); err != nil {
				return uuid.Nil, err
			}
		}
		for _, name := range input.Allergies.Animal {
			if _, err := tx.Exec(ctx, query, childID, "animal", name); err != nil {
				return uuid.Nil, err
			}
		}
		for _, name := range input.Allergies.Others {
			if _, err := tx.Exec(ctx, query, childID, "others", name); err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert favorite food sub-profiles
	if len(input.FavoriteFoods) > 0 {
		var query string = `INSERT INTO favorite_food_profile (child_id, food_name) VALUES ($1, $2)`

		for _, food := range input.FavoriteFoods {
			if _, err := tx.Exec(ctx, query, childID, food); err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert favorite texture sub-profiles
	if len(input.FavoriteTextures) > 0 {
		var query string = `INSERT INTO favorite_texture_profile (child_id, texture_name) VALUES ($1, $2)`

		for _, texture := range input.FavoriteTextures {
			if _, err := tx.Exec(ctx, query, childID, texture); err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert diet sub-profiles
	if len(input.Diets) > 0 {
		var query string = `INSERT INTO child_diet_profile (child_id, diet_name) VALUES ($1, $2)`

		for _, diet := range input.Diets {
			if _, err := tx.Exec(ctx, query, childID, diet); err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert chronic disease sub-profiles
	if len(input.ChronicDiseases) > 0 {
		var query string = `INSERT INTO child_chronic_disease_profile (child_id, disease_name) VALUES ($1, $2)`

		for _, disease := range input.ChronicDiseases {
			if _, err := tx.Exec(ctx, query, childID, disease); err != nil {
				return uuid.Nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return childID, nil
}

// insertChild inserts the main child row and returns the new UUID. Must be called within a transaction.
func (r *ChildRepository) insertChild(ctx context.Context, tx pgx.Tx, motherProfileID uuid.UUID, input *models.CreateChildInput) (uuid.UUID, error) {
	birthDate, err := time.Parse(models.DateLayout, input.BirthDate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid birth_date format, expected DD-MM-YYYY: %w", err)
	}

	var childID uuid.UUID
	var query string = `
		INSERT INTO child (
			mother_profile_id,
			full_name, 
			gender, 
			birth_date, 
			photo_url,
			food_frequency_profile,
			food_goal_profile,
			notes_profile,
			upload_streak_days
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0)
		RETURNING id
	`
	err = tx.QueryRow(ctx, query,
		motherProfileID,
		input.FullName,
		input.Gender,
		birthDate,
		input.PhotoURL,
		input.FoodFrequency,
		input.FoodGoal,
		input.Notes,
	).Scan(&childID)

	return childID, err
}

// Update performs a partial update of the main child row and replace-all for sub-profiles, within an internal transaction.
// A nil pointer field means the field was not sent → keep existing data.
func (r *ChildRepository) Update(ctx context.Context, childID uuid.UUID, input *models.UpdateChildInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := r.updateChildRow(ctx, tx, childID, input); err != nil {
		return err
	}

	if input.Allergies != nil {
		var delQuery string = `
			DELETE FROM child_allergy_profile 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		var query string = `
			INSERT INTO child_allergy_profile (child_id, category, allergen_name)
			VALUES ($1, $2::child_allergy_category, $3)`

		for _, name := range input.Allergies.Food {
			if _, err := tx.Exec(ctx, query, childID, "food", name); err != nil {
				return err
			}
		}
		for _, name := range input.Allergies.Medicine {
			if _, err := tx.Exec(ctx, query, childID, "medicine", name); err != nil {
				return err
			}
		}
		for _, name := range input.Allergies.Animal {
			if _, err := tx.Exec(ctx, query, childID, "animal", name); err != nil {
				return err
			}
		}
		for _, name := range input.Allergies.Others {
			if _, err := tx.Exec(ctx, query, childID, "others", name); err != nil {
				return err
			}
		}
	}

	if input.FavoriteFoods != nil {
		var delQuery string = `
			DELETE FROM favorite_food_profile 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.FavoriteFoods) > 0 {
			var insQuery string = `INSERT INTO favorite_food_profile (child_id, food_name) VALUES ($1, $2)`

			for _, food := range *input.FavoriteFoods {
				if _, err := tx.Exec(ctx, insQuery, childID, food); err != nil {
					return err
				}
			}
		}
	}

	if input.FavoriteTextures != nil {
		var delQuery string = `
			DELETE FROM favorite_texture_profile 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.FavoriteTextures) > 0 {
			var insQuery string = `INSERT INTO favorite_texture_profile (child_id, texture_name) VALUES ($1, $2)`

			for _, texture := range *input.FavoriteTextures {
				if _, err := tx.Exec(ctx, insQuery, childID, texture); err != nil {
					return err
				}
			}
		}
	}

	if input.Diets != nil {
		var delQuery string = `
			DELETE FROM child_diet_profile 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.Diets) > 0 {
			var insQuery string = `INSERT INTO child_diet_profile (child_id, diet_name) VALUES ($1, $2)`

			for _, diet := range *input.Diets {
				if _, err := tx.Exec(ctx, insQuery, childID, diet); err != nil {
					return err
				}
			}
		}
	}

	if input.ChronicDiseases != nil {
		var delQuery string = `
			DELETE FROM child_chronic_disease_profile 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.ChronicDiseases) > 0 {
			var insQuery string = `INSERT INTO child_chronic_disease_profile (child_id, disease_name) VALUES ($1, $2)`
			for _, disease := range *input.ChronicDiseases {
				if _, err := tx.Exec(ctx, insQuery, childID, disease); err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit(ctx)
}

// updateChildRow performs a partial update on the main child row. Must be called within a transaction.
func (r *ChildRepository) updateChildRow(ctx context.Context, tx pgx.Tx, childID uuid.UUID, input *models.UpdateChildInput) error {
	setParts := []string{}
	args := []any{}
	idx := 1

	if input.FullName != nil {
		setParts = append(setParts, fmt.Sprintf("full_name = $%d", idx))
		args = append(args, *input.FullName)
		idx++
	}
	if input.Gender != nil {
		setParts = append(setParts, fmt.Sprintf("gender = $%d", idx))
		args = append(args, *input.Gender)
		idx++
	}
	if input.BirthDate != nil {
		birthDate, err := time.Parse(models.DateLayout, *input.BirthDate)
		if err != nil {
			return fmt.Errorf("invalid birth_date format, expected DD-MM-YYYY: %w", err)
		}
		setParts = append(setParts, fmt.Sprintf("birth_date = $%d", idx))
		args = append(args, birthDate)
		idx++
	}
	if input.PhotoURL != nil {
		setParts = append(setParts, fmt.Sprintf("photo_url = $%d", idx))
		args = append(args, *input.PhotoURL)
		idx++
	}
	if input.FoodFrequency != nil {
		setParts = append(setParts, fmt.Sprintf("food_frequency_profile = $%d", idx))
		args = append(args, *input.FoodFrequency)
		idx++
	}
	if input.FoodGoal != nil {
		setParts = append(setParts, fmt.Sprintf("food_goal_profile = $%d", idx))
		args = append(args, *input.FoodGoal)
		idx++
	}
	if input.Notes != nil {
		setParts = append(setParts, fmt.Sprintf("notes_profile = $%d", idx))
		args = append(args, *input.Notes)
		idx++
	}

	if len(setParts) == 0 {
		return nil // nothing to update in the main row
	}

	args = append(args, childID)

	var query string = fmt.Sprintf(`
		UPDATE child 
		SET %s 
		WHERE id = $%d 
			AND deleted_at IS NULL`,
		strings.Join(setParts, ", "),
		idx,
	)

	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// CheckOwnership returns true if the child belongs to the given mother profile and is not deleted.
func (r *ChildRepository) CheckOwnership(ctx context.Context, childID, motherProfileID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists bool
	var query string = `
		SELECT EXISTS(
			SELECT 1 FROM child
			WHERE id = $1 
				AND mother_profile_id = $2 
				AND deleted_at IS NULL
		)
	`
	err := r.pool.QueryRow(ctx, query, childID, motherProfileID).Scan(&exists)
	return exists, err
}

// GetSimpleByID returns simple child detail.
func (r *ChildRepository) GetSimpleByID(ctx context.Context, childID uuid.UUID) (*models.ChildSimpleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var child models.ChildSimpleResponse
	var birthDate time.Time

	query := `
		SELECT 
			id, full_name, birth_date, gender, photo_url, upload_streak_days
		FROM child
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := r.pool.QueryRow(ctx, query, childID).Scan(
		&child.ID,
		&child.FullName,
		&birthDate,
		&child.Gender,
		&child.PhotoURL,
		&child.UploadStreakDays,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	child.BirthDate = birthDate.Format(models.DateLayout)
	return &child, nil
}

// GetByID returns full child detail with all sub-profiles assembled.
func (r *ChildRepository) GetByID(ctx context.Context, childID uuid.UUID) (*models.ChildDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Fetch main child row
	var c models.Child
	var query string = `
		SELECT 
			id,
			full_name,
			gender,
			birth_date,
			photo_url,
			food_frequency_profile,
			food_goal_profile,
			notes_profile,
			upload_streak_days
		FROM child
		WHERE id = $1 
			AND deleted_at IS NULL
	`
	err := r.pool.QueryRow(ctx, query, childID).Scan(
		&c.ID,
		&c.FullName,
		&c.Gender,
		&c.BirthDate,
		&c.PhotoURL,
		&c.FoodFrequencyProfile,
		&c.FoodGoalProfile,
		&c.NotesProfile,
		&c.UploadStreakDays,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	resp := &models.ChildDetailResponse{
		ID:               c.ID.String(),
		FullName:         c.FullName,
		BirthDate:        c.BirthDate.Format(models.DateLayout),
		Gender:           c.Gender,
		PhotoURL:         c.PhotoURL,
		UploadStreakDays: c.UploadStreakDays,
		FoodFrequency:    c.FoodFrequencyProfile,
		FoodGoal:         c.FoodGoalProfile,
		Notes:            c.NotesProfile,
		// Initialise slices to empty arrays (not null) for consistent JSON output
		Allergies:        models.Allergies{Food: []string{}, Medicine: []string{}, Animal: []string{}, Others: []string{}},
		ChronicDiseases:  []string{},
		Diets:            []string{},
		FavoriteFoods:    []string{},
		FavoriteTextures: []string{},
	}

	// 2. Fetch allergies
	allergies, err := r.getChildAllergies(ctx, childID)
	if err != nil {
		return nil, err
	}
	resp.Allergies = allergies

	// 3-6. Fetch the remaining string-list sub-profiles
	if err := r.populateChildStringProfiles(ctx, childID, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// getChildAllergies is a helper to fetch and group allergies by category for a child.
func (r *ChildRepository) getChildAllergies(ctx context.Context, childID uuid.UUID) (models.Allergies, error) {
	allergies := models.Allergies{
		Food:     []string{},
		Medicine: []string{},
		Animal:   []string{},
		Others:   []string{},
	}

	var allergyQuery string = `SELECT category, allergen_name FROM child_allergy_profile WHERE child_id = $1`
	rows, err := r.pool.Query(ctx, allergyQuery, childID)
	if err != nil {
		return allergies, err
	}
	defer rows.Close()

	for rows.Next() {
		var category, name string
		if err := rows.Scan(&category, &name); err != nil {
			return allergies, err
		}
		switch category {
		case "food":
			allergies.Food = append(allergies.Food, name)
		case "medicine":
			allergies.Medicine = append(allergies.Medicine, name)
		case "animal":
			allergies.Animal = append(allergies.Animal, name)
		case "others":
			allergies.Others = append(allergies.Others, name)
		}
	}
	return allergies, rows.Err()
}

// populateChildStringProfiles fetches and populates favorite foods, textures, diets, and diseases for a child.
func (r *ChildRepository) populateChildStringProfiles(ctx context.Context, childID uuid.UUID, resp *models.ChildDetailResponse) error {
	queries := []struct {
		query string
		dest  *[]string
	}{
		{
			`
			SELECT food_name 
			FROM favorite_food_profile 
			WHERE child_id = $1
			`, &resp.FavoriteFoods,
		},
		{
			`
			SELECT texture_name 
			FROM favorite_texture_profile 
			WHERE child_id = $1
			`, &resp.FavoriteTextures,
		},
		{
			`
			SELECT diet_name 
			FROM child_diet_profile 
			WHERE child_id = $1
			`, &resp.Diets,
		},
		{
			`
			SELECT disease_name 
			FROM child_chronic_disease_profile 
			WHERE child_id = $1
			`, &resp.ChronicDiseases,
		},
	}

	for _, q := range queries {
		rows, err := r.pool.Query(ctx, q.query, childID)
		if err != nil {
			return err
		}

		var result []string
		for rows.Next() {
			var s string
			if err := rows.Scan(&s); err != nil {
				rows.Close()
				return err
			}
			result = append(result, s)
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			return err
		}
		*q.dest = result
	}
	return nil
}

// GetByMotherProfileID returns a lightweight list of children for a mother.
func (r *ChildRepository) GetByMotherProfileID(ctx context.Context, motherProfileID uuid.UUID) ([]models.ChildListItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, full_name, birth_date, gender, photo_url
		FROM child
		WHERE mother_profile_id = $1 
			AND deleted_at IS NULL
		ORDER BY created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []models.ChildListItem{}
	for rows.Next() {
		var item models.ChildListItem
		var birthDate time.Time
		if err := rows.Scan(&item.ID, &item.FullName, &birthDate, &item.Gender, &item.PhotoURL); err != nil {
			return nil, err
		}
		item.BirthDate = birthDate.Format(models.DateLayout)
		result = append(result, item)
	}
	return result, rows.Err()
}

// SoftDelete sets deleted_at = now() for the given child.
func (r *ChildRepository) SoftDelete(ctx context.Context, childID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE child 
		SET deleted_at = now() 
		WHERE id = $1 
			AND deleted_at IS NULL`

	cmdTag, err := r.pool.Exec(ctx, query, childID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

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

// CreateChild inserts a new child row and returns the new UUID.
// Must be called within a transaction (pgx.Tx).
func CreateChild(tx pgx.Tx, motherProfileID uuid.UUID, input *models.CreateChildInput) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

// BulkInsertAllergyProfiles inserts all allergy category rows for a child.
// Must be called within a transaction.
func BulkInsertAllergyProfiles(tx pgx.Tx, childID uuid.UUID, allergies *models.Allergies) error {
	if allergies == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type entry struct {
		category string
		names    []string
	}
	entries := []entry{
		{"food", allergies.Food},
		{"medicine", allergies.Medicine},
		{"animal", allergies.Animal},
		{"others", allergies.Others},
	}
	for _, e := range entries {
		for _, name := range e.names {
			var query string = `
				INSERT INTO child_allergy_profile (
					child_id,
					category,
					allergen_name
				)
				VALUES ($1, $2::child_allergy_category, $3)`

			if _, err := tx.Exec(ctx, query, childID, e.category, name); err != nil {
				return err
			}
		}
	}
	return nil
}

// BulkInsertFavoriteFoods inserts favorite food profile rows for a child.
// Must be called within a transaction.
func BulkInsertFavoriteFoods(tx pgx.Tx, childID uuid.UUID, foods []string) error {
	if len(foods) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO favorite_food_profile (child_id, food_name)
		VALUES ($1, $2)`

	for _, food := range foods {
		if _, err := tx.Exec(ctx, query, childID, food); err != nil {
			return err
		}
	}
	return nil
}

// BulkInsertFavoriteTextures inserts favorite texture profile rows for a child.
// Must be called within a transaction.
func BulkInsertFavoriteTextures(tx pgx.Tx, childID uuid.UUID, textures []string) error {
	if len(textures) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO favorite_texture_profile (child_id, texture_name) 
		VALUES ($1, $2)`

	for _, texture := range textures {
		if _, err := tx.Exec(ctx, query, childID, texture); err != nil {
			return err
		}
	}
	return nil
}

// BulkInsertDiets inserts child diet profile rows for a child.
// Must be called within a transaction.
func BulkInsertDiets(tx pgx.Tx, childID uuid.UUID, diets []string) error {
	if len(diets) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO child_diet_profile (child_id, diet_name) 
		VALUES ($1, $2)`

	for _, diet := range diets {
		if _, err := tx.Exec(ctx, query, childID, diet); err != nil {
			return err
		}
	}
	return nil
}

// BulkInsertChronicDiseases inserts chronic disease profile rows for a child.
// Must be called within a transaction.
func BulkInsertChronicDiseases(tx pgx.Tx, childID uuid.UUID, diseases []string) error {
	if len(diseases) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO child_chronic_disease_profile (child_id, disease_name) 
		VALUES ($1, $2)`

	for _, disease := range diseases {
		if _, err := tx.Exec(ctx, query, childID, disease); err != nil {
			return err
		}
	}
	return nil
}

// DeleteAllergyProfiles deletes all allergy profiles for a child.
// Must be called within a transaction.
func DeleteAllergyProfiles(tx pgx.Tx, childID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM child_allergy_profile 
		WHERE child_id = $1`

	_, err := tx.Exec(ctx, query, childID)
	return err
}

// DeleteFavoriteFoods deletes all favorite food profiles for a child.
// Must be called within a transaction.
func DeleteFavoriteFoods(tx pgx.Tx, childID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM favorite_food_profile 
		WHERE child_id = $1`

	_, err := tx.Exec(ctx, query, childID)
	return err
}

// DeleteFavoriteTextures deletes all favorite texture profiles for a child.
// Must be called within a transaction.
func DeleteFavoriteTextures(tx pgx.Tx, childID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM favorite_texture_profile 
		WHERE child_id = $1`

	_, err := tx.Exec(ctx, query, childID)
	return err
}

// DeleteDiets deletes all diet profiles for a child.
// Must be called within a transaction.
func DeleteDiets(tx pgx.Tx, childID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM child_diet_profile 
		WHERE child_id = $1`

	_, err := tx.Exec(ctx, query, childID)
	return err
}

// DeleteChronicDiseases deletes all chronic disease profiles for a child.
// Must be called within a transaction.
func DeleteChronicDiseases(tx pgx.Tx, childID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM child_chronic_disease_profile 
		WHERE child_id = $1`

	_, err := tx.Exec(ctx, query, childID)
	return err
}

// UpdateChild updates the main child row with only non-nil fields (partial update).
// Must be called within a transaction.
func UpdateChild(tx pgx.Tx, childID uuid.UUID, input *models.UpdateChildInput) error {
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
		WHERE id = $%d AND deleted_at IS NULL`,
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

// CheckChildOwnership returns true if the child belongs to the given mother profile and is not deleted.
func CheckChildOwnership(pool *pgxpool.Pool, childID, motherProfileID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	var query string = `
		SELECT EXISTS(
			SELECT 1 FROM child
			WHERE id = $1 AND mother_profile_id = $2 AND deleted_at IS NULL
		)
	`
	err := pool.QueryRow(ctx, query, childID, motherProfileID).Scan(&exists)
	return exists, err
}

// GetChildByID returns full child detail with all sub-profiles assembled.
// Used by endpoint 50 (GET /children/{child_id}).
func GetChildByID(pool *pgxpool.Pool, childID uuid.UUID) (*models.ChildDetailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := pool.QueryRow(ctx, query, childID).Scan(
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
	allergies, err := getChildAllergies(ctx, pool, childID)
	if err != nil {
		return nil, err
	}
	resp.Allergies = allergies

	// 3-6. Fetch the remaining string-list sub-profiles
	if err := populateChildStringProfiles(ctx, pool, childID, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// getChildAllergies is a helper to fetch and group allergies by category for a child.
func getChildAllergies(ctx context.Context, pool *pgxpool.Pool, childID uuid.UUID) (models.Allergies, error) {
	allergies := models.Allergies{
		Food:     []string{},
		Medicine: []string{},
		Animal:   []string{},
		Others:   []string{},
	}

	var allergyQuery string = `SELECT category, allergen_name FROM child_allergy_profile WHERE child_id = $1`
	rows, err := pool.Query(ctx, allergyQuery, childID)
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
func populateChildStringProfiles(ctx context.Context, pool *pgxpool.Pool, childID uuid.UUID, resp *models.ChildDetailResponse) error {
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
		rows, err := pool.Query(ctx, q.query, childID)
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

// GetChildrenByMotherProfileID returns a lightweight list of children for a mother.
// Used by endpoint 51 (GET /mother-profiles/{id}/children).
func GetChildrenByMotherProfileID(pool *pgxpool.Pool, motherProfileID uuid.UUID) ([]models.ChildListItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, full_name, birth_date, gender, photo_url
		FROM child
		WHERE mother_profile_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC
	`
	rows, err := pool.Query(ctx, query, motherProfileID)

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

// SoftDeleteChild sets deleted_at = now() for the given child.
// Used by endpoint 46 (DELETE /children/{child_id}).
func SoftDeleteChild(pool *pgxpool.Pool, childID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE child 
		SET deleted_at = now() 
		WHERE id = $1 AND deleted_at IS NULL`

	cmdTag, err := pool.Exec(ctx, query, childID)

	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

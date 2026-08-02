package repository

import (
	"context"
	"errors"
	"fmt"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/medical"
	"nusagizi_be/internal/models/menu"
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

	// insert main child row
	birthDate, err := time.Parse(models.DateLayout, input.BirthDate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid birth_date format, expected DD-MM-YYYY: %w", err)
	}

	var childID uuid.UUID
	var query string = `
		INSERT INTO children (
			mother_profile_id,
			full_name, 
			gender, 
			birth_date, 
			photo_url,
			notes_profile,
			upload_streak_days
		)
		VALUES ($1, $2, $3, $4, $5, $6, 0)
		RETURNING id
	`
	if err := tx.QueryRow(ctx, query,
		motherProfileID,
		input.FullName,
		input.Gender,
		birthDate,
		input.PhotoURL,
		input.Notes,
	).Scan(&childID); err != nil {
		return uuid.Nil, err
	}
	// Insert allergy sub-profiles
	if input.Allergies != nil {
		var query string = `
			INSERT INTO child_allergy_profiles (child_id, category, allergen_name)
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
		var query string = `INSERT INTO favorite_food_profiles (child_id, food_name) VALUES ($1, $2)`

		for _, food := range input.FavoriteFoods {
			if _, err := tx.Exec(ctx, query, childID, food); err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert favorite texture sub-profiles
	if len(input.FavoriteTextures) > 0 {
		var query string = `INSERT INTO favorite_texture_profiles (child_id, texture_name) VALUES ($1, $2)`

		for _, texture := range input.FavoriteTextures {
			if _, err := tx.Exec(ctx, query, childID, texture); err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert diet sub-profiles
	if len(input.Diets) > 0 {
		var query string = `INSERT INTO child_diet_profiles (child_id, diet_name) VALUES ($1, $2)`

		for _, diet := range input.Diets {
			if _, err := tx.Exec(ctx, query, childID, diet); err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert chronic disease sub-profiles
	if len(input.ChronicDiseases) > 0 {
		var query string = `INSERT INTO child_chronic_disease_profiles (child_id, disease_name) VALUES ($1, $2)`

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
			DELETE FROM child_allergy_profiles 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		var query string = `
			INSERT INTO child_allergy_profiles (child_id, category, allergen_name)
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
			DELETE FROM favorite_food_profiles 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.FavoriteFoods) > 0 {
			var insQuery string = `INSERT INTO favorite_food_profiles (child_id, food_name) VALUES ($1, $2)`

			for _, food := range *input.FavoriteFoods {
				if _, err := tx.Exec(ctx, insQuery, childID, food); err != nil {
					return err
				}
			}
		}
	}

	if input.FavoriteTextures != nil {
		var delQuery string = `
			DELETE FROM favorite_texture_profiles 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.FavoriteTextures) > 0 {
			var insQuery string = `INSERT INTO favorite_texture_profiles (child_id, texture_name) VALUES ($1, $2)`

			for _, texture := range *input.FavoriteTextures {
				if _, err := tx.Exec(ctx, insQuery, childID, texture); err != nil {
					return err
				}
			}
		}
	}

	if input.Diets != nil {
		var delQuery string = `
			DELETE FROM child_diet_profiles 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.Diets) > 0 {
			var insQuery string = `INSERT INTO child_diet_profiles (child_id, diet_name) VALUES ($1, $2)`

			for _, diet := range *input.Diets {
				if _, err := tx.Exec(ctx, insQuery, childID, diet); err != nil {
					return err
				}
			}
		}
	}

	if input.ChronicDiseases != nil {
		var delQuery string = `
			DELETE FROM child_chronic_disease_profiles 
			WHERE child_id = $1`

		if _, err := tx.Exec(ctx, delQuery, childID); err != nil {
			return err
		}
		if len(*input.ChronicDiseases) > 0 {
			var insQuery string = `INSERT INTO child_chronic_disease_profiles (child_id, disease_name) VALUES ($1, $2)`
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
		UPDATE children 
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
			SELECT 1 FROM children
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
		FROM children
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
			notes_profile,
			upload_streak_days
		FROM children
		WHERE id = $1 
			AND deleted_at IS NULL
	`
	err := r.pool.QueryRow(ctx, query, childID).Scan(
		&c.ID,
		&c.FullName,
		&c.Gender,
		&c.BirthDate,
		&c.PhotoURL,
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

	var allergyQuery string = `
		SELECT 
			category, 
			allergen_name 
		FROM child_allergy_profiles 
		WHERE child_id = $1`

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
			FROM favorite_food_profiles 
			WHERE child_id = $1
			`, &resp.FavoriteFoods,
		},
		{
			`
			SELECT texture_name 
			FROM favorite_texture_profiles 
			WHERE child_id = $1
			`, &resp.FavoriteTextures,
		},
		{
			`
			SELECT diet_name 
			FROM child_diet_profiles 
			WHERE child_id = $1
			`, &resp.Diets,
		},
		{
			`
			SELECT disease_name 
			FROM child_chronic_disease_profiles 
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
		FROM children
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
		UPDATE children 
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

// GetChildrenDataForMenu constructs the AI payload for all children of a mother.
// It uses batched queries for allergies, growth, and medical notes
func (r *ChildRepository) GetChildrenDataForMenu(ctx context.Context, motherProfileID uuid.UUID) ([]menu.FoodEnginePayloadChild, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Get all children data from mother
	children, err := r.fetchChildren(ctx, motherProfileID)
	if err != nil {
		return nil, err
	}
	if len(children) == 0 {
		return nil, nil
	}

	// Map children ID to child struct
	childIDs := make([]uuid.UUID, len(children))
	for i, c := range children {
		childIDs[i] = c.ID
	}

	// Get child allergies data in batch
	allergiesByChild, err := r.fetchAllergies(ctx, childIDs)
	if err != nil {
		return nil, err
	}

	// Get child growth data in batch
	growthByChild, err := r.fetchGrowthHistory(ctx, childIDs)
	if err != nil {
		return nil, err
	}

	// Get child medical notes data in batch
	medicalNotesByChild, err := r.fetchActiveMedicalNotes(ctx, childIDs)
	if err != nil {
		return nil, err
	}

	// Map medical notes ID to nutrition targets
	noteIDs := make([]uuid.UUID, 0, len(medicalNotesByChild))
	for _, note := range medicalNotesByChild {
		noteIDs = append(noteIDs, note.ID)
	}

	// Get nutrition targets data in batch
	targetsByNote, err := r.fetchNutritionTargets(ctx, noteIDs)
	if err != nil {
		return nil, err
	}

	// Construct payload for each child
	payload := make([]menu.FoodEnginePayloadChild, 0, len(children))
	for _, c := range children {
		// Calculate age in months
		usiaBulan := calculateAgeInMonths(c.BirthDate, time.Now())

		// Create child payload
		childPayload := menu.FoodEnginePayloadChild{
			ID:                   c.ID.String(),
			Nama:                 c.Nama,
			Sex:                  c.Sex,
			UsiaBulan:            usiaBulan,
			PreferensiHarga:      "seimbang", // temporary hardcode
			JumlahProteinPerHari: 1,          // temporary hardcode
			Alergi:               allergiesByChild[c.ID],
			RiwayatBerat:         [][]any{},
			Kondisi:              []menu.FoodEnginePayloadKondisi{}, // temporary fill with empty array
		}

		// AI expects an empty array [] if there are no allergies
		if childPayload.Alergi == nil {
			childPayload.Alergi = []string{}
		}

		// Attach latest physical measurements and format weight history graph as [AgeInMonths, WeightKg]
		if growth, ok := growthByChild[c.ID]; ok && len(growth) > 0 {
			// Latest measurement is first in the slice (newest)
			latest := growth[0]
			childPayload.BeratKg = latest.WeightKg
			childPayload.PanjangCm = latest.HeightCm
			childPayload.LingkarKepalaCm = latest.HeadCircCm

			// Convert each measurement into a [AgeInMonths, WeightKg] point for the growth curve
			history := make([][]any, len(growth))
			for i, g := range growth {
				history[i] = []any{calculateAgeInMonths(c.BirthDate, g.MeasurementDate), g.WeightKg}
			}
			childPayload.RiwayatBerat = history
		}

		// Initialize default medical prescription (calculates max safe daily sodium based on child's age)
		resepDokter := menu.FoodEnginePayloadResepDokter{
			Tingkat:         nil, // temporary hardcode
			NomorStr:        nil, // temporary hardcode
			NamaDokter:      "-",
			CatatanTambahan: "-",
			KomposisiPerKg: menu.FoodEnginePayloadKomposisi{
				EnergiKkal:    0,
				ProteinG:      0,
				NatriumMgMaks: getNatriumMaks(usiaBulan),
			},
		}
		// Override default prescription if child has an active doctor's note and targeted nutrition
		if note, ok := medicalNotesByChild[c.ID]; ok {
			resepDokter.NamaDokter = note.DoctorName
			resepDokter.CatatanTambahan = note.Recommendation
			if targets, ok := targetsByNote[note.ID]; ok {
				resepDokter.KomposisiPerKg.EnergiKkal = targets.EnergiKkal
				resepDokter.KomposisiPerKg.ProteinG = targets.ProteinG
			}
		}
		childPayload.ResepDokter = resepDokter

		payload = append(payload, childPayload)
	}

	return payload, nil
}

// fetchChildren fetches all active children for a mother.
func (r *ChildRepository) fetchChildren(ctx context.Context, motherProfileID uuid.UUID) ([]models.ChildInfo, error) {
	const query = `
		SELECT 
			id,
			full_name, 
			gender, 
			birth_date
		FROM children 
		WHERE mother_profile_id = $1 
			AND deleted_at IS NULL`
	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, fmt.Errorf("failed to query children: %w", err)
	}
	defer rows.Close()

	var children []models.ChildInfo
	for rows.Next() {
		var c models.ChildInfo
		if err := rows.Scan(&c.ID, &c.Nama, &c.Sex, &c.BirthDate); err != nil {
			return nil, fmt.Errorf("failed to scan child row: %w", err)
		}
		children = append(children, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating children rows: %w", err)
	}
	return children, nil
}

// fetchAllergies batch-loads allergen names for all given children in a single query, grouped by child ID.
func (r *ChildRepository) fetchAllergies(ctx context.Context, childIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	const query = `
		SELECT child_id, allergen_name 
		FROM child_allergy_profiles 
		WHERE child_id = ANY($1)`

	rows, err := r.pool.Query(ctx, query, childIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query allergies: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]string)
	for rows.Next() {
		var childID uuid.UUID
		var allergen string
		if err := rows.Scan(&childID, &allergen); err != nil {
			return nil, fmt.Errorf("failed to scan allergy row: %w", err)
		}
		result[childID] = append(result[childID], allergen)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating allergy rows: %w", err)
	}
	return result, nil
}

// fetchGrowthHistory batch-loads up to the 4 most recent growth reports per child (index 0 = most recent)
// using a window function instead of one "ORDER BY ... LIMIT" query per child.
func (r *ChildRepository) fetchGrowthHistory(ctx context.Context, childIDs []uuid.UUID) (map[uuid.UUID][]models.GrowthPoint, error) {
	const query = `
		SELECT child_id, measurement_date, weight_kg, height_cm, head_circumference_cm
		FROM (
			SELECT
				child_id, measurement_date, weight_kg, height_cm, head_circumference_cm,
				ROW_NUMBER() OVER (PARTITION BY child_id ORDER BY measurement_date DESC) AS rn
			FROM child_growth_reports
			WHERE child_id = ANY($1)
		) ranked
		WHERE rn <= 4
		ORDER BY child_id, rn`
	rows, err := r.pool.Query(ctx, query, childIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query growth history: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]models.GrowthPoint)
	for rows.Next() {
		var childID uuid.UUID
		var g models.GrowthPoint
		if err := rows.Scan(&childID, &g.MeasurementDate, &g.WeightKg, &g.HeightCm, &g.HeadCircCm); err != nil {
			return nil, fmt.Errorf("failed to scan growth row: %w", err)
		}
		result[childID] = append(result[childID], g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating growth rows: %w", err)
	}
	return result, nil
}

// fetchActiveMedicalNotes batch-loads the single most recent active medical note
// per child using a window function, preventing N+1 queries.
func (r *ChildRepository) fetchActiveMedicalNotes(ctx context.Context, childIDs []uuid.UUID) (map[uuid.UUID]models.MedicalNote, error) {
	const query = `
		SELECT child_id, id, doctor_name, recommendation
		FROM (
			SELECT
				child_id, id, doctor_name, recommendation,
				ROW_NUMBER() OVER (P
					ARTITION BY child_id ORDER BY created_at DESC
				) AS rn
			FROM medical_notes
			WHERE child_id = ANY($1) AND valid_date >= CURRENT_DATE
		) ranked
		WHERE rn = 1`
	rows, err := r.pool.Query(ctx, query, childIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query medical notes: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]models.MedicalNote)
	for rows.Next() {
		var childID uuid.UUID
		var note models.MedicalNote
		if err := rows.Scan(&childID, &note.ID, &note.DoctorName, &note.Recommendation); err != nil {
			return nil, fmt.Errorf("failed to scan medical note row: %w", err)
		}
		result[childID] = note
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating medical note rows: %w", err)
	}
	return result, nil
}

// fetchNutritionTargets batch-loads calorie/protein targets for the given medical notes in one query instead of one query per note.
func (r *ChildRepository) fetchNutritionTargets(ctx context.Context, noteIDs []uuid.UUID) (map[uuid.UUID]models.NutritionTargets, error) {
	result := make(map[uuid.UUID]models.NutritionTargets)
	if len(noteIDs) == 0 {
		return result, nil
	}

	const query = `
		SELECT medical_note_id, nutrient_name, quantity 
		FROM daily_nutrition_targets 
		WHERE medical_note_id = ANY($1)`

	rows, err := r.pool.Query(ctx, query, noteIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query nutrition targets: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var noteID uuid.UUID
		var nutrientName string
		var quantity float64
		if err := rows.Scan(&noteID, &nutrientName, &quantity); err != nil {
			return nil, fmt.Errorf("failed to scan nutrition target row: %w", err)
		}
		t := result[noteID]
		switch nutrientName {
		case string(medical.NutrientCalories):
			t.EnergiKkal = quantity
		case string(medical.NutrientProtein):
			t.ProteinG = quantity
		}
		result[noteID] = t
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating nutrition target rows: %w", err)
	}
	return result, nil
}

// ================================== HELPER FUNCTION ===================================
func calculateAgeInMonths(birthDate, targetDate time.Time) int {
	months := (targetDate.Year()-birthDate.Year())*12 + int(targetDate.Month()-birthDate.Month())
	if targetDate.Day() < birthDate.Day() {
		months--
	}
	if months < 0 {
		return 0
	}
	return months
}

func getNatriumMaks(usiaBulan int) float64 {
	if usiaBulan <= 5 {
		return 120
	} else if usiaBulan <= 11 {
		return 370
	} else if usiaBulan <= 36 {
		return 800
	}
	return 900
}

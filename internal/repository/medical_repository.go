package repository

import (
	"context"
	"fmt"
	"time"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/medical"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MedicalRepository struct {
	pool *pgxpool.Pool
}

func NewMedicalRepository(pool *pgxpool.Pool) *MedicalRepository {
	return &MedicalRepository{pool: pool}
}

// GetMedicalNotes gets a list of medical notes (active or history) for a mother's children.
func (r *MedicalRepository) GetMedicalNotes(ctx context.Context, motherProfileID uuid.UUID, status string, month *int) ([]medical.MedicalNoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT 
			m.id, 
			c.full_name AS child_name, 
			m.recommendation, 
			m.valid_until, 
			m.created_at,
			(SELECT COUNT(*) FROM medical_restrictions WHERE medical_note_id = m.id 
				AND restriction_type = 'prohibition') AS prohibition_count,
			(SELECT COUNT(*) FROM medical_restrictions WHERE medical_note_id = m.id 
				AND restriction_type = 'allergy') AS allergy_count
		FROM medical_notes m
		JOIN children c ON m.child_id = c.id
		WHERE c.mother_profile_id = $1
	`

	switch status {
	case "active":
		query += ` AND m.valid_until >= CURRENT_DATE`
	case "history":
		query += ` AND m.valid_until < CURRENT_DATE`
	}

	if month != nil {
		query += fmt.Sprintf(` AND EXTRACT(MONTH FROM m.valid_until) = %d`, *month)
	}

	query += ` ORDER BY m.valid_until DESC`

	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []medical.MedicalNoteResponse
	for rows.Next() {
		var res medical.MedicalNoteResponse
		var validDate time.Time
		if err := rows.Scan(
			&res.ID, &res.ChildName, &res.Recommendation, &validDate, &res.CreatedAt, &res.ProhibitionCount, &res.AllergyCount,
		); err != nil {
			return nil, err
		}
		res.ValidUntil = validDate.Format(models.DateLayout)
		results = append(results, res)
	}
	return results, rows.Err()
}

// GetChildIDByNoteID fetches the child_id associated with a medical note.
func (r *MedicalRepository) GetChildIDByNoteID(ctx context.Context, noteID uuid.UUID) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var childID uuid.UUID
	query := `
		SELECT child_id 
		FROM medical_notes 
		WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, noteID).Scan(&childID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("medical note not found")
		}
		return uuid.Nil, err
	}
	return childID, nil
}

// GetMedicalNoteDetail gets full details of a medical note.
func (r *MedicalRepository) GetMedicalNoteDetail(ctx context.Context, noteID uuid.UUID) (*medical.MedicalNoteDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var detail medical.MedicalNoteDetailResponse
	var validDate time.Time

	// Get medical note details
	queryNote := `
		SELECT 
			m.id, c.full_name, m.doctor_name, m.facility_name, 
			m.recommendation, m.valid_until, m.created_at
		FROM medical_notes m
		JOIN children c ON m.child_id = c.id
		WHERE m.id = $1
	`
	err := r.pool.QueryRow(ctx, queryNote, noteID).Scan(
		&detail.ID, &detail.ChildName, &detail.DoctorName, &detail.FacilityName,
		&detail.Recommendation, &validDate, &detail.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}
	detail.ValidUntil = validDate.Format(models.DateLayout)

	// Get medical restrictions (Prohibitions and Allergies)
	queryRestrictions := `
		SELECT restriction_type, value 
		FROM medical_restrictions 
		WHERE medical_note_id = $1`

	rowsRest, err := r.pool.Query(ctx, queryRestrictions, noteID)
	if err == nil {
		defer rowsRest.Close()
		for rowsRest.Next() {
			var rType, val string
			if err := rowsRest.Scan(&rType, &val); err == nil {
				switch rType {
				case "prohibition":
					detail.Prohibitions = append(detail.Prohibitions, val)
				case "allergy":
					detail.Allergies = append(detail.Allergies, val)
				}
			}
		}
	}
	if detail.Prohibitions == nil {
		detail.Prohibitions = []string{}
	}
	if detail.Allergies == nil {
		detail.Allergies = []string{}
	}

	// Get daily nutrition targets
	queryTargets := `
		SELECT nutrient, quantity 
		FROM daily_nutrition_targets 
		WHERE medical_note_id = $1`

	rowsTgt, err := r.pool.Query(ctx, queryTargets, noteID)
	if err == nil {
		defer rowsTgt.Close()
		for rowsTgt.Next() {
			var tgt medical.DailyNutritionTarget
			if err := rowsTgt.Scan(&tgt.Nutrient, &tgt.Quantity); err == nil {
				detail.DailyNutritionTargets = append(detail.DailyNutritionTargets, tgt)
			}
		}
	}
	if detail.DailyNutritionTargets == nil {
		detail.DailyNutritionTargets = []medical.DailyNutritionTarget{}
	}

	return &detail, nil
}

// CreateMedicalNote creates a medical note with targets and restrictions.
func (r *MedicalRepository) CreateMedicalNote(ctx context.Context, input *medical.CreateMedicalNoteInput) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	// Parse valid date
	validDate, err := time.Parse(models.DateLayout, input.ValidUntil)
	if err != nil {
		return uuid.Nil, err
	}

	// Insert medical note
	newID := uuid.New()
	queryInsertNote := `
		INSERT INTO medical_notes (id, child_id, doctor_name, facility_name, recommendation, valid_until)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, queryInsertNote, newID, input.ChildID, input.DoctorName, input.FacilityName, input.Recommendation, validDate)
	if err != nil {
		return uuid.Nil, err
	}

	// Insert daily nutrition targets
	if len(input.DailyNutritionTargets) > 0 {
		queryTarget := `
			INSERT INTO daily_nutrition_targets (medical_note_id, nutrient, quantity) 
			VALUES ($1, $2, $3)
		`

		for _, tgt := range input.DailyNutritionTargets {
			_, err = tx.Exec(ctx, queryTarget, newID, tgt.Nutrient, tgt.Quantity)
			if err != nil {
				return uuid.Nil, err
			}
		}
	}

	// Insert medical restrictions (Prohibitions and Allergies)
	queryRest := `
		INSERT INTO medical_restrictions (medical_note_id, restriction_type, value) 
		VALUES ($1, $2, $3)`

	for _, p := range input.Prohibitions {
		_, err = tx.Exec(ctx, queryRest, newID, "prohibition", p)
		if err != nil {
			return uuid.Nil, err
		}
	}
	for _, a := range input.Allergies {
		_, err = tx.Exec(ctx, queryRest, newID, "allergy", a)
		if err != nil {
			return uuid.Nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return newID, nil
}

// UpdateMedicalNote uses replace-all logic for targets and restrictions.
func (r *MedicalRepository) UpdateMedicalNote(ctx context.Context, noteID uuid.UUID, input *medical.UpdateMedicalNoteInput) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Parse valid date
	var validDate *time.Time
	if input.ValidUntil != nil {
		vd, err := time.Parse(models.DateLayout, *input.ValidUntil)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
		validDate = &vd
	}

	// Update medical note
	queryUpdate := `
		UPDATE medical_notes
		SET 
			doctor_name = COALESCE($1, doctor_name),
			facility_name = COALESCE($2, facility_name),
			recommendation = COALESCE($3, recommendation),
			valid_until = COALESCE($4, valid_until)
		WHERE id = $5
	`
	res, err := tx.Exec(ctx, queryUpdate, input.DoctorName, input.FacilityName, input.Recommendation, validDate, noteID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}

	// Update daily nutrition targets
	if input.DailyNutritionTargets != nil {
		// Delete all existing daily nutrition targets
		queryDeleteTgt := `
			DELETE FROM daily_nutrition_targets 
			WHERE medical_note_id = $1`

		_, err = tx.Exec(ctx, queryDeleteTgt, noteID)
		if err != nil {
			return err
		}

		// Insert daily nutrition targets
		queryTarget := `
			INSERT INTO daily_nutrition_targets (medical_note_id, nutrient, quantity) 
			VALUES ($1, $2, $3)
		`

		for _, tgt := range *input.DailyNutritionTargets {
			_, err = tx.Exec(ctx, queryTarget, noteID, tgt.Nutrient, tgt.Quantity)
			if err != nil {
				return err
			}
		}
	}

	// For restrictions, we replace all if either prohibitions OR allergies is provided.
	// But it's better to process them if either is provided.
	if input.Prohibitions != nil {
		// Delete all existing prohibitions
		queryDeleteProh := `
			DELETE FROM medical_restrictions 
			WHERE medical_note_id = $1 
				AND restriction_type = 'prohibition'`

		_, err = tx.Exec(ctx, queryDeleteProh, noteID)
		if err != nil {
			return err
		}

		// Insert prohibitions
		queryRest := `
			INSERT INTO medical_restrictions (medical_note_id, restriction_type, value) 
			VALUES ($1, $2, $3)`

		for _, p := range *input.Prohibitions {
			_, err = tx.Exec(ctx, queryRest, noteID, "prohibition", p)
			if err != nil {
				return err
			}
		}
	}

	if input.Allergies != nil {
		// Delete all existing allergies
		queryDeleteAllergy := `
			DELETE FROM medical_restrictions 
			WHERE medical_note_id = $1 
				AND restriction_type = 'allergy'`
		_, err = tx.Exec(ctx, queryDeleteAllergy, noteID)
		if err != nil {
			return err
		}

		// Insert allergies
		queryRest := `
			INSERT INTO medical_restrictions (medical_note_id, restriction_type, value) 
			VALUES ($1, $2, $3)`

		for _, a := range *input.Allergies {
			_, err = tx.Exec(ctx, queryRest, noteID, "allergy", a)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

// DeleteMedicalNote hard deletes a note.
func (r *MedicalRepository) DeleteMedicalNote(ctx context.Context, noteID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	queryDelete := `
		DELETE FROM medical_notes 
		WHERE id = $1`

	res, err := r.pool.Exec(ctx, queryDelete, noteID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

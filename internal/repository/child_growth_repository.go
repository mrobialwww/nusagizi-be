package repository

import (
	"context"
	"errors"
	"fmt"
	"nusagizi_be/internal/models"
	child_growth "nusagizi_be/internal/models/child_growth"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChildGrowthRepository struct {
	pool *pgxpool.Pool
}

func NewChildGrowthRepository(pool *pgxpool.Pool) *ChildGrowthRepository {
	return &ChildGrowthRepository{pool: pool}
}

// GetLatestGrowthReport returns the most recent growth report for a child (ordered by measured_at DESC).
// It includes child's gender and age_months needed for WHO Z-Score calculation.
func (r *ChildGrowthRepository) GetLatestGrowthReport(ctx context.Context, childID uuid.UUID) (*child_growth.RawGrowthReportFull, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var raw child_growth.RawGrowthReportFull
	query := `
		SELECT
			r.id,
			r.measured_at,
			r.weight_kg,
			r.height_cm,
			r.head_circumference_cm,
			(EXTRACT(YEAR  FROM AGE(r.measured_at, c.birth_date)) * 12
			+ EXTRACT(MONTH FROM AGE(r.measured_at, c.birth_date)))::int AS age_months,
			c.gender
		FROM child_growth_reports r
		JOIN child c ON c.id = r.child_id
		WHERE r.child_id = $1
		ORDER BY r.measured_at DESC
		LIMIT 1
	`
	err := r.pool.QueryRow(ctx, query, childID).Scan(
		&raw.ID,
		&raw.MeasuredAt,
		&raw.WeightKg,
		&raw.HeightCm,
		&raw.HeadCircumferenceCm,
		&raw.AgeMonths,
		&raw.Gender,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &raw, err
}

// GetGrowthMeasurements returns raw measurements and child gender filtered by age bounds (in months).
func (r *ChildGrowthRepository) GetGrowthMeasurements(ctx context.Context, childID uuid.UUID, minAge, maxAge int, analysisType string) ([]child_growth.RawGrowthMeasurement, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var selectFields string
	switch analysisType {
	case "weight_for_age":
		selectFields = "r.weight_kg"
	case "height_for_age":
		selectFields = "r.height_cm"
	case "weight_for_height", "bmi_for_age":
		selectFields = "r.weight_kg, r.height_cm"
	case "head_circumference_for_age":
		selectFields = "r.head_circumference_cm"
	default:
		selectFields = "r.weight_kg, r.height_cm, r.head_circumference_cm"
	}

	query := fmt.Sprintf(`
		SELECT
			%s,
			(EXTRACT(YEAR  FROM AGE(r.measured_at, c.birth_date)) * 12 + EXTRACT(MONTH FROM AGE(r.measured_at, c.birth_date)))::int AS age_months,
			c.gender
		FROM child_growth_reports r
		JOIN child c ON c.id = r.child_id
		WHERE r.child_id = $1
			AND r.measured_at >= c.birth_date + make_interval(months => $2)
			AND r.measured_at <  c.birth_date + make_interval(months => $3 + 1)
		ORDER BY r.measured_at ASC
	`, selectFields)

	rows, err := r.pool.Query(ctx, query, childID, minAge, maxAge)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var results []child_growth.RawGrowthMeasurement
	var gender string
	for rows.Next() {
		var res child_growth.RawGrowthMeasurement

		var scanArgs []interface{}
		switch analysisType {
		case "weight_for_age":
			scanArgs = append(scanArgs, &res.WeightKg)
		case "height_for_age":
			scanArgs = append(scanArgs, &res.HeightCm)
		case "weight_for_height", "bmi_for_age":
			scanArgs = append(scanArgs, &res.WeightKg, &res.HeightCm)
		case "head_circumference_for_age":
			scanArgs = append(scanArgs, &res.HeadCircumferenceCm)
		default:
			scanArgs = append(scanArgs, &res.WeightKg, &res.HeightCm, &res.HeadCircumferenceCm)
		}
		scanArgs = append(scanArgs, &res.AgeMonths, &res.Gender)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, "", err
		}
		gender = res.Gender
		results = append(results, res)
	}

	return results, gender, rows.Err()
}

// GetGrowthReports returns all growth reports for a child, ordered by measured_at DESC.
// It includes child's gender and age_months needed for WHO Z-Score calculation.
func (r *ChildGrowthRepository) GetGrowthReports(ctx context.Context, childID uuid.UUID) ([]child_growth.RawGrowthReportFull, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			r.id,
			r.measured_at,
			r.weight_kg,
			r.height_cm,
			r.head_circumference_cm,
			(EXTRACT(YEAR  FROM AGE(r.measured_at, c.birth_date)) * 12
			+ EXTRACT(MONTH FROM AGE(r.measured_at, c.birth_date)))::int AS age_months,
			c.gender
		FROM child_growth_reports r
		JOIN child c ON c.id = r.child_id
		WHERE r.child_id = $1
		ORDER BY r.measured_at DESC
	`
	rows, err := r.pool.Query(ctx, query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []child_growth.RawGrowthReportFull
	for rows.Next() {
		var raw child_growth.RawGrowthReportFull
		if err := rows.Scan(
			&raw.ID,
			&raw.MeasuredAt,
			&raw.WeightKg,
			&raw.HeightCm,
			&raw.HeadCircumferenceCm,
			&raw.AgeMonths,
			&raw.Gender,
		); err != nil {
			return nil, err
		}
		reports = append(reports, raw)
	}
	return reports, rows.Err()
}

// CreateGrowthReport creates a new growth report.
func (r *ChildGrowthRepository) CreateGrowthReport(ctx context.Context, childID uuid.UUID, input *child_growth.CreateGrowthReportInput) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	measuredAt, err := time.Parse(models.DateLayout, input.MeasuredAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid measured_at format: %w", err)
	}

	var newID uuid.UUID
	query := `
		INSERT INTO child_growth_reports (child_id, measured_at, weight_kg, height_cm, head_circumference_cm)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = r.pool.QueryRow(ctx, query, childID, measuredAt, input.WeightKg, input.HeightCm, input.HeadCircumferenceCm).Scan(&newID)
	return newID, err
}

// UpdateGrowthReport updates an existing growth report partially.
func (r *ChildGrowthRepository) UpdateGrowthReport(ctx context.Context, reportID uuid.UUID, input *child_growth.UpdateGrowthReportInput) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// COALESCE approach for partial updates
	query := `
		UPDATE child_growth_reports
		SET 
			measured_at = COALESCE($1, measured_at),
			weight_kg = COALESCE($2, weight_kg),
			height_cm = COALESCE($3, height_cm),
			head_circumference_cm = COALESCE($4, head_circumference_cm)
		WHERE id = $5
	`

	var measuredAt *time.Time
	if input.MeasuredAt != nil {
		t, err := time.Parse(models.DateLayout, *input.MeasuredAt)
		if err != nil {
			return fmt.Errorf("invalid measured_at format: %w", err)
		}
		measuredAt = &t
	}

	res, err := r.pool.Exec(ctx, query, measuredAt, input.WeightKg, input.HeightCm, input.HeadCircumferenceCm, reportID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteGrowthReport deletes a growth report.
func (r *ChildGrowthRepository) DeleteGrowthReport(ctx context.Context, reportID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM child_growth_reports 
		WHERE id = $1`

	res, err := r.pool.Exec(ctx, query, reportID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

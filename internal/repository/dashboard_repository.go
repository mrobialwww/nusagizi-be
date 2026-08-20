package repository

import (
	"context"
	"time"

	"nusagizi_be/internal/models/dashboard"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	pool *pgxpool.Pool
}

func NewDashboardRepository(pool *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{pool: pool}
}

// GetDashboardSummary retrieves aggregated summary for all children belonging to a mother.
func (r *DashboardRepository) GetDashboardSummary(ctx context.Context, motherProfileID uuid.UUID) ([]dashboard.ChildSummaryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `
		SELECT 
			c.id, 
			c.full_name, 
			c.photo_url,
			c.birth_date, 
			c.gender,
			c.streak_days,
			c.last_streak_date,
			gr.height_cm, 
			gr.weight_kg,
			gr.head_circumference_cm,
			gr.measured_at,
			dr.kpsp_score,
			dr.kpsp_answers_count,
			nr.calories,
			nr.target_calories,
			nr.protein, 
			nr.target_protein,
			nr.fat,
			nr.target_fat,
			nr.carbohydrate,
			nr.target_carbohydrate
		FROM children c
		LEFT JOIN LATERAL (
			SELECT height_cm, weight_kg, head_circumference_cm, measured_at 
			FROM child_growth_reports 
			WHERE child_id = c.id 
			ORDER BY measured_at DESC LIMIT 1
		) gr ON true
		LEFT JOIN LATERAL (
			SELECT cdr.kpsp_score,
				(SELECT COUNT(*) FROM assessment_kpsp_answers aka WHERE aka.child_development_report_id = cdr.id) as kpsp_answers_count
			FROM child_development_reports cdr
			WHERE cdr.child_id = c.id 
			ORDER BY cdr.created_at DESC LIMIT 1
		) dr ON true
		LEFT JOIN LATERAL (
			SELECT calories, target_calories, protein, target_protein, fat, target_fat, carbohydrate, target_carbohydrate
			FROM child_nutrition_reports 
			WHERE child_id = c.id 
			ORDER BY created_at DESC LIMIT 1
		) nr ON true
		WHERE c.mother_profile_id = $1
	`
	rows, err := r.pool.Query(ctx, query, motherProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboard.ChildSummaryResponse
	for rows.Next() {
		var res dashboard.ChildSummaryResponse
		var birthDate time.Time
		if err := rows.Scan(
			&res.ID, &res.FullName, &res.PhotoUrl, &birthDate, &res.Gender, &res.Streak, &res.LastStreakDateRaw,
			&res.HeightCm, &res.WeightKg, &res.HeadCircumferenceCm, &res.GrowthMeasuredAt,
			&res.KpspScore, &res.KpspAnswersCount,
			&res.Calories, &res.TargetCalories,
			&res.Protein, &res.TargetProtein,
			&res.Fat, &res.TargetFat,
			&res.Carbohydrate, &res.TargetCarbohydrate,
		); err != nil {
			return nil, err
		}
		res.BirthDateRaw = birthDate
		results = append(results, res)
	}
	return results, rows.Err()
}

package repository

import (
	"context"
	"time"
	
	"nusagizi_be/internal/models"
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
			c.birth_date, 
			c.gender,
			c.upload_streak_days,
			gr.height_cm, 
			gr.weight_kg,
			gr.head_circumference_cm,
			dr.kpsp_score,
			nr.protein, 
			nr.target_protein
		FROM children c
		LEFT JOIN LATERAL (
			SELECT height_cm, weight_kg, head_circumference_cm 
			FROM child_growth_reports 
			WHERE child_id = c.id 
			ORDER BY measured_at DESC LIMIT 1
		) gr ON true
		LEFT JOIN LATERAL (
			SELECT kpsp_score 
			FROM child_development_reports 
			WHERE child_id = c.id 
			ORDER BY created_at DESC LIMIT 1
		) dr ON true
		LEFT JOIN LATERAL (
			SELECT protein, target_protein 
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
			&res.ID, &res.FullName, &birthDate, &res.Gender, &res.Streak,
			&res.HeightCm, &res.WeightKg, &res.HeadCircumferenceCm,
			&res.KpspScore,
			&res.Protein, &res.TargetProtein,
		); err != nil {
			return nil, err
		}
		res.BirthDate = birthDate.Format(models.DateLayout)
		results = append(results, res)
	}
	return results, rows.Err()
}

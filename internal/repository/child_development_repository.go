package repository

import (
	"context"
	"fmt"
	child_dev "nusagizi_be/internal/models/child_development"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChildDevelopmentRepository struct {
	pool *pgxpool.Pool
}

func NewChildDevelopmentRepository(pool *pgxpool.Pool) *ChildDevelopmentRepository {
	return &ChildDevelopmentRepository{pool: pool}
}

// GetLatestDevelopmentReport returns the latest KPSP report with aggregated domains.
func (r *ChildDevelopmentRepository) GetLatestDevelopmentReport(ctx context.Context, childID uuid.UUID) (*child_dev.DevelopmentReportResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Fetch latest report ID
	var report child_dev.DevelopmentReportResponse
	var reportID uuid.UUID

	queryReport := `
		SELECT id, kpsp_score, next_check_date, created_at
		FROM child_development_reports
		WHERE child_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var nextCheckDate *time.Time
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, queryReport, childID).Scan(&reportID, &report.KPSPScore, &nextCheckDate, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, err
	}
	report.ID = reportID

	if nextCheckDate != nil {
		ncd := nextCheckDate.Format("02-01-2006")
		report.NextCheckDate = &ncd
	} else {
		done := "done"
		report.NextCheckDate = &done
	}

	report.CreatedAt = &createdAt

	// Fetch domains aggregation
	queryDomains := `
		SELECT 
			q.developmental_domain,
			COUNT(q.id) as total_question,
			COUNT(a.id) FILTER (WHERE a.assessment_kpsp_answer = true) as true_answer
		FROM assessment_kpsp_answers a
		JOIN assessment_kpsp_questions q ON a.assessment_kpsp_question_id = q.id
		WHERE a.child_development_report_id = $1
		GROUP BY q.developmental_domain
	`
	rows, err := r.pool.Query(ctx, queryDomains, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d child_dev.DomainAggregate
		if err := rows.Scan(&d.DevelopmentalDomain, &d.TotalQuestion, &d.TrueAnswer); err != nil {
			return nil, err
		}
		report.Domains = append(report.Domains, d)
	}

	return &report, nil
}

// GetDevelopmentReports returns a summary list of past development reports.
func (r *ChildDevelopmentRepository) GetDevelopmentReports(ctx context.Context, childID uuid.UUID) ([]child_dev.DevelopmentReportResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, kpsp_score, month_target, created_at
		FROM child_development_reports
		WHERE child_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, childID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []child_dev.DevelopmentReportResponse
	for rows.Next() {
		var rep child_dev.DevelopmentReportResponse
		var monthTarget int
		var createdAt time.Time

		if err := rows.Scan(&rep.ID, &rep.KPSPScore, &monthTarget, &createdAt); err != nil {
			return nil, err
		}
		rep.MonthTarget = &monthTarget
		rep.CreatedAt = &createdAt
		reports = append(reports, rep)
	}
	return reports, nil
}

// GetDevelopmentReportByID fetches full detail for a specific report.
func (r *ChildDevelopmentRepository) GetDevelopmentReportByID(ctx context.Context, reportID uuid.UUID) (*child_dev.DevelopmentReportResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var report child_dev.DevelopmentReportResponse

	// Fetch specific report by ID
	queryReport := `
		SELECT id, kpsp_score, month_target
		FROM child_development_reports
		WHERE id = $1
	`
	var monthTarget int

	err := r.pool.QueryRow(ctx, queryReport, reportID).Scan(&report.ID, &report.KPSPScore, &monthTarget)
	if err != nil {
		return nil, err
	}
	report.MonthTarget = &monthTarget

	// Fetch domains
	queryDomains := `
		SELECT 
			q.developmental_domain,
			COUNT(q.id) as total_question,
			COUNT(a.id) FILTER (WHERE a.assessment_kpsp_answer = true) as true_answer
		FROM assessment_kpsp_answers a
		JOIN assessment_kpsp_questions q ON a.assessment_kpsp_question_id = q.id
		WHERE a.child_development_report_id = $1
		GROUP BY q.developmental_domain
	`
	rows, err := r.pool.Query(ctx, queryDomains, reportID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var d child_dev.DomainAggregate
			if err := rows.Scan(&d.DevelopmentalDomain, &d.TotalQuestion, &d.TrueAnswer); err != nil {
				return nil, err
			}
			report.Domains = append(report.Domains, d)
		}
	}

	return &report, nil
}

// GetRecommendationsByDevelopmentalDomain fetches recommended actions for a given domain
func (r *ChildDevelopmentRepository) GetRecommendationsByDevelopmentalDomain(ctx context.Context, reportID uuid.UUID, domain string) ([]child_dev.RecommendedAction, error) {
	var recs []child_dev.RecommendedAction

	queryRecs := `
		SELECT ra.id, ra.assessment_kpsp_question_id, ra.title, ra.action_text
		FROM development_report_recommendations aar
		JOIN recommended_actions ra ON aar.recommended_action_id = ra.id
		JOIN assessment_kpsp_questions q ON ra.assessment_kpsp_question_id = q.id
		WHERE aar.child_development_report_id = $1 AND q.developmental_domain = $2
	`

	rows, err := r.pool.Query(ctx, queryRecs, reportID, domain)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rec child_dev.RecommendedAction
		if err := rows.Scan(&rec.ID, &rec.AssessmentKPSPQuestionID, &rec.Title, &rec.ActionText); err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

// GetKPSPQuestions returns KPSP questions filtered by monthTarget.
func (r *ChildDevelopmentRepository) GetKPSPQuestions(ctx context.Context, monthTarget int) ([]child_dev.AssessmentKPSPQuestion, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, "order", developmental_domain, month_target, question, created_at, updated_at
		FROM assessment_kpsp_questions
		WHERE month_target = $1
		ORDER BY "order" ASC`

	rows, err := r.pool.Query(ctx, query, monthTarget)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []child_dev.AssessmentKPSPQuestion
	for rows.Next() {
		var q child_dev.AssessmentKPSPQuestion
		if err := rows.Scan(
			&q.ID, &q.Order, &q.DevelopmentalDomain, &q.MonthTarget, 
			&q.Question, &q.CreatedAt, &q.UpdatedAt,
		); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// CreateDevelopmentReport creates report, answers, and recommendations in a single transaction.
func (r *ChildDevelopmentRepository) CreateDevelopmentReport(
	ctx context.Context,
	childID uuid.UUID,
	monthTarget int,
	kpspScore int,
	nextCheckDate *time.Time,
	answers []child_dev.KPSPAnswerInput,
) (uuid.UUID, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	reportID := uuid.New()
	queryReport := `
		INSERT INTO child_development_reports (id, child_id, kpsp_score, month_target, next_check_date)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, queryReport, reportID, childID, kpspScore, monthTarget, nextCheckDate)
	if err != nil {
		return uuid.Nil, err
	}

	// Insert answers and conditionally fetch recommendations for false answers
	queryAnswer := `
		INSERT INTO assessment_kpsp_answers (id, child_development_report_id, assessment_kpsp_question_id, assessment_kpsp_answer)
		VALUES ($1, $2, $3, $4)
	`
	queryRecFetch := `
		SELECT id
		FROM recommended_actions 
		WHERE assessment_kpsp_question_id = $1
	`
	queryRecInsert := `
		INSERT INTO development_report_recommendations (child_development_report_id, recommended_action_id)
		VALUES ($1, $2)`

	for _, ans := range answers {
		ansID := uuid.New()
		_, err = tx.Exec(ctx, queryAnswer, ansID, reportID, ans.AssessmentKPSPQuestionID, ans.Answer)
		if err != nil {
			return uuid.Nil, err
		}

		if !ans.Answer {
			// Find corresponding recommended action master data
			var recID uuid.UUID
			err = tx.QueryRow(ctx, queryRecFetch, ans.AssessmentKPSPQuestionID).Scan(&recID)
			if err == nil {
				// Insert attempt recommendation
				_, err = tx.Exec(ctx, queryRecInsert, reportID, recID)
				if err != nil {
					return uuid.Nil, err
				}
			} else if err != pgx.ErrNoRows {
				return uuid.Nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return reportID, nil
}

// UpdateDevelopmentReport updates answers in a transaction and recalculates score and recommendations.
func (r *ChildDevelopmentRepository) UpdateDevelopmentReport(
	ctx context.Context,
	reportID uuid.UUID,
	answers []child_dev.KPSPAnswerInput,
	nextCheckDate *time.Time,
) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update the answers. We just update existing ones based on question_id
	updateAnswerQuery := `
		UPDATE assessment_kpsp_answers 
		SET assessment_kpsp_answer = $1 
		WHERE child_development_report_id = $2 
			AND assessment_kpsp_question_id = $3
	`
	for _, ans := range answers {
		_, err = tx.Exec(ctx, updateAnswerQuery, ans.Answer, reportID, ans.AssessmentKPSPQuestionID)
		if err != nil {
			return err
		}
	}

	// Recalculate score
	var newScore int
	queryRecalculateScore := `
		SELECT COUNT(*) 
		FROM assessment_kpsp_answers 
		WHERE child_development_report_id = $1 
			AND assessment_kpsp_answer = true
	`
	err = tx.QueryRow(ctx, queryRecalculateScore, reportID).Scan(&newScore)
	if err != nil {
		return err
	}

	// Delete all existing recommendations for this report to regenerate them
	queryDeleteRecs := `
		DELETE FROM development_report_recommendations 
		WHERE child_development_report_id = $1
	`
	_, err = tx.Exec(ctx, queryDeleteRecs, reportID)
	if err != nil {
		return err
	}

	// Regenerate recommendations for current false answers
	queryFalseAnswers := `
		SELECT assessment_kpsp_question_id 
		FROM assessment_kpsp_answers 
		WHERE child_development_report_id = $1 
			AND assessment_kpsp_answer = false`
	rows, err := tx.Query(ctx, queryFalseAnswers, reportID)
	if err != nil {
		return err
	}
	var falseQuestions []uuid.UUID
	for rows.Next() {
		var qID uuid.UUID
		if err := rows.Scan(&qID); err != nil {
			rows.Close()
			return err
		}
		falseQuestions = append(falseQuestions, qID)
	}
	rows.Close()

	queryRecFetch := `
		SELECT id
		FROM recommended_actions 
		WHERE assessment_kpsp_question_id = $1
	`
	queryRecInsert := `
		INSERT INTO development_report_recommendations (child_development_report_id, recommended_action_id) 
		VALUES ($1, $2)
	`

	for _, qID := range falseQuestions {
		var recID uuid.UUID
		err = tx.QueryRow(ctx, queryRecFetch, qID).Scan(&recID)
		if err == nil {
			_, err = tx.Exec(ctx, queryRecInsert, reportID, recID)
			if err != nil {
				return err
			}
		}
	}

	// Update the report's score
	queryUpdateScore := `
		UPDATE child_development_reports 
		SET kpsp_score = $1, next_check_date = $2
		WHERE id = $3
	`
	_, err = tx.Exec(ctx, queryUpdateScore, newScore, nextCheckDate, reportID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeleteDevelopmentReport performs a hard delete on a development report.
func (r *ChildDevelopmentRepository) DeleteDevelopmentReport(ctx context.Context, reportID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	queryDeleteReport := `
		DELETE FROM child_development_reports 
		WHERE id = $1
	`
	res, err := r.pool.Exec(ctx, queryDeleteReport, reportID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

// GetChecklistMilestoneTasks returns checklist tasks with an is_checked boolean via LEFT JOIN.
func (r *ChildDevelopmentRepository) GetChecklistMilestoneTasks(ctx context.Context, childID uuid.UUID, monthTarget int) ([]child_dev.ChecklistMilestoneTaskResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			t.id, t.developmental_domain, t.task_description,
			CASE WHEN p.checklist_milestone_task_id IS NOT NULL THEN true ELSE false END as is_checked
		FROM checklist_milestone_tasks t
		LEFT JOIN checklist_milestone_progress p 
			ON t.id = p.checklist_milestone_task_id 
			AND p.child_id = $1
		WHERE t.month_target = $2
	`
	rows, err := r.pool.Query(ctx, query, childID, monthTarget)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []child_dev.ChecklistMilestoneTaskResponse
	for rows.Next() {
		var t child_dev.ChecklistMilestoneTaskResponse
		if err := rows.Scan(
			&t.ID, &t.DevelopmentalDomain, &t.TaskDescription, &t.IsChecked,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// UpdateChecklistMilestone performs a replace-all operation on checklist progress for a child.
func (r *ChildDevelopmentRepository) UpdateChecklistMilestone(ctx context.Context, childID uuid.UUID, taskIDs []uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if len(taskIDs) == 0 {
		// If empty, just delete all for this child
		queryDeleteAll := `
			DELETE FROM checklist_milestone_progress 
			WHERE child_id = $1
		`
		_, err = tx.Exec(ctx, queryDeleteAll, childID)
		if err != nil {
			return err
		}
	} else {
		// 1. Delete tasks that are NOT in the new list (un-checked by user)
		queryDeleteUnchecked := `
			DELETE FROM checklist_milestone_progress 
			WHERE child_id = $1 
				AND checklist_milestone_task_id != ALL($2)
		`
		_, err = tx.Exec(ctx, queryDeleteUnchecked, childID, taskIDs)
		if err != nil {
			return err
		}

		// 2. Insert the new ones, ignoring if they already exist (ON CONFLICT DO NOTHING)
		queryInsertProgress := `
			INSERT INTO checklist_milestone_progress (child_id, checklist_milestone_task_id) 
			VALUES ($1, $2)
			ON CONFLICT (child_id, checklist_milestone_task_id) DO NOTHING
		`
		for _, taskID := range taskIDs {
			_, err = tx.Exec(ctx, queryInsertProgress, childID, taskID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

// GetRecommendationsByReportID fetches all recommendation items for a given child_development_report_id.
func (r *ChildDevelopmentRepository) GetRecommendationsByReportID(ctx context.Context, reportID uuid.UUID) ([]child_dev.RecommendationItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT
			akq.developmental_domain,
			ra.action_text
		FROM development_report_recommendations aar
		JOIN recommended_actions ra
			ON ra.id = aar.recommended_action_id
		JOIN assessment_kpsp_questions akq
			ON akq.id = ra.assessment_kpsp_question_id
		WHERE aar.child_development_report_id = $1
		ORDER BY akq.developmental_domain
	`

	rows, err := r.pool.Query(ctx, query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []child_dev.RecommendationItem
	for rows.Next() {
		var item child_dev.RecommendationItem
		if err := rows.Scan(&item.DevelopmentalDomain, &item.ActionText); err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	if results == nil {
		results = []child_dev.RecommendationItem{}
	}
	return results, rows.Err()
}

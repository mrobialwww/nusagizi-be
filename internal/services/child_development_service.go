package services

import (
	"context"
	"fmt"
	"time"

	"nusagizi_be/internal/models"
	child_dev "nusagizi_be/internal/models/child_development"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/utils"

	"github.com/google/uuid"
)

type ChildDevelopmentService struct {
	repo        *repository.ChildDevelopmentRepository
	childRepo   *repository.ChildRepository
	motherRepo  *repository.MotherProfileRepository
	r2PublicURL string
}

func NewChildDevelopmentService(
	repo *repository.ChildDevelopmentRepository,
	childRepo *repository.ChildRepository,
	motherRepo *repository.MotherProfileRepository,
	r2PublicURL string,
) *ChildDevelopmentService {
	return &ChildDevelopmentService{
		repo:        repo,
		childRepo:   childRepo,
		motherRepo:  motherRepo,
		r2PublicURL: r2PublicURL,
	}
}

// GetLatestDevelopmentReport (Endpoint: 19)
func (s *ChildDevelopmentService) GetLatestDevelopmentReport(ctx context.Context, userID string, childID uuid.UUID) (*child_dev.DevelopmentReportResponse, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}

	report, err := s.repo.GetLatestDevelopmentReport(ctx, childID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, nil
	}
	status := DetermineKPSPStatus(report.KPSPScore)
	report.Status = &status
	return report, nil
}

// GetDevelopmentReports (Endpoint: 20)
func (s *ChildDevelopmentService) GetDevelopmentReports(ctx context.Context, userID string, childID uuid.UUID) ([]child_dev.DevelopmentReportResponse, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}
	reports, err := s.repo.GetDevelopmentReports(ctx, childID)
	if err != nil {
		return nil, err
	}
	for i := range reports {
		status := DetermineKPSPStatus(reports[i].KPSPScore)
		reports[i].Status = &status
	}
	return reports, nil
}

// GetDevelopmentReportByID (Endpoint: 21)
func (s *ChildDevelopmentService) GetDevelopmentReportByID(ctx context.Context, userID string, reportID uuid.UUID) (*child_dev.DevelopmentReportResponse, error) {
	// Fetch report first (child_id is embedded in the response)
	report, err := s.repo.GetDevelopmentReportByID(ctx, reportID)
	if err != nil {
		return nil, err
	}

	// Use child_id from the report itself for ownership check
	if report.ChildID == nil {
		return nil, fmt.Errorf("record not found")
	}
	if err := checkMotherOwnership(ctx, userID, *report.ChildID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}

	// Find the domain with the lowest ratio
	var worstNerve string
	minRatio := float64(2.0)

	for _, d := range report.Domains {
		if d.TotalQuestion > 0 {
			ratio := float64(d.TrueAnswer) / float64(d.TotalQuestion)
			if ratio < minRatio {
				minRatio = ratio
				worstNerve = d.DevelopmentalDomain
			}
		}
	}

	// Fetch recommendations only for the worst nerve
	if worstNerve != "" {
		recs, err := s.repo.GetRecommendationsByDevelopmentalDomain(ctx, reportID, worstNerve)
		if err == nil {
			report.RecommendedActions = recs
		}
	}

	status := DetermineKPSPStatus(report.KPSPScore)
	report.Status = &status
	return report, nil
}

// GetKPSPQuestions (Endpoint: 22 - Global Data, no childID required)
func (s *ChildDevelopmentService) GetKPSPQuestions(ctx context.Context, monthTarget int) ([]child_dev.AssessmentKPSPQuestion, error) {
	resp, err := s.repo.GetKPSPQuestions(ctx, monthTarget)
	if err != nil {
		return nil, err
	}
	for i := range resp {
		resp[i].ImageURL = utils.ResolvePublicPhotoURL(s.r2PublicURL, resp[i].ImageURL)
	}
	return resp, nil
}

// CreateDevelopmentReport (Endpoint: 23)
func (s *ChildDevelopmentService) CreateDevelopmentReport(ctx context.Context, userID string, childID uuid.UUID, input *child_dev.CreateDevelopmentReportInput) (uuid.UUID, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return uuid.Nil, err
	}

	if len(input.ListAnswer) != 10 {
		return uuid.Nil, fmt.Errorf("answers must be exactly 10")
	}

	childInfo, err := s.childRepo.GetChild(ctx, childID)
	if err != nil {
		return uuid.Nil, err
	}
	birthDate, err := time.Parse(models.DateLayout, childInfo.BirthDate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid birth_date format in db: %w", err)
	}

	expectedMonthTarget, nextCheckDate := calculateKPSPParams(birthDate)
	if input.MonthTarget != expectedMonthTarget {
		return uuid.Nil, fmt.Errorf("month_target mismatch, expected %d", expectedMonthTarget)
	}

	// Calculate score
	score := 0
	for _, ans := range input.ListAnswer {
		if ans.Answer {
			score++
		}
	}

	return s.repo.CreateDevelopmentReport(ctx, childID, input.MonthTarget, score, nextCheckDate, input.ListAnswer)
}

// UpdateDevelopmentReport (Endpoint: 24)
func (s *ChildDevelopmentService) UpdateDevelopmentReport(ctx context.Context, userID string, childID uuid.UUID, reportID uuid.UUID, input *child_dev.UpdateDevelopmentReportInput) (uuid.UUID, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return uuid.Nil, err
	}

	if len(input.ListAnswer) == 0 {
		return uuid.Nil, fmt.Errorf("answers cannot be empty")
	}

	childInfo, err := s.childRepo.GetChild(ctx, childID)
	if err != nil {
		return uuid.Nil, err
	}
	birthDate, err := time.Parse(models.DateLayout, childInfo.BirthDate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid birth_date format in db: %w", err)
	}

	_, nextCheckDate := calculateKPSPParams(birthDate)

	if _, err := s.repo.UpdateDevelopmentReport(ctx, reportID, input.ListAnswer, nextCheckDate); err != nil {
		return uuid.Nil, err
	}
	return reportID, nil
}

// DeleteDevelopmentReport (Endpoint: 28)
func (s *ChildDevelopmentService) DeleteDevelopmentReport(ctx context.Context, userID string, reportID uuid.UUID) error {
	report, err := s.repo.GetDevelopmentReportByID(ctx, reportID)
	if err != nil {
		return err
	}
	if report.ChildID == nil {
		return fmt.Errorf("record not found")
	}

	if err := checkMotherOwnership(ctx, userID, *report.ChildID, s.motherRepo, s.childRepo); err != nil {
		return err
	}

	return s.repo.DeleteDevelopmentReport(ctx, reportID)
}

// GetChecklistMilestoneTasks (Endpoint: 25)
func (s *ChildDevelopmentService) GetChecklistMilestoneTasks(ctx context.Context, userID string, childID uuid.UUID, monthTarget int) ([]child_dev.ChecklistMilestoneResponse, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetChecklistMilestoneTasks(ctx, childID, monthTarget)
}

// UpdateChecklistMilestone (Endpoint: 27)
func (s *ChildDevelopmentService) UpdateChecklistMilestone(ctx context.Context, userID string, childID uuid.UUID, questionIDs []uuid.UUID) error {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.UpdateChecklistMilestone(ctx, childID, questionIDs)
}

// GetRecommendations (Endpoint: 26)
func (s *ChildDevelopmentService) GetRecommendations(ctx context.Context, userID string, childID uuid.UUID, reportID uuid.UUID) ([]child_dev.RecommendationItem, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetRecommendationsByReportID(ctx, reportID)
}

// ================================== HELPER FUNCTION ===================================
// Calculate KPSP parameters (used for checking consistency and calculating next check date)
var kpspPeriods = []int{3, 6, 9, 12, 15, 18, 21, 24, 30, 36, 42, 48, 54, 60}

func calculateKPSPParams(birthDate time.Time) (int, *time.Time) {
	now := time.Now()

	// age_months = umur anak sekarang, dihitung dari birth_date
	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())
	ageMonths := years*12 + months
	if now.Day() < birthDate.Day() {
		ageMonths--
	}
	if ageMonths < 0 {
		ageMonths = 0
	}

	// month_target = periode TERBESAR di KPSP_PERIODS yang <= age_months
	// jika age_months < 3 -> month_target = 3
	monthTarget := 3
	for i := len(kpspPeriods) - 1; i >= 0; i-- {
		if kpspPeriods[i] <= ageMonths {
			monthTarget = kpspPeriods[i]
			break
		}
	}

	// next_period  = periode TERKECIL di KPSP_PERIODS yang > age_months
	var nextCheckDate *time.Time
	found := false
	var nextPeriod int
	for _, p := range kpspPeriods {
		if p > ageMonths {
			nextPeriod = p
			found = true
			break
		}
	}

	if found {
		// birth_date + next_period bulan (format DD-MM-YYYY)
		ncd := birthDate.AddDate(0, nextPeriod, 0)
		nextCheckDate = &ncd
	}

	return monthTarget, nextCheckDate
}

// Determine KPSP status based on score
func DetermineKPSPStatus(score int) string {
	if score >= 9 {
		return "Sesuai Usia"
	} else if score >= 7 {
		return "Perkembangan Meragukan"
	}
	return "Kemungkinan Penyimpangan"
}

package services

import (
	"context"
	"fmt"

	child_growth "nusagizi_be/internal/models/child_growth"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/who"

	"github.com/google/uuid"
)

type ChildGrowthService struct {
	repo       *repository.ChildGrowthRepository
	childRepo  *repository.ChildRepository
	motherRepo *repository.MotherProfileRepository
}

func NewChildGrowthService(
	repo *repository.ChildGrowthRepository,
	childRepo *repository.ChildRepository,
	motherRepo *repository.MotherProfileRepository,
) *ChildGrowthService {
	return &ChildGrowthService{
		repo:       repo,
		childRepo:  childRepo,
		motherRepo: motherRepo,
	}
}

// GetLatestGrowthReport (Endpoint: 13)
func (s *ChildGrowthService) GetLatestGrowthReport(ctx context.Context, userID string, childID uuid.UUID) (*child_growth.GrowthReportWithStatus, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}
	raw, err := s.repo.GetLatestGrowthReport(ctx, childID)
	if err != nil {
		return nil, err
	}
	status, desc := who.ComputeGrowthStatus(raw.Gender, raw.AgeMonths, raw.WeightKg, raw.HeightCm, raw.HeadCircumferenceCm)
	return &child_growth.GrowthReportWithStatus{
		ID:                  raw.ID,
		MeasuredAt:          raw.MeasuredAt,
		WeightKg:            raw.WeightKg,
		HeightCm:            raw.HeightCm,
		HeadCircumferenceCm: raw.HeadCircumferenceCm,
		Status:              status,
		Description:         desc,
	}, nil
}

// GetGrowthAnalyses (Endpoint: 14)
func (s *ChildGrowthService) GetGrowthAnalyses(ctx context.Context, userID string, childID uuid.UUID, analysisType string, ageRange string) (*child_growth.GrowthAnalysisResult, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}

	// Mandatory validation
	if analysisType == "" || ageRange == "" {
		return nil, fmt.Errorf("analysis_type and age_range are required")
	}

	validAnalysisTypes := map[string]bool{
		"weight_for_age":             true,
		"height_for_age":             true,
		"weight_for_height":          true,
		"bmi_for_age":                true,
		"head_circumference_for_age": true,
	}
	if !validAnalysisTypes[analysisType] {
		return nil, fmt.Errorf("invalid analysis_type")
	}

	ageRangeMap := map[string][2]int{
		"0-2":  {0, 2},
		"0-12": {0, 12},
		"0-60": {0, 60},
	}
	bounds, ok := ageRangeMap[ageRange]
	if !ok {
		return nil, fmt.Errorf("invalid age_range, must be one of: 0-2, 0-12, 0-60")
	}

	// - measurements: list of raw data (height, weight, head circumference, and age in months) filtered by age_range.
	// - gender: the child's gender ("male" or "female") used to determine the correct WHO reference table.
	measurements, gender, err := s.repo.GetGrowthMeasurements(ctx, childID, bounds[0], bounds[1], analysisType)
	if err != nil {
		return nil, err
	}

	result := &child_growth.GrowthAnalysisResult{
		DataPoints: make([]child_growth.GrowthDataPoint, 0, len(measurements)),
		Message:    nil,
	}

	if len(measurements) == 0 {
		return result, nil
	}

	// Convert raw measurements into the who.Measurement input type
	whoMeasurements := make([]who.Measurement, len(measurements))
	for i, m := range measurements {
		whoMeasurements[i] = who.Measurement{
			WeightKg:            m.WeightKg,
			HeightCm:            m.HeightCm,
			HeadCircumferenceCm: m.HeadCircumferenceCm,
			AgeMonths:           m.AgeMonths,
		}
	}

	// Delegate all WHO indicator logic to the who package
	points, msg := who.BuildAnalysisPoints(gender, analysisType, whoMeasurements)

	// Map results back to the domain model
	for _, p := range points {
		result.DataPoints = append(result.DataPoints, child_growth.GrowthDataPoint{X: p.X, Y: p.Y})
	}
	if msg != nil {
		result.Message = &child_growth.GrowthMessage{Title: msg.Title, Description: msg.Description}
	}

	return result, nil
}

// GetGrowthReports (Endpoint: 15)
func (s *ChildGrowthService) GetGrowthReports(ctx context.Context, userID string, childID uuid.UUID) ([]child_growth.GrowthReportWithStatus, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return nil, err
	}
	raws, err := s.repo.GetGrowthReports(ctx, childID)
	if err != nil {
		return nil, err
	}
	results := make([]child_growth.GrowthReportWithStatus, 0, len(raws))
	for _, raw := range raws {
		status, desc := who.ComputeGrowthStatus(raw.Gender, raw.AgeMonths, raw.WeightKg, raw.HeightCm, raw.HeadCircumferenceCm)
		results = append(results, child_growth.GrowthReportWithStatus{
			ID:                  raw.ID,
			MeasuredAt:          raw.MeasuredAt,
			WeightKg:            raw.WeightKg,
			HeightCm:            raw.HeightCm,
			HeadCircumferenceCm: raw.HeadCircumferenceCm,
			Status:              status,
			Description:         desc,
		})
	}
	return results, nil
}

// CreateGrowthReport (Endpoint: 16)
func (s *ChildGrowthService) CreateGrowthReport(ctx context.Context, userID string, childID uuid.UUID, input *child_growth.CreateGrowthReportInput) (uuid.UUID, error) {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return uuid.Nil, err
	}

	// Additional validation: all measurement fields are required and must be > 0
	if input.MeasuredAt == "" {
		return uuid.Nil, fmt.Errorf("measured_at is required")
	}
	if input.WeightKg == nil || *input.WeightKg <= 0 {
		return uuid.Nil, fmt.Errorf("weight_kg is required and must be > 0")
	}
	if input.HeightCm == nil || *input.HeightCm <= 0 {
		return uuid.Nil, fmt.Errorf("height_cm is required and must be > 0")
	}
	if input.HeadCircumferenceCm == nil || *input.HeadCircumferenceCm <= 0 {
		return uuid.Nil, fmt.Errorf("head_circumference_cm is required and must be > 0")
	}

	return s.repo.CreateGrowthReport(ctx, childID, input)
}

// UpdateGrowthReport (Endpoint: 17)
func (s *ChildGrowthService) UpdateGrowthReport(ctx context.Context, userID string, childID uuid.UUID, reportID uuid.UUID, input *child_growth.UpdateGrowthReportInput) error {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return err
	}

	// Additional validation
	if input.HeightCm != nil && *input.HeightCm <= 0 {
		return fmt.Errorf("height_cm must be > 0")
	}
	if input.WeightKg != nil && *input.WeightKg <= 0 {
		return fmt.Errorf("weight_kg must be > 0")
	}
	if input.HeadCircumferenceCm != nil && *input.HeadCircumferenceCm <= 0 {
		return fmt.Errorf("head_circumference_cm must be > 0")
	}

	return s.repo.UpdateGrowthReport(ctx, reportID, input)
}

// DeleteGrowthReport (Endpoint: 18)
func (s *ChildGrowthService) DeleteGrowthReport(ctx context.Context, userID string, childID uuid.UUID, reportID uuid.UUID) error {
	if err := checkMotherOwnership(ctx, userID, childID, s.motherRepo, s.childRepo); err != nil {
		return err
	}

	return s.repo.DeleteGrowthReport(ctx, reportID)
}

package services

import (
	"context"
	"fmt"
	"math"

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
	status, desc := computeOverallStatus(*raw)
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
		"0-6":  {0, 6},
		"0-60": {0, 60},
	}
	bounds, ok := ageRangeMap[ageRange]
	if !ok {
		return nil, fmt.Errorf("invalid age_range, must be one of: 0-2, 0-6, 0-60")
	}

	// GetGrowthMeasurements returns:
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

	// For the message, we only calculate Z-Score for the latest measurement
	for i, m := range measurements {
		var x, y float64
		var lms who.LMSRow
		var found bool

		isLatest := (i == len(measurements)-1)

		switch analysisType {
		case "weight_for_age":
			x = float64(m.AgeMonths)
			y = m.WeightKg
			if gender == "male" {
				lms, found = who.WFABoys[m.AgeMonths]
			} else {
				lms, found = who.WFAGirls[m.AgeMonths]
			}

		case "height_for_age":
			x = float64(m.AgeMonths)
			y = m.HeightCm
			if gender == "male" {
				lms, found = who.LHFABoys[m.AgeMonths]
			} else {
				lms, found = who.LHFAGirls[m.AgeMonths]
			}

		case "weight_for_height":
			roundedHeight := math.Round(m.HeightCm*2) / 2
			x = roundedHeight
			y = m.WeightKg
			if m.AgeMonths < 24 {
				if gender == "male" {
					lms, found = who.WFLBoys[roundedHeight]
				} else {
					lms, found = who.WFLGirls[roundedHeight]
				}
			} else {
				if gender == "male" {
					lms, found = who.WFHBoys[roundedHeight]
				} else {
					lms, found = who.WFHGirls[roundedHeight]
				}
			}

		case "bmi_for_age":
			x = float64(m.AgeMonths)
			if m.HeightCm > 0 {
				heightM := m.HeightCm / 100.0
				y = m.WeightKg / (heightM * heightM)
				if gender == "male" {
					lms, found = who.BFABoys[m.AgeMonths]
				} else {
					lms, found = who.BFAGirls[m.AgeMonths]
				}
			}

		case "head_circumference_for_age":
			x = float64(m.AgeMonths)
			y = m.HeadCircumferenceCm
			if gender == "male" {
				lms, found = who.HCFABoys[m.AgeMonths]
			} else {
				lms, found = who.HCFAGirls[m.AgeMonths]
			}
		}

		// Only append data points if the measurement value is positive.
		if y > 0 {
			result.DataPoints = append(result.DataPoints, child_growth.GrowthDataPoint{
				X: x,
				Y: y,
			})

			// Calculate Z-Score for the latest measurement
			if isLatest && found {
				latestZScore := who.CalculateZScore(y, lms.L, lms.M, lms.S)
				title, desc := who.ClassifyZScore(latestZScore, "anak")
				result.Message = &child_growth.GrowthMessage{
					Title:       title,
					Description: desc,
				}
			}
		}
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
		status, desc := computeOverallStatus(raw)
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

// computeOverallStatus calculates Z-Scores for all 5 WHO indicators from a single growth report
func computeOverallStatus(raw child_growth.RawGrowthReportFull) (status, description string) {
	var scores []float64

	var lms who.LMSRow
	var found bool

	isMale := raw.Gender == "male"
	age := raw.AgeMonths

	if raw.WeightKg > 0 {
		if isMale {
			lms, found = who.WFABoys[age]
		} else {
			lms, found = who.WFAGirls[age]
		}
		if found {
			scores = append(scores, who.CalculateZScore(raw.WeightKg, lms.L, lms.M, lms.S))
		}
	}

	if raw.HeightCm > 0 {
		if isMale {
			lms, found = who.LHFABoys[age]
		} else {
			lms, found = who.LHFAGirls[age]
		}
		if found {
			scores = append(scores, who.CalculateZScore(raw.HeightCm, lms.L, lms.M, lms.S))
		}
	}

	if raw.WeightKg > 0 && raw.HeightCm > 0 {
		roundedH := math.Round(raw.HeightCm*2) / 2 // round to nearest 0.5 cm
		if age < 24 {
			if isMale {
				lms, found = who.WFLBoys[roundedH]
			} else {
				lms, found = who.WFLGirls[roundedH]
			}
		} else {
			if isMale {
				lms, found = who.WFHBoys[roundedH]
			} else {
				lms, found = who.WFHGirls[roundedH]
			}
		}
		if found {
			scores = append(scores, who.CalculateZScore(raw.WeightKg, lms.L, lms.M, lms.S))
		}
	}

	if raw.WeightKg > 0 && raw.HeightCm > 0 {
		heightM := raw.HeightCm / 100
		bmi := raw.WeightKg / (heightM * heightM)
		if isMale {
			lms, found = who.BFABoys[age]
		} else {
			lms, found = who.BFAGirls[age]
		}
		if found {
			scores = append(scores, who.CalculateZScore(bmi, lms.L, lms.M, lms.S))
		}
	}

	if raw.HeadCircumferenceCm > 0 {
		if isMale {
			lms, found = who.HCFABoys[age]
		} else {
			lms, found = who.HCFAGirls[age]
		}
		if found {
			scores = append(scores, who.CalculateZScore(raw.HeadCircumferenceCm, lms.L, lms.M, lms.S))
		}
	}

	// If no indicator produced a score (e.g. age is outside all WHO table ranges),
	// return a neutral fallback instead of panicking or returning a misleading status.
	if len(scores) == 0 {
		return "Tidak Tersedia", "Data pengukuran tidak mencukupi untuk menghitung status pertumbuhan."
	}
	return who.ClassifyOverallStatus(scores)
}

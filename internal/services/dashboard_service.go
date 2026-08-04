package services

import (
	"context"
	"time"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/dashboard"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/who"
)

type DashboardService struct {
	repo       *repository.DashboardRepository
	motherRepo *repository.MotherProfileRepository
}

func NewDashboardService(
	repo *repository.DashboardRepository,
	motherRepo *repository.MotherProfileRepository,
) *DashboardService {
	return &DashboardService{
		repo:       repo,
		motherRepo: motherRepo,
	}
}

// GetDashboardSummary (endpoint: 64)
func (s *DashboardService) GetDashboardSummary(ctx context.Context, userID string) ([]dashboard.ChildSummaryResponse, error) {
	motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	results, err := s.repo.GetDashboardSummary(ctx, motherProfileID)
	if err != nil {
		return nil, err
	}

	for i := range results {
		// Calculate AgeMonths
		var ageMonths int
		birthDate, err := time.Parse(models.DateLayout, results[i].BirthDate)
		if err == nil {
			now := time.Now()
			years := now.Year() - birthDate.Year()
			months := int(now.Month()) - int(birthDate.Month())
			ageMonths = years*12 + months
			if now.Day() < birthDate.Day() {
				ageMonths--
			}
			if ageMonths < 0 {
				ageMonths = 0
			}
		}

		// Calculate Growth Status
		weight := 0.0
		if results[i].WeightKg != nil {
			weight = *results[i].WeightKg
		}
		height := 0.0
		if results[i].HeightCm != nil {
			height = *results[i].HeightCm
		}
		headCm := 0.0
		if results[i].HeadCircumferenceCm != nil {
			headCm = *results[i].HeadCircumferenceCm
		}
		status, _ := who.ComputeGrowthStatus(results[i].Gender, ageMonths, weight, height, headCm)
		results[i].StatusGrowth = status

		// Calculate Development Status
		if results[i].KpspScore != nil {
			results[i].StatusDevelopment = DetermineKPSPStatus(*results[i].KpspScore)
		} else {
			results[i].StatusDevelopment = "Tidak Diketahui"
		}
	}

	return results, nil
}

package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	child_nutri "nusagizi_be/internal/models/child_nutrition"
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
		birthDate := results[i].BirthDateRaw
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

		ageY := ageMonths / 12
		ageM := ageMonths % 12
		ageStr := ""
		if ageY > 0 {
			ageStr += fmt.Sprintf("%d tahun ", ageY)
		}
		if ageM > 0 {
			ageStr += fmt.Sprintf("%d bulan", ageM)
		}
		if ageStr == "" {
			ageStr = "0 bulan"
		}
		results[i].Age = strings.TrimSpace(ageStr)

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

		// Calculate Nutrition Status
		if results[i].TargetCalories != nil && results[i].Calories != nil &&
			results[i].TargetProtein != nil && results[i].Protein != nil &&
			results[i].TargetFat != nil && results[i].Fat != nil &&
			results[i].TargetCarbohydrate != nil && results[i].Carbohydrate != nil {
			
			// Reuse helper from child_nutrition_service by building a temporary struct
			dummyResp := &child_nutri.TodayNutritionReportResponse{
				Calories:           *results[i].Calories,
				TargetCalories:     *results[i].TargetCalories,
				Protein:            *results[i].Protein,
				TargetProtein:      *results[i].TargetProtein,
				Fat:                *results[i].Fat,
				TargetFat:          *results[i].TargetFat,
				Carbohydrate:       *results[i].Carbohydrate,
				TargetCarbohydrate: *results[i].TargetCarbohydrate,
			}
			results[i].StatusNutrition = DetermineNutritionStatus(dummyResp)
		} else {
			results[i].StatusNutrition = "Tidak Diketahui"
		}
	}

	return results, nil
}

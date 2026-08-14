package services

import (
	"context"
	"time"

	child_nutri "nusagizi_be/internal/models/child_nutrition"
	"nusagizi_be/internal/models/dashboard"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/utils"
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

// GetDashboardSummary (Endpoint: 64)
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
		birthDate := results[i].BirthDateRaw
		results[i].Age = utils.FormatAgeString(birthDate)

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

		growthMeasuredAt := time.Now()
		if results[i].GrowthMeasuredAt != nil {
			growthMeasuredAt = *results[i].GrowthMeasuredAt
		}
		growthAgeMonths := utils.CalculateAgeInMonths(birthDate, growthMeasuredAt)

		status, _ := who.ComputeGrowthStatus(results[i].Gender, growthAgeMonths, weight, height, headCm)
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

		// Calculate Combined Status
		redCount, yellowCount := 0, 0

		// countColor maps the status string to its respective severity color
		countColor := func(s string) bool {
			switch s {
			case "Sangat Buruk", "Kemungkinan Penyimpangan":
				redCount++
			case "Berisiko", "Perkembangan Meragukan", "Kurang Optimal":
				yellowCount++
			case "Normal", "Sesuai Usia":
				// Green indicators do not increase the severity count
			default:
				return false
			}
			return true
		}

		// Evaluate all 3 domains. Short-circuits if any data is missing.
		if countColor(results[i].StatusGrowth) && countColor(results[i].StatusDevelopment) && countColor(results[i].StatusNutrition) {
			// Determine final status based on the highest severity counts
			switch redCount {
			case 3:
				results[i].Status = "Perlu Konsultasi"
			case 2:
				results[i].Status = "Perlu Pendampingan"
			case 1:
				results[i].Status = "Perlu Perhatian"
			default:
				// If there are no red indicators, evaluate the yellow counts
				if yellowCount == 3 {
					results[i].Status = "Perlu Perhatian"
				} else if yellowCount >= 1 {
					results[i].Status = "Perkembangan Baik"
				} else {
					results[i].Status = "Tumbuh Optimal"
				}
			}
		} else {
			results[i].Status = "Data Belum Lengkap"
		}
	}

	return results, nil
}

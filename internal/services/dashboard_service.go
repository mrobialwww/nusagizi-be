package services

import (
	"context"
	"time"

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
		// Calculate Valid Streak
		if results[i].LastStreakDateRaw != nil {
			today := time.Now().UTC().Truncate(24 * time.Hour)
			lastDate := (*results[i].LastStreakDateRaw).UTC().Truncate(24 * time.Hour)
			diff := today.Sub(lastDate).Hours() / 24
			if diff > 1 {
				results[i].Streak = 0
			}
		} else {
			results[i].Streak = 0
		}

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

		calories := 0.0
		if results[i].Calories != nil {
			calories = *results[i].Calories
		}
		targetCalories := 0.0
		if results[i].TargetCalories != nil {
			targetCalories = *results[i].TargetCalories
		}
		protein := 0.0
		if results[i].Protein != nil {
			protein = *results[i].Protein
		}
		targetProtein := 0.0
		if results[i].TargetProtein != nil {
			targetProtein = *results[i].TargetProtein
		}
		fat := 0.0
		if results[i].Fat != nil {
			fat = *results[i].Fat
		}
		targetFat := 0.0
		if results[i].TargetFat != nil {
			targetFat = *results[i].TargetFat
		}
		carbohydrate := 0.0
		if results[i].Carbohydrate != nil {
			carbohydrate = *results[i].Carbohydrate
		}
		targetCarbohydrate := 0.0
		if results[i].TargetCarbohydrate != nil {
			targetCarbohydrate = *results[i].TargetCarbohydrate
		}
		results[i].StatusNutrition = DetermineNutritionStatus(
			calories, targetCalories,
			protein, targetProtein,
			fat, targetFat,
			carbohydrate, targetCarbohydrate,
		)

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

package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	child_nutri "nusagizi_be/internal/models/child_nutrition"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/utils"

	"github.com/google/uuid"
)

type ChildNutritionService struct {
	repo          *repository.ChildNutritionRepository
	childRepo     *repository.ChildRepository
	motherRepo    *repository.MotherProfileRepository
	caregiverRepo *repository.CaregiverRepository
	r2PublicURL   string
}

func NewChildNutritionService(
	repo *repository.ChildNutritionRepository,
	childRepo *repository.ChildRepository,
	motherRepo *repository.MotherProfileRepository,
	caregiverRepo *repository.CaregiverRepository,
	r2PublicURL string,
) *ChildNutritionService {
	return &ChildNutritionService{
		repo:          repo,
		childRepo:     childRepo,
		motherRepo:    motherRepo,
		caregiverRepo: caregiverRepo,
		r2PublicURL:   r2PublicURL,
	}
}

// GetTodayNutritionReport (Endpoint: 29)
func (s *ChildNutritionService) GetTodayNutritionReport(ctx context.Context, userID string, childID uuid.UUID, reportDate time.Time) (*child_nutri.TodayNutritionReportResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}

	resp, err := s.repo.GetTodayNutritionReport(ctx, childID, reportDate)
	if err != nil {
		return nil, err
	}

	// Add flag CanRegenerate based on IsReportReused
	isReused, _ := s.repo.IsReportReused(ctx, resp.ID)
	resp.CanRegenerate = !isReused

	if resp != nil {
		resp.Status = DetermineNutritionStatus(resp.Calories, resp.TargetCalories, resp.Protein, resp.TargetProtein, resp.Fat, resp.TargetFat, resp.Carbohydrate, resp.TargetCarbohydrate)
		// Resolve images
		for i := range resp.Menu.Recipes {
			resp.Menu.Recipes[i].ImageURL = utils.ResolvePublicPhotoURL(s.r2PublicURL, resp.Menu.Recipes[i].ImageURL)
		}
	}
	return resp, nil
}

// GetReportMenuByID (Endpoint: 40)
func (s *ChildNutritionService) GetReportMenuByID(ctx context.Context, userID string, childID uuid.UUID, reportID uuid.UUID) (*child_nutri.MenuResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	resp, err := s.repo.GetReportMenuByID(ctx, childID, reportID)
	if err != nil {
		return nil, err
	}
	for i := range resp.Recipes {
		resp.Recipes[i].ImageURL = utils.ResolvePublicPhotoURL(s.r2PublicURL, resp.Recipes[i].ImageURL)
	}
	return resp, nil
}

// UpdateRecipeCompleteStatus (Endpoint: 31)
func (s *ChildNutritionService) UpdateRecipeCompleteStatus(ctx context.Context, userID string, recipeID uuid.UUID, portionsConsumed float64, reportDate time.Time) error {
	childID, err := s.repo.GetChildIDByRecipeID(ctx, recipeID)
	if err != nil {
		return err
	}
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}

	todayReport, err := s.repo.GetReportByChildAndDate(ctx, childID, reportDate)
	if err != nil {
		return fmt.Errorf("failed to get nutrition report for the specified date: %w", err)
	}

	return s.repo.UpdateRecipeCompleteStatus(ctx, recipeID, todayReport.ID, portionsConsumed)
}

// UpdateRecipeBookmarkStatus (Endpoint: 32)
func (s *ChildNutritionService) UpdateRecipeBookmarkStatus(ctx context.Context, userID string, recipeID uuid.UUID, isBookmarked bool) error {
	childID, err := s.repo.GetChildIDByRecipeID(ctx, recipeID)
	if err != nil {
		return err
	}
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.UpdateRecipeBookmarkStatus(ctx, recipeID, isBookmarked)
}

// GetRecipeDetail (Endpoint: 33)
func (s *ChildNutritionService) GetRecipeDetail(ctx context.Context, userID string, recipeID uuid.UUID) (*child_nutri.RecipeDetailResponse, error) {
	childID, err := s.repo.GetChildIDByRecipeID(ctx, recipeID)
	if err != nil {
		return nil, err
	}
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	resp, err := s.repo.GetRecipeDetail(ctx, recipeID)
	if err != nil {
		return nil, err
	}

	resp.ImageURL = utils.ResolvePublicPhotoURL(s.r2PublicURL, resp.ImageURL)
	for i := range resp.MainIngredients {
		resp.MainIngredients[i].Ingredient.ImageURL = utils.ResolvePublicPhotoURL(s.r2PublicURL, resp.MainIngredients[i].Ingredient.ImageURL)
	}
	for i := range resp.CookingSteps {
		if resolved := utils.ResolvePublicPhotoURL(s.r2PublicURL, &resp.CookingSteps[i].ImageURL); resolved != nil {
			resp.CookingSteps[i].ImageURL = *resolved
		}
	}
	return resp, nil
}

// SwapMainIngredientPriority (Endpoint: 34)
func (s *ChildNutritionService) SwapMainIngredientPriority(ctx context.Context, userID string, requests []child_nutri.SwapPriorityRequest) error {
	// Phase 1: Pre-flight check (Access Validation)
	recipeChildMap := make(map[uuid.UUID]uuid.UUID) // Cache childID per recipeID to avoid redundant queries for duplicate recipeIDs.

	for _, req := range requests {
		if _, exists := recipeChildMap[req.RecipeID]; !exists {
			childID, err := s.repo.GetChildIDByRecipeID(ctx, req.RecipeID)
			if err != nil {
				return err
			}
			if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
				return err
			}
			recipeChildMap[req.RecipeID] = childID
		}
	}

	// Phase 2: Atomic Execution
	return s.repo.SwapMainIngredientPriority(ctx, requests)
}

// GetBookmarkedRecipes (Endpoint: 35)
func (s *ChildNutritionService) GetBookmarkedRecipes(ctx context.Context, userID string, childID uuid.UUID) ([]child_nutri.RecipeResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	resp, err := s.repo.GetBookmarkedRecipes(ctx, childID)
	if err != nil {
		return nil, err
	}
	for i := range resp {
		resp[i].ImageURL = utils.ResolvePublicPhotoURL(s.r2PublicURL, resp[i].ImageURL)
	}
	return resp, nil
}

// GetNutritionReportsByMonth (Endpoint: 36)
func (s *ChildNutritionService) GetNutritionReportsByMonth(ctx context.Context, userID string, childID uuid.UUID, month int, year int) ([]child_nutri.NutritionReportSummaryResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month")
	}
	if year < 2000 {
		return nil, fmt.Errorf("invalid year")
	}
	results, err := s.repo.GetNutritionReportsByMonth(ctx, childID, month, year)
	if err != nil {
		return nil, err
	}
	for i := range results {
		res := &results[i]
		res.Status = DetermineNutritionStatus(res.Calories, res.TargetCalories, res.Protein, res.TargetProtein, res.Fat, res.TargetFat, res.Carbohydrate, res.TargetCarbohydrate)
	}
	return results, nil
}

// GetDailyShopIngredients (Endpoint: 37)
func (s *ChildNutritionService) GetDailyShopIngredients(ctx context.Context, userID string, targetDate time.Time) ([]child_nutri.DailyShopIngredient, error) {
	// Try mother profile first
	if motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID); err == nil {
		return s.repo.GetDailyShopByMother(ctx, motherProfileID, targetDate)
	}

	// If not mother, try caregiver profile
	if caregiverProfileID, err := s.caregiverRepo.GetByUserID(ctx, userID); err == nil {
		return s.repo.GetDailyShopByCaregiver(ctx, caregiverProfileID, targetDate)
	}

	// If neither, return error
	return nil, fmt.Errorf("%w: user has no access (must be mother or caregiver)", repository.ErrForbidden)
}

// GetTodayMenuShopping (Endpoint: 38)
func (s *ChildNutritionService) GetTodayMenuShopping(ctx context.Context, userID string, childID uuid.UUID, targetDate time.Time) (*child_nutri.MenuShoppingResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetTodayMenuShopping(ctx, childID, targetDate)
}

// ReuseRecipe (Endpoint: 30)
func (s *ChildNutritionService) ReuseRecipe(ctx context.Context, requesterID string, childID uuid.UUID, sourceRecipeID uuid.UUID, targetDate time.Time) error {
	if err := checkChildAccess(ctx, requesterID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}

	// Get target CNR
	targetCNR, err := s.repo.GetReportByChildAndDate(ctx, childID, targetDate)
	if err != nil {
		return fmt.Errorf("failed to get target nutrition report: %w", err)
	}
	if targetCNR == nil {
		return errors.New("nutrition report not found for the specified date")
	}

	// Ensure the source recipe actually belongs to the target child
	sourceChildID, err := s.repo.GetChildIDByRecipeID(ctx, sourceRecipeID)
	if err != nil {
		return fmt.Errorf("failed to verify source recipe: %w", err)
	}
	if sourceChildID != childID {
		return errors.New("invalid source recipe: recipe does not belong to this child")
	}

	// Get source recipe meal_time
	targetMealTime, err := s.repo.GetRecipeMealTime(ctx, sourceRecipeID)
	if err != nil {
		return fmt.Errorf("failed to get source recipe meal time: %w", err)
	}

	// Execute Repo Tx
	return s.repo.ReuseRecipeTx(ctx, targetCNR.ID, sourceRecipeID, targetMealTime)
}

// ================================== HELPER FUNCTION ===================================
// DetermineNutritionStatus evaluates macro targets and returns a descriptive status.
// (4 macros >= 90% target = Normal, 3 = Kurang Optimal, 2 = Berisiko, <=1 = Sangat Buruk)
func DetermineNutritionStatus(calories, targetCalories, protein, targetProtein, fat, targetFat, carbohydrate, targetCarbohydrate float64) string {
	count := 0
	if targetCalories > 0 && calories >= 0.9*targetCalories {
		count++
	}
	if targetProtein > 0 && protein >= 0.9*targetProtein {
		count++
	}
	if targetFat > 0 && fat >= 0.9*targetFat {
		count++
	}
	if targetCarbohydrate > 0 && carbohydrate >= 0.9*targetCarbohydrate {
		count++
	}

	switch count {
	case 4:
		return "Normal"
	case 3:
		return "Kurang Optimal"
	case 2:
		return "Berisiko"
	default:
		return "Sangat Buruk"
	}
}

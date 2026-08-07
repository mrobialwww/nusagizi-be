package services

import (
	"context"
	"fmt"

	child_nutri "nusagizi_be/internal/models/child_nutrition"
	"nusagizi_be/internal/repository"

	"github.com/google/uuid"
)

type ChildNutritionService struct {
	repo          *repository.ChildNutritionRepository
	childRepo     *repository.ChildRepository
	motherRepo    *repository.MotherProfileRepository
	caregiverRepo *repository.CaregiverRepository
}

func NewChildNutritionService(
	repo *repository.ChildNutritionRepository,
	childRepo *repository.ChildRepository,
	motherRepo *repository.MotherProfileRepository,
	caregiverRepo *repository.CaregiverRepository,
) *ChildNutritionService {
	return &ChildNutritionService{
		repo:          repo,
		childRepo:     childRepo,
		motherRepo:    motherRepo,
		caregiverRepo: caregiverRepo,
	}
}

// GetTodayNutritionReport (Endpoint: 29)
func (s *ChildNutritionService) GetTodayNutritionReport(ctx context.Context, userID string, childID uuid.UUID) (*child_nutri.TodayNutritionReportResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	resp, err := s.repo.GetTodayNutritionReport(ctx, childID)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		resp.Status = DetermineNutritionStatus(resp)
	}
	return resp, nil
}

// GetTodayDailyMenu (Endpoint: 30)
func (s *ChildNutritionService) GetTodayDailyMenu(ctx context.Context, userID string, childID uuid.UUID) (*child_nutri.MenuResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetTodayDailyMenu(ctx, childID)
}

// UpdateRecipeCompleteStatus (Endpoint: 31)
func (s *ChildNutritionService) UpdateRecipeCompleteStatus(ctx context.Context, userID string, recipeID uuid.UUID, portionsConsumed float64) error {
	childID, err := s.repo.GetChildIDByRecipeID(ctx, recipeID)
	if err != nil {
		return err
	}
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.UpdateRecipeCompleteStatus(ctx, recipeID, portionsConsumed)
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
	return s.repo.GetRecipeDetail(ctx, recipeID)
}

// SwapMainIngredientPriority (Endpoint: 34)
func (s *ChildNutritionService) SwapMainIngredientPriority(ctx context.Context, userID string, recipeID uuid.UUID, slot string, priority int) error {
	childID, err := s.repo.GetChildIDByRecipeID(ctx, recipeID)
	if err != nil {
		return err
	}
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.SwapMainIngredientPriority(ctx, recipeID, slot, priority)
}

// GetBookmarkedRecipes (Endpoint: 35)
func (s *ChildNutritionService) GetBookmarkedRecipes(ctx context.Context, userID string, childID uuid.UUID) ([]child_nutri.RecipeResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetBookmarkedRecipes(ctx, childID)
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
	return s.repo.GetNutritionReportsByMonth(ctx, childID, month, year)
}

// GetDailyShopIngredients (Endpoint: 37)
func (s *ChildNutritionService) GetDailyShopIngredients(ctx context.Context, userID string) ([]child_nutri.DailyShopIngredient, error) {
	// 1. Coba mother profile terlebih dahulu
	if motherProfileID, err := s.motherRepo.GetByUserID(ctx, userID); err == nil {
		return s.repo.GetDailyShopByMother(ctx, motherProfileID)
	}

	// 2. Jika bukan mother, coba caregiver profile
	if caregiverProfileID, err := s.caregiverRepo.GetByUserID(ctx, userID); err == nil {
		return s.repo.GetDailyShopByCaregiver(ctx, caregiverProfileID)
	}

	// 3. Jika bukan keduanya, return error
	return nil, fmt.Errorf("%w: user has no access (must be mother or caregiver)", repository.ErrForbidden)
}

// GetTodayMenuShopping (Endpoint: 38)
func (s *ChildNutritionService) GetTodayMenuShopping(ctx context.Context, userID string, childID uuid.UUID) (*child_nutri.MenuShoppingResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetTodayMenuShopping(ctx, childID)
}

// ================================== HELPER FUNCTION ===================================
// DetermineNutritionStatus evaluates macro targets and returns a descriptive status
func DetermineNutritionStatus(resp *child_nutri.TodayNutritionReportResponse) string {
	if resp == nil {
		return "sangat buruk"
	}
	count := 0
	if resp.TargetCalories > 0 && float64(resp.Calories) >= 0.9*float64(resp.TargetCalories) {
		count++
	}
	if resp.TargetProtein > 0 && float64(resp.Protein) >= 0.9*float64(resp.TargetProtein) {
		count++
	}
	if resp.TargetFat > 0 && float64(resp.Fat) >= 0.9*float64(resp.TargetFat) {
		count++
	}
	if resp.TargetCarbohydrate > 0 && float64(resp.Carbohydrate) >= 0.9*float64(resp.TargetCarbohydrate) {
		count++
	}

	switch count {
	case 4:
		return "normal"
	case 3:
		return "kurang optimal"
	case 2:
		return "beresiko"
	default:
		return "sangat buruk"
	}
}

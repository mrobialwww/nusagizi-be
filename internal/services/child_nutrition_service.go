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

// GetTodayNutritionReport (Endpoint: 28)
func (s *ChildNutritionService) GetTodayNutritionReport(ctx context.Context, userID string, childID uuid.UUID) (*child_nutri.TodayNutritionReportResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetTodayNutritionReport(ctx, childID)
}

// GetTodayDailyMenu (Endpoint: 29)
func (s *ChildNutritionService) GetTodayDailyMenu(ctx context.Context, userID string, childID uuid.UUID) (*child_nutri.MenuResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetTodayDailyMenu(ctx, childID)
}

// UpdateRecipeCompleteStatus (Endpoint: 30)
func (s *ChildNutritionService) UpdateRecipeCompleteStatus(ctx context.Context, userID string, childID uuid.UUID, recipeID uuid.UUID, isCompleted bool) error {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.UpdateRecipeCompleteStatus(ctx, recipeID, isCompleted)
}

// UpdateRecipeBookmarkStatus (Endpoint: 31)
func (s *ChildNutritionService) UpdateRecipeBookmarkStatus(ctx context.Context, userID string, childID uuid.UUID, recipeID uuid.UUID, isBookmarked bool) error {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.UpdateRecipeBookmarkStatus(ctx, recipeID, isBookmarked)
}

// GetRecipeDetail (Endpoint: 32)
func (s *ChildNutritionService) GetRecipeDetail(ctx context.Context, userID string, childID uuid.UUID, recipeID uuid.UUID) (*child_nutri.RecipeDetailResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetRecipeDetail(ctx, recipeID)
}

// SwapMainIngredientPriority (Endpoint: 33)
func (s *ChildNutritionService) SwapMainIngredientPriority(ctx context.Context, userID string, childID uuid.UUID, recipeID uuid.UUID, mainIngredientID uuid.UUID, slot string) error {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	if slot == "" {
		return fmt.Errorf("slot cannot be empty")
	}
	return s.repo.SwapMainIngredientPriority(ctx, recipeID, mainIngredientID, slot)
}

// GetBookmarkedRecipes (Endpoint: 34)
func (s *ChildNutritionService) GetBookmarkedRecipes(ctx context.Context, userID string, childID uuid.UUID) ([]child_nutri.RecipeResponse, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return nil, err
	}
	return s.repo.GetBookmarkedRecipes(ctx, childID)
}


// GetNutritionReportsByMonth (Endpoint: 35)
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

// UpdateDailyShoppingStatus (Endpoint: 36)
func (s *ChildNutritionService) UpdateDailyShoppingStatus(ctx context.Context, userID string, childID uuid.UUID, shoppingID uuid.UUID, isCompleted bool) error {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.UpdateDailyShoppingStatus(ctx, shoppingID, isCompleted)
}

// AddShoppingItem (Endpoint: 37)
func (s *ChildNutritionService) AddShoppingItem(ctx context.Context, userID string, childID uuid.UUID, shoppingID uuid.UUID, input *child_nutri.ShoppingItemInput) (uuid.UUID, error) {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return uuid.Nil, err
	}
	if input.Quantity <= 0 {
		return uuid.Nil, fmt.Errorf("quantity must be greater than 0")
	}
	if input.Unit == "" {
		return uuid.Nil, fmt.Errorf("unit cannot be empty")
	}
	return s.repo.AddShoppingItem(ctx, shoppingID, input)
}

// UpdateShoppingItem (Endpoint: 38)
func (s *ChildNutritionService) UpdateShoppingItem(ctx context.Context, userID string, childID uuid.UUID, itemID uuid.UUID, input *child_nutri.UpdateShoppingItemInput) error {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	if input.Quantity != nil && *input.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	if input.Unit != nil && *input.Unit == "" {
		return fmt.Errorf("unit cannot be empty")
	}
	return s.repo.UpdateShoppingItem(ctx, itemID, input)
}

// DeleteShoppingItem (Endpoint: 39)
func (s *ChildNutritionService) DeleteShoppingItem(ctx context.Context, userID string, childID uuid.UUID, itemID uuid.UUID) error {
	if err := checkChildAccess(ctx, userID, childID, s.motherRepo, s.caregiverRepo, s.childRepo); err != nil {
		return err
	}
	return s.repo.DeleteShoppingItem(ctx, itemID)
}

// SearchIngredients (Endpoint: 40 - Global data, no childID required)
func (s *ChildNutritionService) SearchIngredients(ctx context.Context, nameQuery string) ([]child_nutri.Ingredient, error) {
	if len(nameQuery) < 2 {
		return nil, fmt.Errorf("search query must be at least 2 characters long")
	}
	return s.repo.SearchIngredients(ctx, nameQuery)
}

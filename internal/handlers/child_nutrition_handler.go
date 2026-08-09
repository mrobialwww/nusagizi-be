package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChildNutritionHandler struct {
	service *services.ChildNutritionService
}

func NewChildNutritionHandler(service *services.ChildNutritionService) *ChildNutritionHandler {
	return &ChildNutritionHandler{service: service}
}

// GetDailyShopIngredients (Endpoint: 37)
func (h *ChildNutritionHandler) GetDailyShopIngredients(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	resp, err := h.service.GetDailyShopIngredients(c.Request.Context(), requester.ID)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this resource"}})
			return
		}
		slog.Error("GetDailyShopIngredients failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTodayMenuShopping (Endpoint: 38)
func (h *ChildNutritionHandler) GetTodayMenuShopping(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childIDStr := c.Param("child_id")
	childID, err := uuid.Parse(childIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetTodayMenuShopping(c.Request.Context(), requester.ID, childID)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this child"}})
			return
		}
		slog.Error("GetTodayMenuShopping failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTodayNutritionReport (Endpoint: 29)
func (h *ChildNutritionHandler) GetTodayNutritionReport(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childIDStr := c.Param("child_id")
	childID, err := uuid.Parse(childIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetTodayNutritionReport(c.Request.Context(), requester.ID, childID)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this child"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetTodayDailyMenu (Endpoint: 30)
func (h *ChildNutritionHandler) GetTodayDailyMenu(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childIDStr := c.Param("child_id")
	childID, err := uuid.Parse(childIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetTodayDailyMenu(c.Request.Context(), requester.ID, childID)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this child"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDailyMenuByID (Endpoint: 40)
func (h *ChildNutritionHandler) GetDailyMenuByID(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childIDStr := c.Param("child_id")
	childID, err := uuid.Parse(childIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	dailyMenuIDStr := c.Param("daily_menu_id")
	dailyMenuID, err := uuid.Parse(dailyMenuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid daily_menu_id format"}})
		return
	}

	resp, err := h.service.GetDailyMenuByID(c.Request.Context(), requester.ID, childID, dailyMenuID)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this child"}})
			return
		}
		if err.Error() == "not found" || strings.Contains(err.Error(), "no rows") || errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "daily menu not found"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateRecipeCompleteStatus (Endpoint: 31)
func (h *ChildNutritionHandler) UpdateRecipeCompleteStatus(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	recipeIDStr := c.Param("recipe_id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid recipe_id format"}})
		return
	}

	// Get portions_consumed from request body
	var reqBody struct {
		PortionsConsumed float64 `json:"portions_consumed"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid request body"}})
		return
	}

	// Validate portions_consumed
	switch reqBody.PortionsConsumed {
	case 0, 0.33, 0.66, 1:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "portions_consumed must be one of: 0, 0.33, 0.66, 1"}})
		return
	}

	err = h.service.UpdateRecipeCompleteStatus(c.Request.Context(), requester.ID, recipeID, reqBody.PortionsConsumed)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this recipe"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.Status(http.StatusOK)
}

// UpdateRecipeBookmarkStatus (Endpoint: 32)
func (h *ChildNutritionHandler) UpdateRecipeBookmarkStatus(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	recipeIDStr := c.Param("recipe_id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid recipe_id format"}})
		return
	}

	// Get is_bookmarked from request body
	var reqBody struct {
		IsBookmarked bool `json:"is_bookmarked"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid request body"}})
		return
	}

	err = h.service.UpdateRecipeBookmarkStatus(c.Request.Context(), requester.ID, recipeID, reqBody.IsBookmarked)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this recipe"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.Status(http.StatusOK)
}

// GetRecipeDetail (Endpoint: 33)
func (h *ChildNutritionHandler) GetRecipeDetail(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	recipeIDStr := c.Param("recipe_id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid recipe_id format"}})
		return
	}

	resp, err := h.service.GetRecipeDetail(c.Request.Context(), requester.ID, recipeID)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this recipe"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SwapMainIngredientPriority (Endpoint: 34)
func (h *ChildNutritionHandler) SwapMainIngredientPriority(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	recipeIDStr := c.Param("recipe_id")
	recipeID, err := uuid.Parse(recipeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid recipe_id format"}})
		return
	}

	// Get slot and priority from request body
	var reqBody struct {
		Slot     string `json:"slot"`
		Priority int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid request body"}})
		return
	}

	// Validate requested slot and priority
	if reqBody.Slot == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "slot is required"}})
		return
	}
	if reqBody.Priority <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "priority must be greater than 1 to be swapped to main ingredient"}})
		return
	}

	err = h.service.SwapMainIngredientPriority(c.Request.Context(), requester.ID, recipeID, reqBody.Slot, reqBody.Priority)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this recipe"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.Status(http.StatusOK)
}

// GetBookmarkedRecipes (Endpoint: 35)
func (h *ChildNutritionHandler) GetBookmarkedRecipes(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childIDStr := c.Param("child_id")
	childID, err := uuid.Parse(childIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetBookmarkedRecipes(c.Request.Context(), requester.ID, childID)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this child"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetNutritionReportsByMonth (Endpoint: 36)
func (h *ChildNutritionHandler) GetNutritionReportsByMonth(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childIDStr := c.Param("child_id")
	childID, err := uuid.Parse(childIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	// Get month and year from query params
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	// Parse month and year query params into integers
	var month, year int
	if _, err := fmt.Sscanf(monthStr, "%d", &month); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid month format"}})
		return
	}
	if _, err := fmt.Sscanf(yearStr, "%d", &year); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid year format"}})
		return
	}

	resp, err := h.service.GetNutritionReportsByMonth(c.Request.Context(), requester.ID, childID, month, year)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this child"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, resp)
}

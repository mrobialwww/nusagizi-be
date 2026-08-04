package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/dashboard"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetDashboardSummary (endpoint: 64)
func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	summary, err := h.service.GetDashboardSummary(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		default:
			slog.Error("GetDashboardSummary failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if summary == nil {
		summary = []dashboard.ChildSummaryResponse{}
	}
	c.JSON(http.StatusOK, summary)
}

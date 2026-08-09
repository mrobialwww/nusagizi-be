package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/caregiver"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaregiverHandler struct {
	service *services.CaregiverService
}

func NewCaregiverHandler(service *services.CaregiverService) *CaregiverHandler {
	return &CaregiverHandler{service: service}
}

// GetMotherCaregiverEngagements (endpoint: 52)
func (h *CaregiverHandler) GetCaregiverEngagements(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	engagements, err := h.service.GetCaregiverEngagements(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		default:
			slog.Error("GetMotherCaregiverEngagements failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if engagements == nil {
		engagements = []caregiver.CaregiverEngagementResponse{}
	}
	c.JSON(http.StatusOK, engagements)
}

// GetMotherCaregiverEngagementsRevoked (endpoint: 53)
func (h *CaregiverHandler) GetCaregiverEngagementsRevoked(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	engagements, err := h.service.GetRevokedEngagements(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		default:
			slog.Error("GetMotherCaregiverEngagementsRevoked failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if engagements == nil {
		engagements = []caregiver.CaregiverEngagementResponse{}
	}
	c.JSON(http.StatusOK, engagements)
}

// DeleteCaregiverEngagement (endpoint: 54)
func (h *CaregiverHandler) DeleteCaregiverEngagement(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	engagementID, err := uuid.Parse(c.Param("caregiver_engagement_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid caregiver_engagement_id format"}})
		return
	}

	err = h.service.DeleteCaregiverEngagement(c.Request.Context(), requester.ID, engagementID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have permission to delete this engagement"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "engagement not found"}})
		default:
			slog.Error("DeleteCaregiverEngagement failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetCaregiverChildren (endpoint: 56)
func (h *CaregiverHandler) GetCaregiverChildren(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	children, err := h.service.GetCaregiverChildren(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a caregiver profile"}})
		default:
			slog.Error("GetCaregiverChildren failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if children == nil {
		children = []caregiver.ChildSimpleResponse{}
	}
	c.JSON(http.StatusOK, children)
}

package handlers

import (
	"errors"
	"net/http"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/checkin"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CheckinHandler struct {
	service *services.CheckinService
}

func NewCheckinHandler(service *services.CheckinService) *CheckinHandler {
	return &CheckinHandler{service: service}
}

// GenerateToken (Endpoint: 66)
func (h *CheckinHandler) Generate(c *gin.Context) {
	var req checkin.GenerateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	token, expiresAt, err := h.service.GenerateToken(c.Request.Context(), req.ChildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "gagal membuat token"}})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"expiresAt": expiresAt.Format(time.RFC3339),
	})
}

// ValidateToken (Endpoint: 67)
func (h *CheckinHandler) Validate(c *gin.Context) {
	// Retrieve user from JWT context
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate token, record log, and automatically create engagement in CheckinService
	engagementID, err := h.service.ValidateToken(c.Request.Context(), req.Token, requester.ID)

	isNewEngagement := true
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			isNewEngagement = false
			err = nil // Clear error since we consider this a success
		}
	}

	if err != nil {
		status := http.StatusInternalServerError
		code := "INTERNAL_ERROR"
		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
			code = "NOT_FOUND"
		} else if errors.Is(err, repository.ErrForbidden) {
			status = http.StatusForbidden
			code = "FORBIDDEN"
		}

		c.JSON(status, gin.H{"error": gin.H{"code": code, "message": err.Error()}})
		return
	}

	// Success, either new engagement or conflict (active engagement already exists)
	var engagement string
	if engagementID != uuid.Nil {
		engagement = engagementID.String()
	}

	c.JSON(http.StatusOK, gin.H{
		"engagementId":    engagement,
		"isNewEngagement": isNewEngagement,
		"checkedInAt":     time.Now().Format(time.RFC3339),
	})
}

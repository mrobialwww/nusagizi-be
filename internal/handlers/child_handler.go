package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChildHandler struct {
	service *services.ChildService
}

func NewChildHandler(service *services.ChildService) *ChildHandler {
	return &ChildHandler{service: service}
}

// CreateChildProfile (Endpoint: 7)
func (h *ChildHandler) CreateChildProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	var input models.CreateChildInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input name
	if strings.TrimSpace(input.FullName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "full_name cannot be empty"}})
		return
	}

	// Validate input birth_date format
	if _, err := time.Parse(models.DateLayout, input.BirthDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid birth_date format, must be DD-MM-YYYY"}})
		return
	}

	// Validate input gender
	if input.Gender != "male" && input.Gender != "female" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "gender must be 'male' or 'female'"}})
		return
	}

	childID, err := h.service.CreateChild(c.Request.Context(), requester.ID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "no active mother profile for this user"}})
		default:
			slog.Error("CreateChild failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": childID})
}

// UpdateChildProfile (Endpoint: 8)
func (h *ChildHandler) UpdateChildProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	var input models.UpdateChildInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input name
	if input.FullName != nil && strings.TrimSpace(*input.FullName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "full_name cannot be empty"}})
		return
	}

	// Validate input birth_date format
	if input.BirthDate != nil {
		if _, err := time.Parse(models.DateLayout, *input.BirthDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid birth_date format, must be DD-MM-YYYY"}})
			return
		}
	}

	// Validate input gender
	if input.Gender != nil && *input.Gender != "male" && *input.Gender != "female" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "gender must be 'male' or 'female'"}})
		return
	}

	err = h.service.UpdateChild(c.Request.Context(), requester.ID, childID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found or already deleted"}})
		default:
			slog.Error("UpdateChild failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteChildProfile (Endpoint: 9)
func (h *ChildHandler) DeleteChildProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	err = h.service.DeleteChild(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found or already deleted"}})
		default:
			slog.Error("DeleteChild failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetChildProfile (Endpoint: 10)
func (h *ChildHandler) GetChildProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetChildDetail(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found or already deleted"}})
		default:
			slog.Error("GetChildDetail failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetListChild (Endpoint: 11)
func (h *ChildHandler) GetListChild(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	// ID is automatically deduced from the requester's token in the service layer

	children, err := h.service.GetListChild(c.Request.Context(), requester.ID)
	if err != nil {
		slog.Error("GetListChild failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		return
	}

	c.JSON(http.StatusOK, children)
}

// GetChild (Endpoint: 12)
func (h *ChildHandler) GetChild(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetChild(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found or already deleted"}})
		default:
			slog.Error("GetChild failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

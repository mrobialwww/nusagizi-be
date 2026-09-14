package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetUserProfile (Endpoint: 3)
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	resp, err := h.service.GetUser(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you can only access your own profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "user not found"}})
		default:
			slog.Error("GetUser failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUserProfile (Endpoint: 1)
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input name
	if input.FullName != nil && strings.TrimSpace(*input.FullName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "full_name cannot be empty"}})
		return
	}

	// Validate input email
	if input.Email != nil {
		email := strings.TrimSpace(*input.Email)
		if email == "" || !strings.Contains(email, "@") {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid email format"}})
			return
		}
	}

	// Validate input phone number
	if input.PhoneNumber != nil && strings.TrimSpace(*input.PhoneNumber) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "phone_number cannot be empty"}})
		return
	}

	// Validate input gender
	if input.Gender != nil && *input.Gender != "male" && *input.Gender != "female" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "gender must be 'male' or 'female'"}})
		return
	}

	err := h.service.UpdateUser(c.Request.Context(), requester.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you can only update your own profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "user not found"}})
		case errors.Is(err, repository.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT", "message": "email already in use"}})
		default:
			slog.Error("UpdateUser failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteUserProfile (Endpoint: 2)
func (h *UserHandler) DeleteUserProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	err := h.service.SoftDeleteUser(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you can only delete your own account"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "user not found or already deleted"}})
		default:
			slog.Error("DeleteUser failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetMotherProfile (Endpoint: 6)
func (h *UserHandler) GetMotherProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	resp, err := h.service.GetMotherProfile(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "mother profile not found"}})
		default:
			slog.Error("GetMotherProfile failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCaregiverProfile (endpoint: 55)
func (h *UserHandler) GetCaregiverProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	resp, err := h.service.GetCaregiverProfile(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a caregiver profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "caregiver profile not found"}})
		default:
			slog.Error("GetCaregiverProfile failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateMotherProfile (Endpoint: 4)
func (h *UserHandler) CreateMotherProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	motherProfileID, err := h.service.CreateMotherProfile(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT", "message": "mother profile already exists for this user"}})
		default:
			slog.Error("CreateMotherProfile failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": motherProfileID})
}

// CreateCaregiverProfile (Endpoint: 5)
func (h *UserHandler) CreateCaregiverProfile(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	caregiverProfileID, err := h.service.CreateCaregiverProfile(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT", "message": "caregiver profile already exists for this user"}})
		default:
			slog.Error("CreateCaregiverProfile failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": caregiverProfileID})
}

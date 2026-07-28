package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
)

// getUserFromCtx extracts the authenticated *models.User set by the auth middleware.
func getUserFromCtx(c *gin.Context) (*models.User, bool) {
	v, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	user, ok := v.(*models.User)
	return user, ok
}

// errResp builds the standard error JSON body.
func errResp(code, message string) gin.H {
	return gin.H{"error": gin.H{"code": code, "message": message}}
}

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Get handles GET /users/:user_id (endpoint 3).
func (h *UserHandler) Get(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	ctx := c.Request.Context()
	resp, err := h.svc.GetUser(ctx, requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you can only access your own profile"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "user not found"))
		default:
			slog.Error("GetUser failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Update handles PATCH /users/:user_id (endpoint 1).
func (h *UserHandler) Update(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", err.Error()))
		return
	}

	if input.Gender != nil && *input.Gender != "male" && *input.Gender != "female" {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", "gender must be 'male' or 'female'"))
		return
	}

	ctx := c.Request.Context()
	err := h.svc.UpdateUser(ctx, requester.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you can only update your own profile"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "user not found"))
		case errors.Is(err, repository.ErrConflict):
			c.JSON(http.StatusConflict, errResp("CONFLICT", "email already in use"))
		default:
			slog.Error("UpdateUser failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.Status(http.StatusOK)
}

// Delete handles DELETE /users/:user_id (endpoint 2).
func (h *UserHandler) Delete(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	ctx := c.Request.Context()
	err := h.svc.SoftDeleteUser(ctx, requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you can only delete your own account"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "user not found or already deleted"))
		default:
			slog.Error("DeleteUser failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetMotherProfile handles GET /mother-profiles (endpoint 6).
func (h *UserHandler) GetMotherProfile(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	ctx := c.Request.Context()
	resp, err := h.svc.GetMotherProfile(ctx, requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you do not have a mother profile"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "mother profile not found"))
		default:
			slog.Error("GetMotherProfile failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCaregiverProfile handles GET /caregiver-profiles (endpoint 56).
func (h *UserHandler) GetCaregiverProfile(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	ctx := c.Request.Context()
	resp, err := h.svc.GetCaregiverProfile(ctx, requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you do not have a caregiver profile"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "caregiver profile not found"))
		default:
			slog.Error("GetCaregiverProfile failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChildHandler struct {
	svc *services.ChildService
}

func NewChildHandler(svc *services.ChildService) *ChildHandler {
	return &ChildHandler{svc: svc}
}

// Create handles POST /children (endpoint 7).
func (h *ChildHandler) Create(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	var input models.CreateChildInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", err.Error()))
		return
	}

	childID, err := h.svc.CreateChild(c.Request.Context(), requester.ID, &input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "no active mother profile for this user"))
		default:
			slog.Error("CreateChild failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": childID})
}

// Update handles PATCH /children/:child_id (endpoint 8).
func (h *ChildHandler) Update(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", "invalid child_id format"))
		return
	}

	var input models.UpdateChildInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", err.Error()))
		return
	}

	if input.Gender != nil && *input.Gender != "male" && *input.Gender != "female" {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", "gender must be 'male' or 'female'"))
		return
	}

	err = h.svc.UpdateChild(c.Request.Context(), requester.ID, childID, &input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you do not own this child profile"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "child not found or already deleted"))
		default:
			slog.Error("UpdateChild failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.Status(http.StatusOK)
}

// Delete handles DELETE /children/:child_id (endpoint 9).
func (h *ChildHandler) Delete(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", "invalid child_id format"))
		return
	}

	err = h.svc.DeleteChild(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you do not own this child profile"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "child not found or already deleted"))
		default:
			slog.Error("DeleteChild failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.Status(http.StatusOK)
}

// Get handles GET /children/:child_id (endpoint 10).
func (h *ChildHandler) Get(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", "invalid child_id format"))
		return
	}

	resp, err := h.svc.GetChildDetail(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you do not own this child profile"))
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, errResp("NOT_FOUND", "child not found or already deleted"))
		default:
			slog.Error("GetChildDetail failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListByMother handles GET /mother-profiles/:mother_profile_id/children (endpoint 11).
func (h *ChildHandler) ListByMother(c *gin.Context) {
	requester, ok := getUserFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
		return
	}

	// ID is automatically deduced from the requester's token in the service layer

	children, err := h.svc.GetChildrenByMother(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, errResp("FORBIDDEN", "you do not own this mother profile"))
		default:
			slog.Error("GetChildrenByMother failed", "error", err)
			c.JSON(http.StatusInternalServerError, errResp("INTERNAL_ERROR", "internal server error"))
		}
		return
	}

	c.JSON(http.StatusOK, children)
}

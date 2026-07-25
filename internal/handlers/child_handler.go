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
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateChildHandler handles POST /children (endpoint 22).
func CreateChildHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		childID, err := services.CreateChild(pool, requester.ID, &input)
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
}

// UpdateChildHandler handles PATCH /children/:child_id (endpoint 24).
func UpdateChildHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		err = services.UpdateChild(pool, requester.ID, childID, &input)
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
}

// DeleteChildHandler handles DELETE /children/:child_id (endpoint 46).
func DeleteChildHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		err = services.DeleteChild(pool, requester.ID, childID)
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
}

// GetChildHandler handles GET /children/:child_id (endpoint 50).
func GetChildHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		resp, err := services.GetChildDetail(pool, requester.ID, childID)
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
}

// GetChildrenHandler handles GET /mother-profiles/:mother_profile_id/children (endpoint 51).
func GetChildrenHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester, ok := getUserFromCtx(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
			return
		}

		motherProfileID, err := uuid.Parse(c.Param("mother_profile_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, errResp("BAD_REQUEST", "invalid mother_profile_id format"))
			return
		}

		children, err := services.GetChildrenByMother(pool, requester.ID, motherProfileID)
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
}

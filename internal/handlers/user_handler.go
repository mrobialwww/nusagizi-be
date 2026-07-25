package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"nusagizi_be/internal/models"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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

// GetUserHandler handles GET /users/:user_id (endpoint 45).
func GetUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester, ok := getUserFromCtx(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
			return
		}

		resp, err := services.GetUser(pool, requester.ID, c.Param("user_id"))
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
}

// UpdateUserHandler handles PATCH /users/:user_id (endpoint 23).
func UpdateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		err := services.UpdateUser(pool, requester.ID, c.Param("user_id"), input)
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
}

// DeleteUserHandler handles DELETE /users/:user_id (endpoint 25).
func DeleteUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester, ok := getUserFromCtx(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, errResp("UNAUTHORIZED", "user not found in context"))
			return
		}

		err := services.SoftDeleteUser(pool, requester.ID, c.Param("user_id"))
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
}

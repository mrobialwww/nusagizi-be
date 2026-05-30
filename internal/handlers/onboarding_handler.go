package handlers

import (
	"log/slog"
	"net/http"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/config"
	"nusagizi_be/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OnboardingRequest struct {
	Role string `json:"role" binding:"required"`
}

func OnboardingHandler(pool *pgxpool.Pool, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil payload dari request
		var req OnboardingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
			return
		}

		// 2. Ambil user_id (Auth0 Sub) dari context yang di-set middleware
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: user_id not found in context"})
			return
		}
		auth0ID, ok := userIDInterface.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: invalid user_id format"})
			return
		}

		// 3. Update database: simpan role
		err := repository.UpdateUserOnboarding(pool, auth0ID, req.Role)
		if err != nil {
			slog.Error("Failed to update user onboarding", "auth0_id", auth0ID, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update onboarding data"})
			return
		}

		// 4. Step B: Tentukan role ID berdasarkan pilihan user
		var roleID string
		switch req.Role {
		case "mother":
			roleID = cfg.Auth0RoleIDMother
		case "caregiver":
			roleID = cfg.Auth0RoleIDCaregiver
		case "doctor":
			roleID = cfg.Auth0RoleIDDoctor
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role: must be mother, caregiver, or doctor"})
			return
		}

		// 5. Step C: Assign role ke user di Auth0
		err = auth.AssignRoleToUser(
			cfg.Auth0Domain,
			cfg.Auth0ClientID,
			cfg.Auth0ClientSecret,
			auth0ID,
			roleID,
		)
		if err != nil {
			slog.Error("Failed to assign role in Auth0", "auth0_id", auth0ID, "role", req.Role, "error", err)

			// Rollback: kembalikan DB ke kondisi semula agar tidak inkonsisten
			if rbErr := repository.RollbackUserOnboarding(pool, auth0ID); rbErr != nil {
				slog.Error("CRITICAL: Rollback also failed", "auth0_id", auth0ID, "error", rbErr)
			}

			c.JSON(http.StatusBadGateway, gin.H{"error": "Gagal mendaftarkan role, silakan coba lagi"})
			return
		}

		// 6. Berhasil — Flutter boleh refresh JWT sekarang
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Data onboarding berhasil disimpan! Silakan refresh token.",
		})
	}
}

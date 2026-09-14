package handlers

import (
	"log/slog"
	"net/http"
	"nusagizi_be/internal/auth"
	"nusagizi_be/internal/config"
	"nusagizi_be/internal/repository"

	"github.com/gin-gonic/gin"
)

type OnboardingRequest struct {
	Role string `json:"role" binding:"required"`
}

type OnboardingHandler struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewOnboardingHandler(userRepo *repository.UserRepository, cfg *config.Config) *OnboardingHandler {
	return &OnboardingHandler{userRepo: userRepo, cfg: cfg}
}

// Handle handles POST /onboarding.
func (h *OnboardingHandler) Handle(c *gin.Context) {
	// 1. Ambil payload dari request
	var req OnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// 2. Ambil auth0_id (Auth0 Sub) dari context yang di-set middleware
	auth0IDInterface, exists := c.Get("auth0_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: auth0_id not found in context"})
		return
	}
	auth0ID, ok := auth0IDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error: invalid auth0_id format"})
		return
	}

	// 3. Update database: simpan role (cari user berdasarkan auth0_id)
	err := h.userRepo.UpdateOnboarding(c.Request.Context(), auth0ID, req.Role)
	if err != nil {
		slog.Error("Failed to update user onboarding", "auth0_id", auth0ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update onboarding data"})
		return
	}

	// 4. Step B: Tentukan role ID berdasarkan pilihan user
	var roleID string
	switch req.Role {
	case "mother":
		roleID = h.cfg.Auth0RoleIDMother
	case "caregiver":
		roleID = h.cfg.Auth0RoleIDCaregiver
	case "doctor":
		roleID = h.cfg.Auth0RoleIDDoctor
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role: must be mother, caregiver, or doctor"})
		return
	}

	// 5. Step C: Assign role ke user di Auth0
	err = auth.AssignRoleToUser(h.cfg.Auth0Domain, h.cfg.M2MClientID, h.cfg.M2MClientSecret, auth0ID, roleID)
	if err != nil {
		slog.Error("Failed to assign Auth0 role", "auth0_id", auth0ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan role di Auth0"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data onboarding berhasil disimpan! Silakan refresh token.",
	})
}

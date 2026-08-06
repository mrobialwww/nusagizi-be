package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nusagizi_be/internal/services"
)

// ConfirmEmailHandler handles POST /internal/confirm-email.
//
// This is an unprotected endpoint used exclusively during registration to 
// translate "email ownership proof" (ID Token from Passwordless OTP) into 
// email_verified=true on the user's database account via Management API.
type ConfirmEmailHandler struct {
	service *services.EmailOTPService
}

func NewConfirmEmailHandler(service *services.EmailOTPService) *ConfirmEmailHandler {
	return &ConfirmEmailHandler{service: service}
}

type confirmEmailRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func (h *ConfirmEmailHandler) Handle(c *gin.Context) {
	var req confirmEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token wajib diisi"})
		return
	}

	ctx := c.Request.Context()
	err := h.service.ConfirmEmail(ctx, req.IDToken)
	if err != nil {
		switch err.Error() {
		case "id_token tidak valid", "email belum terverifikasi oleh Auth0":
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		case "id_token tidak memiliki klaim email":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "akun database untuk email ini tidak ditemukan":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "verified"})
}

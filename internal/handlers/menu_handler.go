package handlers

import (
	"net/http"
	"strings"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuService *services.MenuService
}

func NewMenuHandler(menuService *services.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

// GenerateMenu (Endpoint: 39)
func (h *MenuHandler) GenerateMenu(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user not found in context"})
		return
	}

	err := h.menuService.GenerateMenu(c.Request.Context(), requester.ID)
	if err == nil {
		c.JSON(http.StatusCreated, gin.H{
			"status":  "created",
			"message": "Menu berhasil di-generate",
		})
		return
	}

	errStr := err.Error()
	status := http.StatusInternalServerError

	if strings.Contains(errStr, "mother profile not found") || strings.Contains(errStr, "no children found for this mother") {
		status = http.StatusNotFound
	} else if strings.Contains(errStr, "ai engine is currently unavailable") ||
		strings.Contains(errStr, "ai engine returned an invalid response") ||
		strings.Contains(errStr, "ai engine returned no menu data") {
		status = http.StatusBadGateway
	}

	c.JSON(status, gin.H{"error": errStr})
}

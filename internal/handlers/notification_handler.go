package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/notification"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// GetNotifications (endpoint: 63)
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	notifications, err := h.service.GetNotifications(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "no notifications found"}})
		default:
			slog.Error("GetNotifications failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if notifications == nil {
		notifications = []notification.NotificationResponse{}
	}
	c.JSON(http.StatusOK, notifications)
}

// GetLatestNotification (Endpoint: 64)
func (h *NotificationHandler) GetLatestNotification(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	notificationRes, err := h.service.GetLatestNotification(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "no notifications found"}})
		default:
			slog.Error("GetLatestNotification failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, notificationRes)
}

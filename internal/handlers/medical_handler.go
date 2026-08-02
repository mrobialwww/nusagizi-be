package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/medical"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MedicalHandler struct {
	service *services.MedicalService
}

func NewMedicalHandler(service *services.MedicalService) *MedicalHandler {
	return &MedicalHandler{service: service}
}

// GetMedicalNotes (Endpoint: 60, 61)
func (h *MedicalHandler) GetMedicalNotes(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	// Get status and month from query params
	status := c.Query("status")
	monthStr := c.Query("month")

	// Validate input status
	switch status {
	case "active", "history":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "status must be active or history"}})
		return
	}

	// Validate input month
	var month *int
	if monthStr != "" {
		m, err := strconv.Atoi(monthStr)
		if err != nil || m < 1 || m > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "month must be between 1 and 12"}})
			return
		}
		month = &m
	}

	notes, err := h.service.GetMedicalNotes(c.Request.Context(), requester.ID, status, month)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		default:
			slog.Error("GetMedicalNotes failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if notes == nil {
		notes = []medical.MedicalNoteResponse{}
	}
	c.JSON(http.StatusOK, notes)
}

// GetMedicalNoteDetail (Endpoint: 62)
func (h *MedicalHandler) GetMedicalNoteDetail(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	noteID, err := uuid.Parse(c.Param("medical_note_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid medical_note_id format"}})
		return
	}

	note, err := h.service.GetMedicalNoteDetail(c.Request.Context(), requester.ID, noteID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have permission to view this medical note"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "medical note not found"}})
		default:
			slog.Error("GetMedicalNoteDetail failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, note)
}

// CreateMedicalNote (Endpoint: 63)
func (h *MedicalHandler) CreateMedicalNote(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	var input medical.CreateMedicalNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input child_id
	if input.ChildID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "child_id cannot be empty"}})
		return
	}

	// Validate input doctor_name
	if strings.TrimSpace(input.DoctorName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "doctor_name cannot be empty"}})
		return
	}

	// Validate input recommendation
	if strings.TrimSpace(input.Recommendation) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "recommendation cannot be empty"}})
		return
	}

	// Validate input valid_date
	if strings.TrimSpace(input.ValidDate) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "valid_date cannot be empty"}})
		return
	}
	if _, err := time.Parse(models.DateLayout, input.ValidDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid valid_date format, must be DD-MM-YYYY"}})
		return
	}

	noteID, err := h.service.CreateMedicalNote(c.Request.Context(), requester.ID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found"}})
		default:
			slog.Error("CreateMedicalNote failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": noteID})
}

// UpdateMedicalNote (Endpoint: 64)
func (h *MedicalHandler) UpdateMedicalNote(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	noteID, err := uuid.Parse(c.Param("medical_note_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid medical_note_id format"}})
		return
	}

	var input medical.UpdateMedicalNoteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input doctor_name
	if input.DoctorName != nil && strings.TrimSpace(*input.DoctorName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "doctor_name cannot be empty"}})
		return
	}

	// Validate input recommendation
	if input.Recommendation != nil && strings.TrimSpace(*input.Recommendation) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "recommendation cannot be empty"}})
		return
	}

	// Validate input valid_date
	if input.ValidDate != nil {
		if strings.TrimSpace(*input.ValidDate) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "valid_date cannot be empty"}})
			return
		}
		if _, err := time.Parse(models.DateLayout, *input.ValidDate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid valid_date format, must be DD-MM-YYYY"}})
			return
		}
	}

	err = h.service.UpdateMedicalNote(c.Request.Context(), requester.ID, noteID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have permission to edit this note"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "medical note not found"}})
		default:
			slog.Error("UpdateMedicalNote failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteMedicalNote (Endpoint: 65)
func (h *MedicalHandler) DeleteMedicalNote(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	noteID, err := uuid.Parse(c.Param("medical_note_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid medical_note_id format"}})
		return
	}

	err = h.service.DeleteMedicalNote(c.Request.Context(), requester.ID, noteID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have permission to delete this note"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "medical note not found"}})
		default:
			slog.Error("DeleteMedicalNote failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

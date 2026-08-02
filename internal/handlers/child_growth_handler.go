package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"nusagizi_be/internal/models"
	child_growth "nusagizi_be/internal/models/child_growth"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChildGrowthHandler struct {
	service *services.ChildGrowthService
}

func NewChildGrowthHandler(service *services.ChildGrowthService) *ChildGrowthHandler {
	return &ChildGrowthHandler{service: service}
}

// GetLatestGrowthReport (Endpoint: 13)
func (h *ChildGrowthHandler) GetLatestGrowthReport(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetLatestGrowthReport(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found or no growth report available"}})
		default:
			slog.Error("GetLatestGrowthReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetGrowthAnalyses (Endpoint: 14)
func (h *ChildGrowthHandler) GetGrowthAnalyses(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	// Get analysis_type and age_range from query params
	analysisType := c.Query("analysis_type")
	ageRange := c.Query("age_range")

	// Validate input analysis_type and age_range
	if analysisType == "" || ageRange == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "analysis_type and age_range are required"}})
		return
	}

	// Validate input analysis_type format
	switch analysisType {
	case "weight_for_age", "height_for_age", "weight_for_height", "bmi_for_age", "head_circumference_for_age":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid analysis_type"}})
		return
	}

	// Validate input age_range format
	switch ageRange {
	case "0-2", "0-6", "0-60":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid age_range, must be one of: 0-2, 0-6, 0-60"}})
		return
	}

	resp, err := h.service.GetGrowthAnalyses(c.Request.Context(), requester.ID, childID, analysisType, ageRange)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found"}})
		default:
			slog.Error("GetGrowthAnalyses failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetGrowthReports (Endpoint: 15)
func (h *ChildGrowthHandler) GetGrowthReports(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	resp, err := h.service.GetGrowthReports(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found"}})
		default:
			slog.Error("GetGrowthReports failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if resp == nil {
		resp = []child_growth.GrowthReportWithStatus{}
	}
	c.JSON(http.StatusOK, resp)
}

// CreateGrowthReport (Endpoint: 16)
func (h *ChildGrowthHandler) CreateGrowthReport(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	var input child_growth.CreateGrowthReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input measured_at
	if strings.TrimSpace(input.MeasuredAt) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "measured_at is required"}})
		return
	}
	if _, err := time.Parse(models.DateLayout, input.MeasuredAt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid measured_at format, must be DD-MM-YYYY"}})
		return
	}

	// Validate input weight_kg
	if input.WeightKg != nil && *input.WeightKg <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "weight_kg must be > 0"}})
		return
	}

	// Validate input height_cm
	if input.HeightCm != nil && *input.HeightCm <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "height_cm must be > 0"}})
		return
	}

	// Validate input head_circumference_cm
	if input.HeadCircumferenceCm != nil && *input.HeadCircumferenceCm <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "head_circumference_cm must be > 0"}})
		return
	}

	reportID, err := h.service.CreateGrowthReport(c.Request.Context(), requester.ID, childID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found"}})
		default:
			slog.Error("CreateGrowthReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": reportID})
}

// UpdateGrowthReport (Endpoint: 17)
func (h *ChildGrowthHandler) UpdateGrowthReport(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	reportID, err := uuid.Parse(c.Param("child_growth_report_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_growth_report_id format"}})
		return
	}

	var input child_growth.UpdateGrowthReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input measured_at
	if input.MeasuredAt != nil {
		if strings.TrimSpace(*input.MeasuredAt) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "measured_at cannot be empty"}})
			return
		}
		if _, err := time.Parse(models.DateLayout, *input.MeasuredAt); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid measured_at format, must be DD-MM-YYYY"}})
			return
		}
	}

	// Validate input weight_kg
	if input.WeightKg != nil && *input.WeightKg <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "weight_kg must be > 0"}})
		return
	}

	// Validate input height_cm
	if input.HeightCm != nil && *input.HeightCm <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "height_cm must be > 0"}})
		return
	}

	// Validate input head_circumference_cm
	if input.HeadCircumferenceCm != nil && *input.HeadCircumferenceCm <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "head_circumference_cm must be > 0"}})
		return
	}

	err = h.service.UpdateGrowthReport(c.Request.Context(), requester.ID, childID, reportID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "report not found"}})
		default:
			slog.Error("UpdateGrowthReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteGrowthReport (Endpoint: 18)
func (h *ChildGrowthHandler) DeleteGrowthReport(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	childID, err := uuid.Parse(c.Param("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	reportID, err := uuid.Parse(c.Param("child_growth_report_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_growth_report_id format"}})
		return
	}

	err = h.service.DeleteGrowthReport(c.Request.Context(), requester.ID, childID, reportID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "report not found"}})
		default:
			slog.Error("DeleteGrowthReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

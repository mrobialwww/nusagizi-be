package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"nusagizi_be/internal/models"
	child_dev "nusagizi_be/internal/models/child_development"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChildDevelopmentHandler struct {
	service *services.ChildDevelopmentService
}

func NewChildDevelopmentHandler(service *services.ChildDevelopmentService) *ChildDevelopmentHandler {
	return &ChildDevelopmentHandler{service: service}
}

// GetLatestDevelopmentReport (Endpoint: 19)
func (h *ChildDevelopmentHandler) GetLatestDevelopmentReport(c *gin.Context) {
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

	resp, err := h.service.GetLatestDevelopmentReport(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child development report not found"}})
		default:
			slog.Error("GetLatestDevelopmentReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetDevelopmentReports (Endpoint: 20)
func (h *ChildDevelopmentHandler) GetDevelopmentReports(c *gin.Context) {
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

	resp, err := h.service.GetDevelopmentReports(c.Request.Context(), requester.ID, childID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child development reports not found"}})
		default:
			slog.Error("GetDevelopmentReports failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}
	if resp == nil {
		resp = []child_dev.DevelopmentReportResponse{}
	}
	c.JSON(http.StatusOK, resp)
}

// GetDevelopmentReportByID (Endpoint: 21)
func (h *ChildDevelopmentHandler) GetDevelopmentReportByID(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}
	reportID, err := uuid.Parse(c.Param("child_development_report_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid report_id format"}})
		return
	}

	resp, err := h.service.GetDevelopmentReportByID(c.Request.Context(), requester.ID, reportID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child development report not found"}})
		default:
			slog.Error("GetDevelopmentReportByID failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetKPSPQuestions (Endpoint: 22)
func (h *ChildDevelopmentHandler) GetKPSPQuestions(c *gin.Context) {

	// Get month_target from query params
	monthTargetStr := c.Query("month_target")
	monthTarget, err := strconv.Atoi(monthTargetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "month_target is required and must be an integer"}})
		return
	}

	// Validate input month_target
	switch monthTarget {
	case 3, 6, 9, 12, 15, 18, 21, 24, 30, 36, 42, 48, 54, 60:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid month_target"}})
		return
	}

	resp, err := h.service.GetKPSPQuestions(c.Request.Context(), monthTarget)
	if err != nil {
		slog.Error("GetKPSPQuestions failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		return
	}

	if resp == nil {
		resp = []child_dev.AssessmentKPSPQuestion{}
	}
	c.JSON(http.StatusOK, resp)
}

// CreateDevelopmentReport (Endpoint: 23)
func (h *ChildDevelopmentHandler) CreateDevelopmentReport(c *gin.Context) {
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

	var input child_dev.CreateDevelopmentReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input month_target
	switch input.MonthTarget {
	case 3, 6, 9, 12, 15, 18, 21, 24, 30, 36, 42, 48, 54, 60:
	default:
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "month_target mismatch"}})
		return
	}

	// Validate input list_answer
	if len(input.ListAnswer) != 10 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "answers must be exactly 10"}})
		return
	}

	reportID, err := h.service.CreateDevelopmentReport(c.Request.Context(), requester.ID, childID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		default:
			slog.Error("CreateDevelopmentReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": reportID})
}

// UpdateDevelopmentReport (Endpoint: 24)
func (h *ChildDevelopmentHandler) UpdateDevelopmentReport(c *gin.Context) {
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

	reportID, err := uuid.Parse(c.Param("child_development_report_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid report_id format"}})
		return
	}

	var input child_dev.UpdateDevelopmentReportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input list_answer
	if len(input.ListAnswer) == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "UNPROCESSABLE_ENTITY", "message": "answers cannot be empty"}})
		return
	}

	updatedID, err := h.service.UpdateDevelopmentReport(c.Request.Context(), requester.ID, childID, reportID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child development report not found"}})
		default:
			slog.Error("UpdateDevelopmentReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": updatedID})
}

// GetChecklistMilestoneTasks (Endpoint: 25)
func (h *ChildDevelopmentHandler) GetChecklistMilestoneTasks(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	// Get child_id from query params
	childID, err := uuid.Parse(c.Query("child_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
		return
	}

	// Get month_target from query params
	monthTargetStr := c.Query("month_target")
	monthTarget, err := strconv.Atoi(monthTargetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "month_target is required and must be an integer"}})
		return
	}

	// Validate input month_target
	switch monthTarget {
	case 3, 6, 9, 12, 15, 18, 21, 24, 30, 36, 42, 48, 54, 60:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid month_target"}})
		return
	}

	resp, err := h.service.GetChecklistMilestoneTasks(c.Request.Context(), requester.ID, childID, monthTarget)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		default:
			slog.Error("GetChecklistMilestoneTasks failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetRecommendations (Endpoint: 26)
func (h *ChildDevelopmentHandler) GetRecommendations(c *gin.Context) {
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

	reportID, err := uuid.Parse(c.Param("report_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid report_id format"}})
		return
	}

	items, err := h.service.GetRecommendations(c.Request.Context(), requester.ID, childID, reportID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "development report not found"}})
		default:
			slog.Error("GetRecommendations failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, items)
}

// UpdateChecklistMilestone (Endpoint: 27)
func (h *ChildDevelopmentHandler) UpdateChecklistMilestone(c *gin.Context) {
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

	// Empty array is explicitly allowed to enable clearing all tasks
	var input struct {
		AssessmentKPSPQuestionIDs []uuid.UUID `json:"assessment_kpsp_question_ids"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	err = h.service.UpdateChecklistMilestone(c.Request.Context(), requester.ID, childID, input.AssessmentKPSPQuestionIDs)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "some task IDs were not found"}})
		default:
			slog.Error("UpdateChecklistMilestone failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeleteDevelopmentReport (Endpoint: 28)
func (h *ChildDevelopmentHandler) DeleteDevelopmentReport(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}
	reportID, err := uuid.Parse(c.Param("child_development_report_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_development_report_id format"}})
		return
	}

	err = h.service.DeleteDevelopmentReport(c.Request.Context(), requester.ID, reportID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child development report not found"}})
		default:
			slog.Error("DeleteDevelopmentReport failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

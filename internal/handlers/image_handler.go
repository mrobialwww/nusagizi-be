package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/services"
)

type ImageHandler struct {
	service *services.ImageService
}

func NewImageHandler(service *services.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

func (h *ImageHandler) PresignUpload(c *gin.Context) {
	var req models.PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate content type
	if req.ContentType != "image/webp" && req.ContentType != "image/jpeg" && req.ContentType != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type, must be image/webp, image/jpeg, or image/png"})
		return
	}

	// Validate category
	if req.Category != "profile" && req.Category != "child-profile" && req.Category != "social" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category. Allowed values: profile, child-profile, social"})
		return
	}

	// Retrieve logged-in user from context
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user not found in context"})
		return
	}

	uploadURL, objectKey, err := h.service.GeneratePresignPutURL(c.Request.Context(), req.Category, req.OwnerID, requester.ID, req.ChildID, req.ContentType)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"object_key": objectKey,
	})
}

func (h *ImageHandler) Confirm(c *gin.Context) {
	var req models.ConfirmUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// We skip verifying the existence of the object right now to make this lightweight.
	c.JSON(http.StatusOK, gin.H{"message": "Upload confirmed"})
}

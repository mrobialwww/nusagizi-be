package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"nusagizi_be/internal/models"
	"nusagizi_be/internal/models/photos_contacts"
	"nusagizi_be/internal/repository"
	"nusagizi_be/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhotosContactsHandler struct {
	service *services.PhotosContactsService
}

func NewPhotosContactsHandler(service *services.PhotosContactsService) *PhotosContactsHandler {
	return &PhotosContactsHandler{service: service}
}

// GetContacts (Endpoint: 41)
func (h *PhotosContactsHandler) GetContacts(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	contacts, err := h.service.GetContacts(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		default:
			slog.Error("GetContacts failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// AddContact (endpoint: 49)
func (h *PhotosContactsHandler) AddContact(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	var input struct {
		RelatedMotherProfileID uuid.UUID `json:"related_mother_profile_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input related_mother_profile_id
	if input.RelatedMotherProfileID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "related_mother_profile_id cannot be empty"}})
		return
	}

	contactID, err := h.service.AddContact(c.Request.Context(), requester.ID, input.RelatedMotherProfileID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "related mother profile not found"}})
		case errors.Is(err, repository.ErrConflict):
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "CONFLICT", "message": "contact already exists"}})
		case err.Error() == "cannot add yourself as a contact":
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		default:
			slog.Error("AddContact failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": contactID})
}

// DeleteContact (Endpoint: 42)
func (h *PhotosContactsHandler) DeleteContact(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	contactID, err := uuid.Parse(c.Param("contact_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid contact_id format"}})
		return
	}

	err = h.service.DeleteContact(c.Request.Context(), requester.ID, contactID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this contact"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "contact not found"}})
		default:
			slog.Error("DeleteContact failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// GetMotherChildPhotos (Endpoint: 43)
func (h *PhotosContactsHandler) GetMotherChildPhotos(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	var childIDPtr *uuid.UUID
	if childIDStr := c.Query("child_id"); childIDStr != "" {
		childID, err := uuid.Parse(childIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_id format"}})
			return
		}
		childIDPtr = &childID
	}

	latestPerChild := c.Query("latest_per_child") == "true"

	photos, err := h.service.GetMotherChildPhotos(c.Request.Context(), requester.ID, childIDPtr, latestPerChild)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		default:
			slog.Error("GetMotherChildPhotos failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if photos == nil {
		photos = []photos_contacts.ChildPhotoResponse{}
	}
	c.JSON(http.StatusOK, photos)
}

// GetContactChildPhotos (Endpoint: 44)
func (h *PhotosContactsHandler) GetContactChildPhotos(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	contactID, err := uuid.Parse(c.Param("contact_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid contact_id format"}})
		return
	}

	photos, err := h.service.GetContactChildPhotos(c.Request.Context(), requester.ID, contactID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "contact is not yours or not found"}})
		default:
			slog.Error("GetContactChildPhotos failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if photos == nil {
		photos = []photos_contacts.ChildPhotoResponse{}
	}
	c.JSON(http.StatusOK, photos)
}

// GetAllChildPhotos (Endpoint: 45)
func (h *PhotosContactsHandler) GetAllChildPhotos(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	photos, err := h.service.GetAllChildPhotos(c.Request.Context(), requester.ID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have a mother profile"}})
		default:
			slog.Error("GetAllChildPhotos failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	if photos == nil {
		photos = []photos_contacts.ChildPhotoResponse{}
	}
	c.JSON(http.StatusOK, photos)
}

// GetPhotoDetail (Endpoint: 46)
func (h *PhotosContactsHandler) GetPhotoDetail(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	photoID, err := uuid.Parse(c.Param("child_photo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_photo_id format"}})
		return
	}

	photo, err := h.service.GetPhotoDetail(c.Request.Context(), requester.ID, photoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have access to this photo"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "photo not found"}})
		default:
			slog.Error("GetPhotoDetail failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusOK, photo)
}

// AddPhotoMother (Endpoint: 47)
func (h *PhotosContactsHandler) AddPhotoMother(c *gin.Context) {
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

	var input photos_contacts.CreatePhotoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input url
	if strings.TrimSpace(input.URL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "url cannot be empty"}})
		return
	}

	// Validate input caption
	if strings.TrimSpace(input.Caption) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "caption cannot be empty"}})
		return
	}

	// Validate input visibility
	switch input.Visibility {
	case "all", "private", "selected_only":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "visibility must be all, private, or selected_only"}})
		return
	}

	// Validate input list_visibility
	if input.Visibility == "selected_only" && len(input.ListVisibility) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "list_visibility cannot be empty if visibility is selected_only"}})
		return
	}

	photoID, err := h.service.AddPhotoMother(c.Request.Context(), requester.ID, childID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this child profile"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found"}})
		default:
			slog.Error("AddPhotoMother failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": photoID})
}

// AddPhotoCaregiver (Endpoint: 48)
func (h *PhotosContactsHandler) AddPhotoCaregiver(c *gin.Context) {
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

	var input photos_contacts.CreatePhotoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input url
	if strings.TrimSpace(input.URL) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "url cannot be empty"}})
		return
	}

	// Validate input is_review_required
	if !input.IsReviewRequired {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "is_review_required must be true for caregiver uploads"}})
		return
	}

	photoID, err := h.service.AddPhotoCaregiver(c.Request.Context(), requester.ID, childID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have an active engagement with this child"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "child not found"}})
		default:
			slog.Error("AddPhotoCaregiver failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": photoID})
}

// UpdatePhoto (endpoint: 50)
func (h *PhotosContactsHandler) UpdatePhoto(c *gin.Context) {
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

	photoID, err := uuid.Parse(c.Param("child_photo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_photo_id format"}})
		return
	}

	var input photos_contacts.UpdatePhotoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": err.Error()}})
		return
	}

	// Validate input visibility
	if input.Visibility != nil {
		switch *input.Visibility {
		case "all", "private", "selected_only":
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "visibility must be all, private, or selected_only"}})
			return
		}

		// Validate input list_visibility
		if *input.Visibility == "selected_only" {
			if input.ListVisibility == nil || len(*input.ListVisibility) == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "list_visibility cannot be empty if visibility is selected_only"}})
				return
			}
		}
	}

	err = h.service.UpdatePhoto(c.Request.Context(), requester.ID, childID, photoID, &input)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not have permission to edit this photo"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "photo not found"}})
		default:
			slog.Error("UpdatePhoto failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

// DeletePhoto (endpoint: 51)
func (h *PhotosContactsHandler) DeletePhoto(c *gin.Context) {
	v, exists := c.Get("user")
	requester, ok := v.(*models.User)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "user not found in context"}})
		return
	}

	photoID, err := uuid.Parse(c.Param("child_photo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid child_photo_id format"}})
		return
	}

	err = h.service.DeletePhoto(c.Request.Context(), requester.ID, photoID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "you do not own this photo"}})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "photo not found"}})
		default:
			slog.Error("DeletePhoto failed", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL_ERROR", "message": "internal server error"}})
		}
		return
	}

	c.Status(http.StatusOK)
}

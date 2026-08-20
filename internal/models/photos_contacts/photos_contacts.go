package photos_contacts

import (
	"time"

	"github.com/google/uuid"
)

// ContactResponse represents a contact (friend) with their user details.
type ContactResponse struct {
	ContactID              uuid.UUID `json:"contact_id"`
	RelatedMotherProfileID uuid.UUID `json:"related_mother_profile_id"`
	FullName               string    `json:"full_name"`
	PhotoURL               *string   `json:"photo_url"`
}

// ChildPhotoResponse represents a photo of a child.
type ChildPhotoResponse struct {
	ID               uuid.UUID `json:"id"`
	ChildID          uuid.UUID `json:"child_id"`
	PhotoURL         string    `json:"photo_url"`
	Caption          string    `json:"caption,omitempty"`
	Visibility       string    `json:"visibility,omitempty"`
	IsReviewRequired bool      `json:"is_review_required"`
	CreatedAt        time.Time `json:"created_at"`
}

// CreatePhotoInput is the payload for adding a new photo.
type CreatePhotoInput struct {
	URL              string      `json:"url"`
	Caption          string      `json:"caption"`
	Visibility       string      `json:"visibility"`
	ListVisibility   []uuid.UUID `json:"list_visibility"`
	IsReviewRequired bool        `json:"is_review_required"`
}

// UpdatePhotoInput is the payload for editing an existing photo.
type UpdatePhotoInput struct {
	Caption          *string      `json:"caption"`
	Visibility       *string      `json:"visibility"`
	ListVisibility   *[]uuid.UUID `json:"list_visibility"`
	IsReviewRequired *bool        `json:"is_review_required"`
}

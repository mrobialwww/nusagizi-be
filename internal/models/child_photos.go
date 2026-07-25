package models

import (
	"time"

	"github.com/google/uuid"
)

// PhotoVisibilityEnum merepresentasikan tipe ENUM untuk visibilitas foto
type PhotoVisibilityEnum string

const (
	PhotoVisibilityAll     PhotoVisibilityEnum = "ALL"
	PhotoVisibilityPrivate PhotoVisibilityEnum = "PRIVATE"
	PhotoVisibilityOnly    PhotoVisibilityEnum = "ONLY"
)

// IsValid memastikan nilai PhotoVisibilityEnum valid sesuai dengan database
func (p PhotoVisibilityEnum) IsValid() bool {
	switch p {
	case PhotoVisibilityAll, PhotoVisibilityPrivate, PhotoVisibilityOnly:
		return true
	}
	return false
}

type ChildPhoto struct {
	ID               uuid.UUID           `json:"child_photo_id" db:"child_photo_id"`
	ChildID          uuid.UUID           `json:"child_id" db:"child_id"`
	URL              string              `json:"url" db:"url"`
	Captions         string              `json:"captions" db:"captions"`
	Visibility       PhotoVisibilityEnum `json:"visibility" db:"visibility"`
	IsReviewRequired bool                `json:"is_review_required" db:"is_review_required"`
	CreatedAt        time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at" db:"updated_at"`
}

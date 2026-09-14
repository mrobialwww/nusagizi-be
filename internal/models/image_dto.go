package models

type PresignUploadRequest struct {
	Category    string  `json:"category" binding:"required"`
	OwnerID     string  `json:"owner_id"` // can be empty for child creation
	ChildID     *string `json:"child_id"` // optional child_id cross-reference
	ContentType string  `json:"content_type" binding:"required"`
}

type ConfirmUploadRequest struct {
	Category  string `json:"category" binding:"required"`
	OwnerID   string `json:"owner_id"`
	ObjectKey string `json:"object_key" binding:"required"`
	// Additional metadata if passed
}

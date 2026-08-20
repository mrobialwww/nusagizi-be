package utils

import (
	"strings"
)

// // ResolvePhotoURL converts an object_key stored in DB to a presigned GET URL.
// // If the value looks like a full URL (starts with "http"), it's returned as-is.
// func ResolvePhotoURL(ctx context.Context, key *string, storage storage.ObjectStorage) *string {
// 	if key == nil || *key == "" {
// 		return nil
// 	}
// 	if strings.HasPrefix(*key, "http") {
// 		return key // legacy URL or external, return as-is
// 	}
// 	url, err := storage.PresignGetURL(ctx, *key, 15*time.Minute)
// 	if err != nil {
// 		return key // fallback: return key as-is on error
// 	}
// 	return &url
// }

// ResolvePublicPhotoURL appends the R2 public base URL to the object key if it is not a full URL.
func ResolvePublicPhotoURL(baseURL string, key *string) *string {
	if key == nil || *key == "" {
		return nil
	}
	if strings.HasPrefix(*key, "http") {
		return key // legacy URL or external, return as-is
	}

	cleanBaseURL := strings.TrimSuffix(baseURL, "/")
	cleanKey := strings.TrimPrefix(*key, "/")

	if cleanBaseURL == "" {
		return key
	}

	url := cleanBaseURL + "/" + cleanKey
	return &url
}

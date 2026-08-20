package storage

import (
	"context"
	"time"
)

type ObjectStorage interface {
	PresignPutURL(ctx context.Context, key, contentType string, expires time.Duration) (string, error)
	PresignGetURL(ctx context.Context, key string, expires time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}

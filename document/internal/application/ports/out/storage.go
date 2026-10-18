package out

import (
	"context"
	"io"

	"github.com/aakashloyar/elevate/document/internal/domain"
)

type Storage interface {
	CreateMultipartUpload(ctx context.Context, key string) (string, error)
	GeneratePresignedPartURL(ctx context.Context, key, uploadID string, partNumber int32) (string, error)
	CompleteMultipartUpload(ctx context.Context, key, uploadID string, parts []domain.Chunk) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	BuildStorageKey(uploadID, fileName string) string
}

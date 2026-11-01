package upload

import (
	"context"
	"time"

	"github.com/aakashloyar/elevate/document/internal/domain"
)

type ListDocumentsInput struct {
	CreatedBy string
	Status    *domain.UploadStatus
	Search    string
	Limit     int
	Offset    int
}

type DocumentSummary struct {
	ID          string              `json:"id"`
	CreatedBy   string              `json:"created_by"`
	FileName    string              `json:"file_name"`
	ContentType string              `json:"content_type"`
	FileSize    int64               `json:"file_size"`
	Status      domain.UploadStatus `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type ListDocumentsOutput struct {
	Documents []DocumentSummary `json:"documents"`
	Total     int               `json:"total"`
	Limit     int               `json:"limit"`
	Offset    int               `json:"offset"`
}

type ListDocumentsService interface {
	Execute(context.Context, ListDocumentsInput) (ListDocumentsOutput, error)
}

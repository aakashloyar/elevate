package in

import (
	"context"

	"github.com/aakashloyar/elevate/document/internal/domain"
)

type ProcessDocumentService interface {
	Execute(ctx context.Context, event domain.DocumentUploadedEvent) error
}

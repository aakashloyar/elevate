package out

import (
	"context"

	"github.com/aakashloyar/elevate/document/internal/domain"
)

type DocumentUploadedMessage struct {
	Event domain.DocumentUploadedEvent
}

type DocumentUploadedPublisher interface {
	PublishDocumentUploaded(ctx context.Context, message DocumentUploadedMessage) error
	Close()
}

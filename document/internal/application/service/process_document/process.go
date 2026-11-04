package process_document

import (
	"context"
	"fmt"
	"io"
	"strings"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type Service struct {
	uploads out.UploadRespository
	storage out.Storage
}

func NewService(uploads out.UploadRespository, storage out.Storage) in.ProcessDocumentService {
	return &Service{uploads: uploads, storage: storage}
}

func (s *Service) Execute(ctx context.Context, event domain.DocumentUploadedEvent) error {
	if strings.TrimSpace(event.DocumentID) == "" {
		return fmt.Errorf("document id is required")
	}
	upload, err := s.uploads.FindByID(event.DocumentID)
	if err != nil {
		return fmt.Errorf("find uploaded document %s: %w", event.DocumentID, err)
	}
	if strings.TrimSpace(upload.StorageKey) == "" {
		return fmt.Errorf("uploaded document %s has no storage key", event.DocumentID)
	}
	if err := s.uploads.Complete(event.DocumentID, domain.StatusProcessing); err != nil {
		return fmt.Errorf("mark document %s processing: %w", event.DocumentID, err)
	}
	body, err := s.storage.Download(ctx, upload.StorageKey)
	if err != nil {
		_ = s.uploads.Complete(event.DocumentID, domain.StatusFailed)
		return fmt.Errorf("download document %s: %w", event.DocumentID, err)
	}
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		_ = s.uploads.Complete(event.DocumentID, domain.StatusFailed)
		return fmt.Errorf("read document %s: %w", event.DocumentID, err)
	}
	text, err := extractText(data, upload.FileName, upload.ContentType)
	if err != nil {
		_ = s.uploads.Complete(event.DocumentID, domain.StatusFailed)
		return fmt.Errorf("extract text from document %s: %w", event.DocumentID, err)
	}
	if strings.TrimSpace(text) == "" {
		_ = s.uploads.Complete(event.DocumentID, domain.StatusFailed)
		return fmt.Errorf("document %s has no extracted text", event.DocumentID)
	}
	if err := s.uploads.SaveExtractedText(event.DocumentID, text); err != nil {
		_ = s.uploads.Complete(event.DocumentID, domain.StatusFailed)
		return fmt.Errorf("save extracted text for document %s: %w", event.DocumentID, err)
	}
	if err := s.uploads.Complete(event.DocumentID, domain.StatusReady); err != nil {
		return fmt.Errorf("mark document %s ready: %w", event.DocumentID, err)
	}
	return nil
}

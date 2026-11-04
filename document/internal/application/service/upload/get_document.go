package service

import (
	"context"
	"fmt"
	"strings"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
)

type GetDocumentService struct {
	repo out.UploadRespository
}

func NewGetDocumentService(repo out.UploadRespository) in.GetDocumentService {
	return &GetDocumentService{repo: repo}
}

func (s *GetDocumentService) Execute(ctx context.Context, input in.GetDocumentInput) (in.GetDocumentOutput, error) {
	_ = ctx
	documentID := strings.TrimSpace(input.DocumentID)
	if documentID == "" {
		return in.GetDocumentOutput{}, fmt.Errorf("document id is required")
	}
	document, err := s.repo.FindByID(documentID)
	if err != nil {
		return in.GetDocumentOutput{}, err
	}
	return in.GetDocumentOutput{Document: in.DocumentSummary{
		ID: document.ID, CreatedBy: document.CreatedBy, FileName: document.FileName,
		ContentType: document.ContentType, FileSize: document.FileSize, Status: document.Status,
		CreatedAt: document.CreatedAt, UpdatedAt: document.UpdatedAt,
	}}, nil
}

package service

import (
	"context"
	"fmt"
	"strings"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
)

type GetDocumentContentService struct {
	repo out.UploadRespository
}

func NewGetDocumentContentService(repo out.UploadRespository) in.GetDocumentContentService {
	return &GetDocumentContentService{repo: repo}
}

func (s *GetDocumentContentService) Execute(ctx context.Context, input in.GetDocumentContentInput) (in.GetDocumentContentOutput, error) {
	_ = ctx
	documentID := strings.TrimSpace(input.DocumentID)
	if documentID == "" {
		return in.GetDocumentContentOutput{}, fmt.Errorf("document id is required")
	}
	document, err := s.repo.FindByID(documentID)
	if err != nil {
		return in.GetDocumentContentOutput{}, err
	}
	if strings.TrimSpace(document.ExtractedText) == "" {
		return in.GetDocumentContentOutput{}, fmt.Errorf("document %s has no extracted content", documentID)
	}
	return in.GetDocumentContentOutput{DocumentID: documentID, Content: document.ExtractedText}, nil
}

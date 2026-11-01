package service

import (
	"context"
	"fmt"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
)

type GetDocumentsService struct {
	repo out.UploadRespository
}

func NewGetDocumentsService(repo out.UploadRespository) in.GetDocumentsService {
	return &GetDocumentsService{repo: repo}
}

func (s *GetDocumentsService) Execute(ctx context.Context, input in.GetDocumentsInput) (in.GetDocumentsOutput, error) {
	_ = ctx
	if len(input.DocumentIDs) == 0 {
		return in.GetDocumentsOutput{Documents: []in.DocumentSummary{}}, nil
	}
	items, err := s.repo.FindByIDs(input.DocumentIDs)
	if err != nil {
		return in.GetDocumentsOutput{}, err
	}
	if len(items) != len(input.DocumentIDs) {
		return in.GetDocumentsOutput{}, fmt.Errorf("one or more documents were not found")
	}
	documents := make([]in.DocumentSummary, 0, len(items))
	for _, item := range items {
		documents = append(documents, in.DocumentSummary{
			ID: item.ID, CreatedBy: item.CreatedBy, FileName: item.FileName,
			ContentType: item.ContentType, FileSize: item.FileSize, Status: item.Status,
			CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}
	return in.GetDocumentsOutput{Documents: documents}, nil
}

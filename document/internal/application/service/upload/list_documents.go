package service

import (
	"context"
	"strings"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
)

type ListDocumentsService struct {
	repo out.UploadRespository
}

func NewListDocumentsService(repo out.UploadRespository) in.ListDocumentsService {
	return &ListDocumentsService{repo: repo}
}

func (s *ListDocumentsService) Execute(ctx context.Context, input in.ListDocumentsInput) (in.ListDocumentsOutput, error) {
	_ = ctx
	var createdBy *string
	if strings.TrimSpace(input.CreatedBy) != "" {
		createdBy = &input.CreatedBy
	}
	var search *string
	if value := strings.TrimSpace(input.Search); value != "" {
		search = &value
	}
	items, total, err := s.repo.List(out.UploadListFilters{
		CreatedBy: createdBy, Status: input.Status, Search: search,
	})
	if err != nil {
		return in.ListDocumentsOutput{}, err
	}
	documents := make([]in.DocumentSummary, 0, len(items))
	for _, item := range items {
		documents = append(documents, in.DocumentSummary{
			ID: item.ID, CreatedBy: item.CreatedBy, FileName: item.FileName,
			ContentType: item.ContentType, FileSize: item.FileSize, Status: item.Status,
			CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
		})
	}
	return in.ListDocumentsOutput{Documents: documents, Total: total}, nil
}

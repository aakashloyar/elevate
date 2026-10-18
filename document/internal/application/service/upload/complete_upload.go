package service

import (
	"context"
	"fmt"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type CompleteUploadService struct {
	repo	out.UploadRespository
	storage out.Storage
}

func NewCompleteUploadService(repo out.UploadRespository, storage out.Storage) in.CompleteUploadService {
	return &CompleteUploadService{repo: repo, storage: storage}
}

func (s *CompleteUploadService) Execute(ctx context.Context, input in.CompleteUploadInput) (in.CompleteUploadOutput, error) {

	// 1. Fetch upload
	upload, err := s.repo.FindByID(input.UploadID)
	if err != nil {
		return in.CompleteUploadOutput{}, err
	}

	// 2. Fetch parts
	parts, err := s.repo.GetChunksByUpload(input.UploadID)
	if err != nil {
		return in.CompleteUploadOutput{}, err
	}

	// 3. Validate completeness
	if len(parts) != upload.TotalChunks {
		return in.CompleteUploadOutput{}, fmt.Errorf("missing chunks")
	}

	// 4. Convert to storage format
	var completedParts []domain.Chunk

	for _, p := range parts {
		completedParts = append(completedParts, domain.Chunk{
			UploadID:    input.UploadID,
			ChunkNumber: p.ChunkNumber,
			ETag:        p.ETag,
		})
	}

	// 5. Complete multipart upload
	key := s.storage.BuildStorageKey(upload.ID, upload.FileName)
	err = s.storage.CompleteMultipartUpload(
		ctx,
		key,
		upload.StorageUploadID,
		completedParts,
	)
	if err != nil {
		return in.CompleteUploadOutput{}, err
	}

	err = s.repo.Complete(input.UploadID, domain.StatusMerging)
	if err != nil {
		return in.CompleteUploadOutput{}, err
	}
	err = s.repo.DeleteUploadParts(input.UploadID)
	if err != nil {
		return in.CompleteUploadOutput{}, err
	}
	
	return in.CompleteUploadOutput{}, nil
}

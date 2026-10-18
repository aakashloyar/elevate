package service

import (
	"context"
	"math"

	"github.com/aakashloyar/elevate/document/config"
	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/upload"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type InitUploadService struct {
	repo    out.UploadRespository
	storage out.Storage
	idGen   domain.IDGenerator
	clock   domain.Clock
}

func NewInitUploadService(repo out.UploadRespository, storage out.Storage, idGen domain.IDGenerator, clock domain.Clock) in.InitUploadService {
	return &InitUploadService{repo: repo, storage: storage, idGen: idGen, clock: clock}
}

func (s *InitUploadService) Execute(ctx context.Context, input in.InitUploadInput) (in.InitUploadOutput, error) {

	maxChunkSize := config.App.Upload.MaxChunkSize
	uploadID := s.idGen.NewID()

	totalChunks := int(math.Ceil(float64(input.FileSize) / float64(maxChunkSize)))

	key := s.storage.BuildStorageKey(uploadID, input.FileName)

	// create multipart upload
	storageUploadID, err := s.storage.CreateMultipartUpload(ctx, key)

	if err != nil {
		return in.InitUploadOutput{}, err
	}

	// generate URLs
	var uploadURLs []in.UploadURL
	for i := 1; i <= totalChunks; i++ {
		url, err := s.storage.GeneratePresignedPartURL(ctx, key, storageUploadID, int32(i))
		if err != nil {
			return in.InitUploadOutput{}, err
		}
		uploadURLs = append(uploadURLs, in.UploadURL{
			ChunkNumber: i,
			URL:         url,
		})
	}
	upload := domain.Upload{
		ID:              uploadID,
		FileName:        input.FileName,
		FileSize:        input.FileSize,
		Status:          domain.StatusInitiated,
		StorageUploadID: storageUploadID,
		StorageKey:      key,
		TotalChunks:     totalChunks,
		CreatedAt:       s.clock.Now().UTC(),
		UpdatedAt:       s.clock.Now().UTC(),
	}

	if err := s.repo.Init(upload); err != nil {
		return in.InitUploadOutput{}, err
	}

	return in.InitUploadOutput{
		UploadID:     upload.ID,
		MaxChunkSize: maxChunkSize,
		UploadURLs:   uploadURLs,
	}, nil
}

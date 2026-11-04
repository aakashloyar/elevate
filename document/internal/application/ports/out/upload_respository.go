package out

import "github.com/aakashloyar/elevate/document/internal/domain"

type UploadRespository interface {
	FindByID(uploadID string) (domain.Upload, error)
	Init(upload domain.Upload) error
	SaveExtractedText(uploadID, text string) error
	Complete(uploadID string, status domain.UploadStatus) error
	GetChunksByUpload(uploadID string) ([]domain.Chunk, error)
	DeleteUploadParts(uploadID string) error
	List(filters UploadListFilters) ([]domain.Upload, int, error)
	FindByIDs(uploadIDs []string) ([]domain.Upload, error)
}

type UploadListFilters struct {
	CreatedBy *string
	Status    *domain.UploadStatus
	Search    *string
}

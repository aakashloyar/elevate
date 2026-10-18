package out
import (
	"github.com/aakashloyar/elevate/document/internal/domain"
)
type UploadRespository interface {
	FindByID(uploadID string) (domain.Upload, error) 
	Init(upload domain.Upload) error
	Complete(uploadID string,status domain.UploadStatus) error
	GetChunksByUpload(uploadID string) ([]domain.Chunk, error)
	DeleteUploadParts(uploadID string) error
}

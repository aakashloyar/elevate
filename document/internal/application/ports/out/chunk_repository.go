package out

import (
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type ChunkRespository interface {
	Mark(chunk domain.Chunk) error
}

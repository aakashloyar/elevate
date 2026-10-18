package out

import "github.com/aakashloyar/elevate/document/internal/domain"

type VectorRepository interface {
	SaveChunk(chunk domain.DocumentChunk) error
	SearchSimilar(embedding []float32, limit int) ([]domain.DocumentChunk, error)
}

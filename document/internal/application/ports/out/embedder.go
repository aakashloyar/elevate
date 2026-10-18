package out

import "context"

// Embedder converts document text into vectors suitable for retrieval.
type Embedder interface {
	EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error)
}

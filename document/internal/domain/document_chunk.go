package domain

import "time"

type DocumentChunk struct {
	ID         string
	DocumentID string
	ChunkIndex int
	Content    string
	Embedding  []float32
	Metadata   []byte
	CreatedAt  time.Time
}

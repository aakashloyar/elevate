package domain

import "time"

type Chunk struct {
	UploadID    string
	ChunkNumber int
	ETag        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

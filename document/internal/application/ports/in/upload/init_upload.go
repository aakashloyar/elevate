package upload

import "context"

type InitUploadInput struct {
	FileName string
	FileSize int64
}

type UploadURL struct {
	ChunkNumber int
	URL         string
}

type InitUploadOutput struct {
	UploadID     string
	MaxChunkSize int64
	UploadURLs   []UploadURL
}

type InitUploadService interface {
	Execute(ctx context.Context, input InitUploadInput) (InitUploadOutput, error)
}

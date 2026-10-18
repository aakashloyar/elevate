package in 

import (
	"context"
)

type MarkChunkInput struct {
	UploadID    string
	ChunkNumber int
	ETag        string
}

type MarkChunkOutput struct {
	
}

type MarkChunkService interface{
	Execute(ctx context.Context, input MarkChunkInput) (MarkChunkOutput, error)
}
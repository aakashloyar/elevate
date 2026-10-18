package service 

import (
	"context"
	in "github.com/aakashloyar/elevate/document/internal/application/ports/in/chunk"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/domain"
)

type MarkChunkService struct {
	clock   domain.Clock
	idGen   domain.IDGenerator
	repo    out.ChunkRespository
}

func NewMarkChunkService(chunkRepo out.ChunkRespository, idGen domain.IDGenerator, clock domain.Clock) in.MarkChunkService {
	return &MarkChunkService{
		clock:   clock,
		idGen:   idGen,
		repo:    chunkRepo,
	}
}

func (s *MarkChunkService) Execute(ctx context.Context, input in.MarkChunkInput) (in.MarkChunkOutput, error) {
	chunk := domain.Chunk {
		UploadID: input.UploadID,
		ChunkNumber: input.ChunkNumber,
		ETag: input.ETag,
		CreatedAt: s.clock.Now().UTC().Unix(),
		UpdatedAt: s.clock.Now().UTC().Unix(),
	}

	err := s.repo.Mark(chunk); if err != nil {
		return in.MarkChunkOutput{}, err
	}
	return in.MarkChunkOutput{}, nil
}

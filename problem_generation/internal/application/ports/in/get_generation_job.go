package in

import (
	"context"
	"time"

	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type GetGenerationJobInput struct {
	JobID string
}

type GetGenerationJobOutput struct {
	ID                 string
	UserID             string
	SingleCorrectCount int
	MultiCorrectCount  int
	NumericalCount     int
	DocumentID         *string
	AssessmentID       *string
	Level              domain.GenerationLevel
	Description        string
	Status             domain.GenerationJobStatus
	TopicIDs           []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type GetGenerationJobService interface {
	Execute(ctx context.Context, input GetGenerationJobInput) (GetGenerationJobOutput, error)
}

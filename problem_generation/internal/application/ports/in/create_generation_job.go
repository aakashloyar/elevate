package in

import (
	"context"

	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type CreateGenerationJobInput struct {
	UserID             string
	SingleCorrectCount int
	MultiCorrectCount  int
	NumericalCount     int
	DocumentID         *string
	AssessmentID       *string
	Level              domain.GenerationLevel
	Description        string
	TopicIDs           []string
}

type CreateGenerationJobOutput struct {
	JobID  string
	Status domain.GenerationJobStatus
}

type CreateGenerationJobService interface {
	Execute(ctx context.Context, input CreateGenerationJobInput) (CreateGenerationJobOutput, error)
}

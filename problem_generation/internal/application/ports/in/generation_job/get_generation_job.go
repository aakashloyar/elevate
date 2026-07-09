package generation_job

import (
	"context"
	"time"
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
	Level              string
	Description        string
	Status             string
	TopicIDs           []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type GetGenerationJobService interface {
	Execute(ctx context.Context, input GetGenerationJobInput) (GetGenerationJobOutput, error)
}

package generation_job

import "context"

type CreateGenerationJobInput struct {
	UserID             string
	SingleCorrectCount int
	MultiCorrectCount  int
	NumericalCount     int
	DocumentID         *string
	AssessmentID       *string
	Level              string
	Description        string
	TopicIDs           []string
}

type CreateGenerationJobOutput struct {
	JobID  string
	Status string
}

type CreateGenerationJobService interface {
	Execute(ctx context.Context, input CreateGenerationJobInput) (CreateGenerationJobOutput, error)
}

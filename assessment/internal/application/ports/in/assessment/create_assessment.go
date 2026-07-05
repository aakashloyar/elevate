package assessment

import "context"

type CreateAssessmentInput struct {
	Title           string
	Description     string
	Status          string
	DurationSeconds int
	CreatedBy       string
}

type CreateAssessmentOutput struct {
	AssessmentID string
}

type CreateAssessmentService interface {
	Execute(ctx context.Context, input CreateAssessmentInput) (CreateAssessmentOutput, error)
}

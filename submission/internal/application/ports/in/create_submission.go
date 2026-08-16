package submission

import "context"

type CreateSubmissionInput struct {
	AssessmentID    string
	UserID          string
	DurationSeconds int
}

type CreateSubmissionOutput struct {
	SubmissionID string
	StartedAt    string
}

type CreateSubmissionService interface {
	Execute(ctx context.Context, input CreateSubmissionInput) (CreateSubmissionOutput, error)
}

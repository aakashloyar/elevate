package submission

import "context"

import "github.com/aakashloyar/elevate/submission/internal/domain"

type CreateSubmissionInput struct {
	AssessmentID    string
	UserID          string
	DurationSeconds int
	MarkingScheme   domain.MarkingScheme
}

type CreateSubmissionOutput struct {
	SubmissionID string
	CreatedAt    string
}

type CreateSubmissionService interface {
	Execute(ctx context.Context, input CreateSubmissionInput) (CreateSubmissionOutput, error)
}

package submission
package submission

import "context"

type CreateSubmissionInput struct {
	AssessmentID string
	UserID       string
}

type CreateSubmissionOutput struct {
	SubmissionID string
	StartedAt    string
}

type CreateSubmissionService interface {
	Execute(ctx context.Context, input CreateSubmissionInput) (CreateSubmissionOutput, error)
}

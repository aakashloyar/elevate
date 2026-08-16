package submission

import "context"

type StartSubmissionInput struct {
	SubmissionID string
}

type StartSubmissionService interface {
	Execute(ctx context.Context, input StartSubmissionInput) error
}

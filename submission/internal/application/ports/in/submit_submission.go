package submission

import "context"

type SubmitSubmissionInput struct {
	SubmissionID string
}

type SubmitSubmissionService interface {
	Execute(ctx context.Context, input SubmitSubmissionInput) error
}

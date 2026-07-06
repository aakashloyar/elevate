package submission

import "context"

type SubmitAnswerInput struct {
	SubmissionID string
	ProblemID    string
	Answer       []string
}

type SubmitAnswerOutput struct{}

type SubmitAnswerService interface {
	Execute(ctx context.Context, input SubmitAnswerInput) error
}

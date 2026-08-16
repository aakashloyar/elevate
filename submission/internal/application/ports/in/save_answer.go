package submission

import "context"

type SaveAnswerInput struct {
	SubmissionID string
	ProblemID    string
	Answer       []string
}

type SaveAnswerOutput struct{}

type SaveAnswerService interface {
	Execute(ctx context.Context, input SaveAnswerInput) error
}
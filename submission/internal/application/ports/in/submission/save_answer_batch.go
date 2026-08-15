package submission

import "context"

type SaveAnswerBatchInput struct {
	SubmissionID string
	Answers      []SaveAnswerBatchItem
}

type SaveAnswerBatchItem struct {
	ProblemID string
	Answer    []string
}

type SaveAnswerBatchOutput struct{}

type SaveAnswerBatchService interface {
	Execute(ctx context.Context, input SaveAnswerBatchInput) error
}

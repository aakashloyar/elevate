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

type SaveAnswerBatchOutput struct {
	SavedCount int
	Errors     []SaveAnswerBatchItemError
}

type SaveAnswerBatchItemError struct {
	ProblemID string
	Message   string
}

type SaveAnswerBatchService interface {
	Execute(ctx context.Context, input SaveAnswerBatchInput) (SaveAnswerBatchOutput, error)
}

package problem

import "context"

type CreateProblemBatchInput struct {
	Problems     []CreateProblemInput
	AssessmentID string
}

type CreateProblemBatchOutput struct {}

type CreateProblemBatchService interface {
	Execute(ctx context.Context, input CreateProblemBatchInput) (CreateProblemBatchOutput, error)
}

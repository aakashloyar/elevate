package assessment

import "context"

type AddProblemsBatchInput struct {
	AssessmentID string
	ProblemIDs   []string
}

type AddProblemsBatchOutput struct{}

type AddProblemsBatchService interface {
	Execute(ctx context.Context, input AddProblemsBatchInput) (AddProblemsBatchOutput, error)
}

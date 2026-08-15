package problem

import "context"

type DeleteProblemInput struct {
	ProblemID string
}

type DeleteProblemOutput struct{}

type DeleteProblemService interface {
	Execute(ctx context.Context, input DeleteProblemInput) error
}

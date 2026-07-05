package problem

import "context"

type UpdateProblemInput struct {
	ProblemID  string
	CreatedBy  string
	Title      string
	Statement  string
	Type       string
	Difficulty string
	SourceType string
	Status     string
	Options    []CreateProblemOptionInput
	Tags       []string
}

type UpdateProblemOutput struct{}

type UpdateProblemService interface {
	Execute(ctx context.Context, input UpdateProblemInput) error
}

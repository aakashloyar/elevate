package problem

import "context"

type CreateProblemInput struct {
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

type CreateProblemOptionInput struct {
	Text      string
	IsCorrect bool
}

type CreateProblemOutput struct {
	ProblemID string
}

type CreateProblemService interface {
	Execute(ctx context.Context, input CreateProblemInput) (CreateProblemOutput, error)
}

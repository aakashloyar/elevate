package problem

import (
	"context"

	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type CreateProblemInput struct {
	CreatedBy  string
	Title      string
	Statement  string
	Type       domain.ProblemType
	Difficulty domain.Difficulty
	SourceType domain.SourceType
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

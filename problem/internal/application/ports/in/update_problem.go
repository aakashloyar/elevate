package problem

import (
	"context"

	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type UpdateProblemInput struct {
	ProblemID  string
	CreatedBy  string
	Title      string
	Statement  string
	Type       domain.ProblemType
	Difficulty domain.Difficulty
	SourceType domain.SourceType
	Options    []CreateProblemOptionInput
	Tags       []string
}

type UpdateProblemOutput struct{}

type UpdateProblemService interface {
	Execute(ctx context.Context, input UpdateProblemInput) error
}

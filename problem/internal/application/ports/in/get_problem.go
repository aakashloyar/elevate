package problem

import (
	"context"
	"time"

	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type GetProblemInput struct {
	ProblemID string
}

type GetProblemOutput struct {
	ID         string
	CreatedBy  string
	Title      string
	Statement  string
	Type       domain.ProblemType
	Difficulty domain.Difficulty
	SourceType domain.SourceType
	Options    []ProblemOptionOutput
	Tags       []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ProblemOptionOutput struct {
	ID        string
	Text      string
	IsCorrect bool
}

type GetProblemService interface {
	Execute(ctx context.Context, input GetProblemInput) (GetProblemOutput, error)
}

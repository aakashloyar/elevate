package problem

import (
	"context"

	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type ListProblemsInput struct {
	Offset     int
	Limit      int
	CreatedBy  string
	Title      string
	Type       string
	Difficulty string
	SourceType string
	Tag        string
}

type ListProblemsOutput struct {
	Problems []ListProblemItem
}

type ListProblemItem struct {
	ID         string
	Title      string
	Type       domain.ProblemType
	Difficulty domain.Difficulty
	CreatedAt  string
}

type ListProblemsService interface {
	Execute(ctx context.Context, input ListProblemsInput) (ListProblemsOutput, error)
}

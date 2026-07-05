package problem

import "context"

type ListProblemsInput struct {
	Offset int
	Limit  int
	Type   string
	Status string
	Tag    string
}

type ListProblemsOutput struct {
	Problems []ListProblemItem
}

type ListProblemItem struct {
	ID         string
	Title      string
	Type       string
	Difficulty string
	Status     string
	CreatedAt  string
}

type ListProblemsService interface {
	Execute(ctx context.Context, input ListProblemsInput) (ListProblemsOutput, error)
}

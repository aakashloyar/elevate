package problem

import (
	"context"
	"time"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in/problem"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
)

type ListProblemsService struct {
	problemRepo out.ProblemRepository
}

func NewListProblemsService(problemRepo out.ProblemRepository) in.ListProblemsService {
	return &ListProblemsService{problemRepo: problemRepo}
}

func (s *ListProblemsService) Execute(ctx context.Context, input in.ListProblemsInput) (in.ListProblemsOutput, error) {
	filters := map[string]string{}
	if input.Type != "" {
		filters["type"] = input.Type
	}
	if input.Status != "" {
		filters["status"] = input.Status
	}
	if input.Tag != "" {
		filters["tag"] = input.Tag
	}

	problems, err := s.problemRepo.List(input.Offset, input.Limit, filters)
	if err != nil {
		return in.ListProblemsOutput{}, err
	}

	items := make([]in.ListProblemItem, 0, len(problems))
	for _, p := range problems {
		items = append(items, in.ListProblemItem{
			ID:         p.ID,
			Title:      p.Title,
			Type:       p.Type,
			Difficulty: p.Difficulty,
			Status:     p.Status,
			CreatedAt:  p.CreatedAt.Format(time.RFC3339),
		})
	}

	return in.ListProblemsOutput{Problems: items}, nil
}

package problem

import (
	"context"
	"time"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in"
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
	if input.CreatedBy != "" {
		filters["created_by"] = input.CreatedBy
	}
	if input.Title != "" {
		filters["title"] = input.Title
	}
	if input.Type != "" {
		filters["type"] = input.Type
	}
	if input.Difficulty != "" {
		filters["difficulty"] = input.Difficulty
	}
	if input.SourceType != "" {
		filters["source_type"] = input.SourceType
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
			CreatedAt:  p.CreatedAt.Format(time.RFC3339),
		})
	}

	return in.ListProblemsOutput{Problems: items}, nil
}

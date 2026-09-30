package problem

import (
	"context"
	"time"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
)

type ListProblemsService struct {
	problemRepo out.ProblemRepository
	userClient  out.UserClient
}

func NewListProblemsService(problemRepo out.ProblemRepository, userClient out.UserClient) in.ListProblemsService {
	return &ListProblemsService{problemRepo: problemRepo, userClient: userClient}
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

	userIDs := make([]string, 0, len(problems))
	for _, p := range problems {
		userIDs = append(userIDs, p.CreatedBy)
	}
	usernames, err := s.userClient.GetUsernames(ctx, userIDs)
	if err != nil {
		return in.ListProblemsOutput{}, err
	}
	items := make([]in.ListProblemItem, 0, len(problems))
	for _, p := range problems {
		items = append(items, in.ListProblemItem{
			ID:            p.ID,
			CreatedBy:     p.CreatedBy,
			CreatedByName: usernames[p.CreatedBy],
			Title:         p.Title,
			Type:          p.Type,
			Difficulty:    p.Difficulty,
			SourceType:    p.SourceType,
			Tags:          p.Tags,
			CreatedAt:     p.CreatedAt.Format(time.RFC3339),
		})
	}

	return in.ListProblemsOutput{Problems: items}, nil
}

package problem

import (
	"context"
	"strings"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in/problem"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type UpdateProblemService struct {
	problemRepo out.ProblemRepository
	clock       out.Clock
}

func NewUpdateProblemService(problemRepo out.ProblemRepository, clock out.Clock) in.UpdateProblemService {
	return &UpdateProblemService{problemRepo: problemRepo, clock: clock}
}

func (s *UpdateProblemService) Execute(ctx context.Context, input in.UpdateProblemInput) error {
	problem := domain.Problem{
		ID:         input.ProblemID,
		CreatedBy:  strings.TrimSpace(input.CreatedBy),
		Title:      strings.TrimSpace(input.Title),
		Statement:  strings.TrimSpace(input.Statement),
		Type:       strings.ToUpper(strings.TrimSpace(input.Type)),
		Difficulty: strings.ToUpper(strings.TrimSpace(input.Difficulty)),
		SourceType: strings.ToUpper(strings.TrimSpace(input.SourceType)),
		Status:     strings.ToUpper(strings.TrimSpace(input.Status)),
		UpdatedAt:  s.clock.Now(),
	}
	return s.problemRepo.Update(problem)
}

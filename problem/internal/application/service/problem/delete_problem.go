package problem

import (
	"context"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in/problem"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
)

type DeleteProblemService struct {
	problemRepo out.ProblemRepository
}

func NewDeleteProblemService(problemRepo out.ProblemRepository) in.DeleteProblemService {
	return &DeleteProblemService{problemRepo: problemRepo}
}

func (s *DeleteProblemService) Execute(ctx context.Context, input in.DeleteProblemInput) error {
	return s.problemRepo.DeleteByID(input.ProblemID)
}

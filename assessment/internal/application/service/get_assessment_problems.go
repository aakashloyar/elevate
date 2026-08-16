package assessment

import (
	"context"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type GetAssessmentProblemsService struct {
	repository out.AssessmentRepository
}

func NewGetAssessmentProblemsService(repository out.AssessmentRepository) in.GetAssessmentProblemsService {
	return &GetAssessmentProblemsService{repository: repository}
}

func (s *GetAssessmentProblemsService) Execute(ctx context.Context, input in.GetAssessmentProblemsInput) (in.GetAssessmentProblemsOutput, error) {
	problemIDs, err := s.repository.FindProblemIDs(input.AssessmentID)
	return in.GetAssessmentProblemsOutput{ProblemIDs: problemIDs}, err
}

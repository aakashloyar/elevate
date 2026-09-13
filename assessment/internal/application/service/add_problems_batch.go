package assessment

import (
	"context"
	"log"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type AddProblemsBatchService struct {
	assessmentRepo out.AssessmentRepository
	problemClient  out.ProblemClient
}

func NewAddProblemsBatchService(assessmentRepo out.AssessmentRepository, problemClient out.ProblemClient) in.AddProblemsBatchService {
	return &AddProblemsBatchService{assessmentRepo: assessmentRepo, problemClient: problemClient}
}

func (s *AddProblemsBatchService) Execute(ctx context.Context, input in.AddProblemsBatchInput) (in.AddProblemsBatchOutput, error) {
	if input.AssessmentID == "" || len(input.ProblemIDs) == 0 {
		return in.AddProblemsBatchOutput{}, nil
	}
	for _, problemID := range input.ProblemIDs {
		if err := s.problemClient.Exists(ctx, problemID); err != nil {
			return in.AddProblemsBatchOutput{}, err
		}
	}
	if err := s.assessmentRepo.AddProblems(input.AssessmentID, input.ProblemIDs); err != nil {
		log.Printf("failed to add problems to assessment %s: %v", input.AssessmentID, err)
		return in.AddProblemsBatchOutput{}, err
	}
	return in.AddProblemsBatchOutput{}, nil
}

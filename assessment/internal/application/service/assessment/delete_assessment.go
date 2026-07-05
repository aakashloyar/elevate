package assessment

import (
	"context"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in/assessment"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type DeleteAssessmentService struct {
	assessmentRepo out.AssessmentRepository
}

func NewDeleteAssessmentService(assessmentRepo out.AssessmentRepository) in.DeleteAssessmentService {
	return &DeleteAssessmentService{assessmentRepo: assessmentRepo}
}

func (s *DeleteAssessmentService) Execute(ctx context.Context, input in.DeleteAssessmentInput) error {
	return s.assessmentRepo.DeleteByID(input.AssessmentID)
}

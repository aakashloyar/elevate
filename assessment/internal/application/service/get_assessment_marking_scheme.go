package assessment

import (
	"context"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
	"github.com/aakashloyar/elevate/assessment/internal/domain"
)

type GetAssessmentMarkingSchemeService struct {
	assessmentRepo out.AssessmentRepository
}

func NewGetAssessmentMarkingSchemeService(assessmentRepo out.AssessmentRepository) in.GetAssessmentMarkingSchemeService {
	return &GetAssessmentMarkingSchemeService{assessmentRepo: assessmentRepo}
}

func (s *GetAssessmentMarkingSchemeService) Execute(ctx context.Context, input in.GetAssessmentMarkingSchemeInput) (domain.AssessmentMarkingScheme, error) {
	return s.assessmentRepo.FindMarkingScheme(input.AssessmentID)
}

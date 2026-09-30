package assessment

import (
	"context"
	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type GetAssessmentsBatchService struct{ repository out.AssessmentRepository }

func NewGetAssessmentsBatchService(repository out.AssessmentRepository) in.GetAssessmentsBatchService {
	return &GetAssessmentsBatchService{repository: repository}
}
func (s *GetAssessmentsBatchService) Execute(_ context.Context, input in.GetAssessmentsBatchInput) (in.GetAssessmentsBatchOutput, error) {
	assessments, err := s.repository.FindByIDs(input.AssessmentIDs)
	if err != nil {
		return in.GetAssessmentsBatchOutput{}, err
	}
	result := make([]in.AssessmentSummary, 0, len(assessments))
	for _, assessment := range assessments {
		result = append(result, in.AssessmentSummary{ID: assessment.ID, Title: assessment.Title})
	}
	return in.GetAssessmentsBatchOutput{Assessments: result}, nil
}

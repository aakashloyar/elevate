package assessment

import (
	"context"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in/assessment"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type GetAssessmentService struct {
	assessmentRepo out.AssessmentRepository
}

func NewGetAssessmentService(assessmentRepo out.AssessmentRepository) in.GetAssessmentService {
	return &GetAssessmentService{assessmentRepo: assessmentRepo}
}

func (s *GetAssessmentService) Execute(ctx context.Context, input in.GetAssessmentInput) (in.GetAssessmentOutput, error) {
	assessment, err := s.assessmentRepo.FindByID(input.AssessmentID)
	if err != nil {
		return in.GetAssessmentOutput{}, err
	}

	return in.GetAssessmentOutput{
		ID:              assessment.ID,
		Title:           assessment.Title,
		Description:     assessment.Description,
		Status:          assessment.Status,
		DurationSeconds: assessment.DurationSeconds,
		CreatedBy:       assessment.CreatedBy,
		CreatedAt:       assessment.CreatedAt,
		UpdatedAt:       assessment.UpdatedAt,
	}, nil
}

package assessment

import (
	"context"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type GetAssessmentService struct {
	assessmentRepo out.AssessmentRepository
	userClient     out.UserClient
}

func NewGetAssessmentService(assessmentRepo out.AssessmentRepository, userClient out.UserClient) in.GetAssessmentService {
	return &GetAssessmentService{assessmentRepo: assessmentRepo, userClient: userClient}
}

func (s *GetAssessmentService) Execute(ctx context.Context, input in.GetAssessmentInput) (in.GetAssessmentOutput, error) {
	assessment, err := s.assessmentRepo.FindByID(input.AssessmentID)
	if err != nil {
		return in.GetAssessmentOutput{}, err
	}
	usernames, err := s.userClient.GetUsernames(ctx, []string{assessment.CreatedBy})
	if err != nil {
		return in.GetAssessmentOutput{}, err
	}

	return in.GetAssessmentOutput{
		ID:              assessment.ID,
		Title:           assessment.Title,
		Description:     assessment.Description,
		DurationSeconds: assessment.DurationSeconds,
		CreatedBy:       assessment.CreatedBy,
		CreatedByName:   usernames[assessment.CreatedBy],
		CreatedAt:       assessment.CreatedAt,
		UpdatedAt:       assessment.UpdatedAt,
	}, nil
}

package assessment

import (
	"context"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out"
)

type ListAssessmentsService struct {
	assessmentRepo out.AssessmentRepository
	userClient     out.UserClient
}

func NewListAssessmentsService(assessmentRepo out.AssessmentRepository, userClient out.UserClient) in.ListAssessmentsService {
	return &ListAssessmentsService{assessmentRepo: assessmentRepo, userClient: userClient}
}

func (s *ListAssessmentsService) Execute(ctx context.Context, input in.ListAssessmentsInput) (in.ListAssessmentsOutput, error) {
	assessments, err := s.assessmentRepo.FindAll(out.FindAllAssessmentFilters{
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		Limit:       input.Limit,
		Offset:      input.Offset,
	})
	if err != nil {
		return in.ListAssessmentsOutput{}, err
	}

	userIDs := make([]string, 0, len(assessments))
	for _, assessment := range assessments {
		userIDs = append(userIDs, assessment.CreatedBy)
	}
	usernames, err := s.userClient.GetUsernames(ctx, userIDs)
	if err != nil {
		return in.ListAssessmentsOutput{}, err
	}
	outputs := make([]in.GetAssessmentOutput, 0, len(assessments))
	for _, assessment := range assessments {
		outputs = append(outputs, in.GetAssessmentOutput{
			ID:              assessment.ID,
			Title:           assessment.Title,
			Description:     assessment.Description,
			DurationSeconds: assessment.DurationSeconds,
			CreatedBy:       assessment.CreatedBy,
			CreatedByName:   usernames[assessment.CreatedBy],
			CreatedAt:       assessment.CreatedAt,
			UpdatedAt:       assessment.UpdatedAt,
		})
	}

	return in.ListAssessmentsOutput{Assessments: outputs}, nil
}

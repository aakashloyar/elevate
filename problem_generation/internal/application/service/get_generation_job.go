package service

import (
	"context"
	"errors"

	in "github.com/aakashloyar/elevate/problem_generation/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
)

type GetGenerationJobService struct {
	jobRepo          out.GenerationJobRepository
	assessmentClient out.AssessmentClient
}

func NewGetGenerationJobService(jobRepo out.GenerationJobRepository, assessmentClient out.AssessmentClient) in.GetGenerationJobService {
	return &GetGenerationJobService{jobRepo: jobRepo, assessmentClient: assessmentClient}
}

func (s *GetGenerationJobService) Execute(ctx context.Context, input in.GetGenerationJobInput) (in.GetGenerationJobOutput, error) {
	if input.JobID == "" {
		return in.GetGenerationJobOutput{}, errors.New("job id is required")
	}

	job, err := s.jobRepo.FindByID(input.JobID)
	if err != nil {
		return in.GetGenerationJobOutput{}, err
	}
	assessmentTitles := map[string]string{}
	if assessmentID := valueOrEmpty(job.AssessmentID); assessmentID != "" {
		assessmentTitles, err = s.assessmentClient.GetTitles(ctx, []string{assessmentID})
		if err != nil {
			return in.GetGenerationJobOutput{}, err
		}
	}

	return in.GetGenerationJobOutput{
		ID:                    job.ID,
		UserID:                job.UserID,
		SingleCorrectCount:    job.SingleCorrectCount,
		MultiCorrectCount:     job.MultiCorrectCount,
		NumericalCount:        job.NumericalCount,
		DocumentID:            job.DocumentID,
		AssessmentID:          job.AssessmentID,
		AssessmentTitle:       assessmentTitles[valueOrEmpty(job.AssessmentID)],
		Level:                 job.Level,
		Description:           job.Description,
		Status:                job.Status,
		TopicIDs:              job.TopicIDs,
		GeneratedProblemCount: job.GeneratedProblemCount,
		CreatedAt:             job.CreatedAt,
		UpdatedAt:             job.UpdatedAt,
	}, nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

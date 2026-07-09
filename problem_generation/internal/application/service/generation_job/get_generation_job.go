package generation_job

import (
	"context"
	"errors"

	in "github.com/aakashloyar/elevate/problem_generation/internal/application/ports/in/generation_job"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
)

type GetGenerationJobService struct {
	jobRepo out.GenerationJobRepository
}

func NewGetGenerationJobService(jobRepo out.GenerationJobRepository) in.GetGenerationJobService {
	return &GetGenerationJobService{jobRepo: jobRepo}
}

func (s *GetGenerationJobService) Execute(ctx context.Context, input in.GetGenerationJobInput) (in.GetGenerationJobOutput, error) {
	if input.JobID == "" {
		return in.GetGenerationJobOutput{}, errors.New("job id is required")
	}

	job, err := s.jobRepo.FindByID(input.JobID)
	if err != nil {
		return in.GetGenerationJobOutput{}, err
	}

	return in.GetGenerationJobOutput{
		ID:                 job.ID,
		UserID:             job.UserID,
		SingleCorrectCount: job.SingleCorrectCount,
		MultiCorrectCount:  job.MultiCorrectCount,
		NumericalCount:     job.NumericalCount,
		DocumentID:         job.DocumentID,
		AssessmentID:       job.AssessmentID,
		Level:              job.Level,
		Description:        job.Description,
		Status:             job.Status,
		TopicIDs:           job.TopicIDs,
		CreatedAt:          job.CreatedAt,
		UpdatedAt:          job.UpdatedAt,
	}, nil
}

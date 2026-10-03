package service

import (
	"context"

	in "github.com/aakashloyar/elevate/problem_generation/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type ListGenerationJobsService struct {
	jobRepo          out.GenerationJobRepository
	assessmentClient out.AssessmentClient
}

func NewListGenerationJobsService(jobRepo out.GenerationJobRepository, assessmentClient out.AssessmentClient) in.ListGenerationJobsService {
	return &ListGenerationJobsService{jobRepo: jobRepo, assessmentClient: assessmentClient}
}

func (s *ListGenerationJobsService) Execute(ctx context.Context, input in.ListGenerationJobsInput) (in.ListGenerationJobsOutput, error) {
	jobs, err := s.jobRepo.FindAll(input.Limit, input.Offset, input.Search)
	if err != nil {
		return in.ListGenerationJobsOutput{}, err
	}

	assessmentIDs := make([]string, 0, len(jobs))
	seenAssessmentIDs := make(map[string]struct{}, len(jobs))
	for _, job := range jobs {
		if job.AssessmentID != nil && *job.AssessmentID != "" {
			if _, exists := seenAssessmentIDs[*job.AssessmentID]; !exists {
				seenAssessmentIDs[*job.AssessmentID] = struct{}{}
				assessmentIDs = append(assessmentIDs, *job.AssessmentID)
			}
		}
	}
	assessmentTitles, err := s.assessmentClient.GetTitles(ctx, assessmentIDs)
	if err != nil {
		return in.ListGenerationJobsOutput{}, err
	}

	result := make([]in.GetGenerationJobOutput, 0, len(jobs))
	for _, job := range jobs {
		output := generationJobOutput(job)
		if job.AssessmentID != nil {
			output.AssessmentTitle = assessmentTitles[*job.AssessmentID]
		}
		result = append(result, output)
	}
	return in.ListGenerationJobsOutput{Jobs: result}, nil
}

func generationJobOutput(job domain.GenerationJob) in.GetGenerationJobOutput {
	return in.GetGenerationJobOutput{
		ID:                    job.ID,
		UserID:                job.UserID,
		SingleCorrectCount:    job.SingleCorrectCount,
		MultiCorrectCount:     job.MultiCorrectCount,
		NumericalCount:        job.NumericalCount,
		DocumentID:            job.DocumentID,
		AssessmentID:          job.AssessmentID,
		Level:                 job.Level,
		Description:           job.Description,
		Status:                job.Status,
		TopicIDs:              job.TopicIDs,
		GeneratedProblemCount: job.GeneratedProblemCount,
		CreatedAt:             job.CreatedAt,
		UpdatedAt:             job.UpdatedAt,
	}
}

package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type ProcessGenerationJobService struct {
	jobRepo   out.GenerationJobRepository
	processor out.GenerationProcessor
}

func NewProcessGenerationJobService(jobRepo out.GenerationJobRepository, processor out.GenerationProcessor) *ProcessGenerationJobService {
	return &ProcessGenerationJobService{jobRepo: jobRepo, processor: processor}
}

func (s *ProcessGenerationJobService) Execute(ctx context.Context, event out.GenerationRequestedEvent) error {
	if strings.TrimSpace(event.JobID) == "" {
		return errors.New("job id is required")
	}

	if err := s.jobRepo.UpdateStatus(event.JobID, domain.GenerationJobStatusProcessing); err != nil {
		return fmt.Errorf("mark generation job processing: %w", err)
	}

	if err := s.processor.ProcessGeneration(ctx, event); err != nil {
		if statusErr := s.jobRepo.UpdateStatus(event.JobID, domain.GenerationJobStatusFailed); statusErr != nil {
			return fmt.Errorf("%w; additionally failed to mark generation job failed: %v", err, statusErr)
		}
		return err
	}

	if err := s.jobRepo.UpdateStatus(event.JobID, domain.GenerationJobStatusCompleted); err != nil {
		return fmt.Errorf("mark generation job completed: %w", err)
	}
	return nil
}

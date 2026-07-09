package generation_job

import (
	"context"
	"errors"
	"fmt"
	"strings"

	in "github.com/aakashloyar/elevate/problem_generation/internal/application/ports/in/generation_job"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type CreateGenerationJobService struct {
	jobRepo        out.GenerationJobRepository
	eventPublisher out.EventPublisher
	idGen          out.IDGenerator
	clock          out.Clock
}

func NewCreateGenerationJobService(jobRepo out.GenerationJobRepository, eventPublisher out.EventPublisher, idGen out.IDGenerator, clock out.Clock) in.CreateGenerationJobService {
	return &CreateGenerationJobService{jobRepo: jobRepo, eventPublisher: eventPublisher, idGen: idGen, clock: clock}
}

func (s *CreateGenerationJobService) Execute(ctx context.Context, input in.CreateGenerationJobInput) (in.CreateGenerationJobOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return in.CreateGenerationJobOutput{}, errors.New("user id is required")
	}

	level := strings.TrimSpace(input.Level)
	if level == "" {
		return in.CreateGenerationJobOutput{}, errors.New("level is required")
	}

	description := strings.TrimSpace(input.Description)
	if input.SingleCorrectCount < 0 || input.MultiCorrectCount < 0 || input.NumericalCount < 0 {
		return in.CreateGenerationJobOutput{}, errors.New("question counts cannot be negative")
	}

	if input.SingleCorrectCount+input.MultiCorrectCount+input.NumericalCount == 0 {
		return in.CreateGenerationJobOutput{}, errors.New("at least one question must be requested")
	}

	now := s.clock.Now()
	job := domain.GenerationJob{
		ID:                 s.idGen.NewID(),
		UserID:             userID,
		SingleCorrectCount: input.SingleCorrectCount,
		MultiCorrectCount:  input.MultiCorrectCount,
		NumericalCount:     input.NumericalCount,
		DocumentID:         normalizeOptional(input.DocumentID),
		AssessmentID:       normalizeOptional(input.AssessmentID),
		Level:              level,
		Description:        description,
		Status:             "pending",
		TopicIDs:           normalizeTopicIDs(input.TopicIDs),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.jobRepo.Save(job); err != nil {
		return in.CreateGenerationJobOutput{}, err
	}

	prompt := buildPrompt(job)
	if err := s.eventPublisher.PublishGenerationRequested(ctx, job, prompt); err != nil {
		return in.CreateGenerationJobOutput{}, fmt.Errorf("publish generation request: %w", err)
	}

	return in.CreateGenerationJobOutput{JobID: job.ID, Status: job.Status}, nil
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizeTopicIDs(topicIDs []string) []string {
	normalized := make([]string, 0, len(topicIDs))
	for _, topicID := range topicIDs {
		trimmed := strings.TrimSpace(topicID)
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
}

func buildPrompt(job domain.GenerationJob) string {
	parts := []string{
		fmt.Sprintf("Generate %d single-correct questions", job.SingleCorrectCount),
		fmt.Sprintf("%d multi-correct questions", job.MultiCorrectCount),
		fmt.Sprintf("%d numerical questions", job.NumericalCount),
		fmt.Sprintf("for level %s", job.Level),
	}

	if job.Description != "" {
		parts = append(parts, fmt.Sprintf("based on: %s", job.Description))
	}

	if job.DocumentID != nil {
		parts = append(parts, fmt.Sprintf("using document %s", *job.DocumentID))
	}

	if job.AssessmentID != nil {
		parts = append(parts, fmt.Sprintf("for assessment %s", *job.AssessmentID))
	}

	if len(job.TopicIDs) > 0 {
		parts = append(parts, fmt.Sprintf("topics: %s", strings.Join(job.TopicIDs, ", ")))
	}

	return strings.Join(parts, "; ")
}

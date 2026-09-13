package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	in "github.com/aakashloyar/elevate/problem_generation/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type CreateGenerationJobService struct {
	jobRepo          out.GenerationJobRepository
	eventPublisher   out.EventPublisher
	userClient       out.UserClient
	assessmentClient out.AssessmentClient
	idGen            out.IDGenerator
	clock            out.Clock
}

func NewCreateGenerationJobService(jobRepo out.GenerationJobRepository, eventPublisher out.EventPublisher, userClient out.UserClient, assessmentClient out.AssessmentClient, idGen out.IDGenerator, clock out.Clock) in.CreateGenerationJobService {
	return &CreateGenerationJobService{jobRepo: jobRepo, eventPublisher: eventPublisher, userClient: userClient, assessmentClient: assessmentClient, idGen: idGen, clock: clock}
}

func (s *CreateGenerationJobService) Execute(ctx context.Context, input in.CreateGenerationJobInput) (in.CreateGenerationJobOutput, error) {
	err := validateInput(input)
	if err != nil {
		return in.CreateGenerationJobOutput{}, err
	}
	if err := s.userClient.Exists(ctx, input.UserID); err != nil {
		return in.CreateGenerationJobOutput{}, fmt.Errorf("user id is invalid: %w", err)
	}
	if input.AssessmentID != nil && strings.TrimSpace(*input.AssessmentID) != "" {
		if err := s.assessmentClient.Exists(ctx, strings.TrimSpace(*input.AssessmentID)); err != nil {
			return in.CreateGenerationJobOutput{}, fmt.Errorf("assessment id is invalid: %w", err)
		}
	}

	level := domain.GenerationLevel(strings.ToLower(strings.TrimSpace(string(input.Level))))

	description := strings.TrimSpace(input.Description)

	now := s.clock.Now()
	job := domain.GenerationJob{
		ID:                 s.idGen.NewID(),
		UserID:             input.UserID,
		SingleCorrectCount: input.SingleCorrectCount,
		MultiCorrectCount:  input.MultiCorrectCount,
		NumericalCount:     input.NumericalCount,
		DocumentID:         normalizeOptional(input.DocumentID),
		AssessmentID:       normalizeOptional(input.AssessmentID),
		Level:              level,
		Description:        description,
		Status:             domain.GenerationJobStatusPending,
		TopicIDs:           normalizeTopicIDs(input.TopicIDs),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.jobRepo.Save(job); err != nil {
		return in.CreateGenerationJobOutput{}, err
	}

	if err := s.eventPublisher.PublishGenerationRequested(ctx, out.GenerationRequestedMessage{
		Event: out.GenerationRequestedEvent{
			JobID: job.ID,
		},
	}); err != nil {
		return in.CreateGenerationJobOutput{}, fmt.Errorf("publish generation request: %w", err)
	}

	return in.CreateGenerationJobOutput{JobID: job.ID, Status: job.Status}, nil
}

func validateInput(input in.CreateGenerationJobInput) error {
	if input.UserID == "" {
		return errors.New("user id is required")
	}

	if input.SingleCorrectCount < 0 || input.MultiCorrectCount < 0 || input.NumericalCount < 0 {
		return errors.New("question counts cannot be negative")
	}

	if input.SingleCorrectCount+input.MultiCorrectCount+input.NumericalCount == 0 {
		return errors.New("at least one question must be requested")
	}

	normalized := domain.GenerationLevel(strings.ToLower(strings.TrimSpace(string(input.Level))))
	if !normalized.IsValid() {
		return errors.New("level must be easy, medium, or hard")
	}
	return nil
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

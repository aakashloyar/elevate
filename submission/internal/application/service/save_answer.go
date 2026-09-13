package submission

import (
	"context"
	"errors"
	"fmt"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SaveAnswerService struct {
	submissionRepo out.SubmissionRepository
	problemClient  out.ProblemClient
	clock          out.Clock
}

func NewSaveAnswerService(submissionRepo out.SubmissionRepository, problemClient out.ProblemClient, clock out.Clock) in.SaveAnswerService {
	return &SaveAnswerService{submissionRepo: submissionRepo, problemClient: problemClient, clock: clock}
}

func (s *SaveAnswerService) Execute(ctx context.Context, input in.SaveAnswerInput) error {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return errors.New("submission id is required")
	}
	if strings.TrimSpace(input.ProblemID) == "" {
		return errors.New("problem id is required")
	}
	if err := s.problemClient.Exists(ctx, input.ProblemID); err != nil {
		return fmt.Errorf("problem id is invalid: %w", err)
	}

	answer := make([]string, 0, len(input.Answer))
	for _, value := range input.Answer {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		answer = append(answer, trimmed)
	}

	now := s.clock.Now()
	saved, err := s.submissionRepo.SaveAnswer(domain.SubmissionAnswer{
		ID:           "",
		SubmissionID: input.SubmissionID,
		ProblemID:    input.ProblemID,
		Answer:       answer,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return err
	}
	if !saved {
		return errors.New("submission must be in IN_PROGRESS state")
	}

	return nil
}

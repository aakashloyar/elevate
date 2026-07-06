package submission

import (
	"context"
	"errors"
	"strings"
	"time"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in/submission"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmitAnswerService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
}

func NewSubmitAnswerService(submissionRepo out.SubmissionRepository, clock out.Clock) in.SubmitAnswerService {
	return &SubmitAnswerService{submissionRepo: submissionRepo, clock: clock}
}

func (s *SubmitAnswerService) Execute(ctx context.Context, input in.SubmitAnswerInput) error {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return errors.New("submission id is required")
	}
	if strings.TrimSpace(input.ProblemID) == "" {
		return errors.New("problem id is required")
	}

	answer := make([]string, 0, len(input.Answer))
	for _, value := range input.Answer {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			answer = append(answer, trimmed)
		}
	}

	if err := s.submissionRepo.SaveAnswer(domain.SubmissionAnswer{
		ID:           "",
		SubmissionID: input.SubmissionID,
		ProblemID:    input.ProblemID,
		Answer:       answer,
		AnsweredAt:   s.clock.Now(),
	}); err != nil {
		return err
	}

	_ = time.Now()
	return nil
}

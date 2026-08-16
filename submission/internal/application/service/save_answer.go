package submission

import (
	"context"
	"errors"
	"strings"
	"time"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SaveAnswerService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
}

func NewSaveAnswerService(submissionRepo out.SubmissionRepository, clock out.Clock) in.SaveAnswerService {
	return &SaveAnswerService{submissionRepo: submissionRepo, clock: clock}
}

func (s *SaveAnswerService) Execute(ctx context.Context, input in.SaveAnswerInput) error {
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
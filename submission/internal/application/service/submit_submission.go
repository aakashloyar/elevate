package submission

import (
	"context"
	"errors"
	"fmt"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

type SubmitSubmissionService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
	publisher      out.EventPublisher
	topic          string
}

func NewSubmitSubmissionService(submissionRepo out.SubmissionRepository, clock out.Clock, publisher out.EventPublisher, topic string) in.SubmitSubmissionService {
	return &SubmitSubmissionService{submissionRepo: submissionRepo, clock: clock, publisher: publisher, topic: topic}
}

func (s *SubmitSubmissionService) Execute(ctx context.Context, input in.SubmitSubmissionInput) error {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return errors.New("submission id is required")
	}

	submitted, err := s.submissionRepo.Submit(input.SubmissionID, s.clock.Now())
	if err != nil {
		return err
	}
	if !submitted {
		return errors.New("submission is already submitted or has expired")
	}

	submission, answers, err := s.submissionRepo.FindByID(input.SubmissionID)
	if err != nil {
		return err
	}
	if err := s.publisher.PublishSubmissionSubmitted(ctx, newSubmissionSubmittedMessage(s.topic, submission, answers)); err != nil {
		return fmt.Errorf("publish submitted submission %s: %w", submission.ID, err)
	}

	return nil
}

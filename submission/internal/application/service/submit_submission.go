package submission

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

type SubmitSubmissionService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
}

func NewSubmitSubmissionService(submissionRepo out.SubmissionRepository, clock out.Clock) in.SubmitSubmissionService {
	return &SubmitSubmissionService{submissionRepo: submissionRepo, clock: clock}
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

	return nil
}

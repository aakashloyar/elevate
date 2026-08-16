package submission

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
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

	submission, _, err := s.submissionRepo.FindByID(input.SubmissionID)
	if err != nil {
		return err
	}
	if !submission.Status.IsSubmittable() {
		return errors.New("submission is already submitted")
	}

	return s.submissionRepo.UpdateSubmissionTime(input.SubmissionID, s.clock.Now(), domain.SubmissionStatusSubmitted)
}

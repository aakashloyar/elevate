package submission

import (
	"context"
	"errors"
	"strings"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type StartSubmissionService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
}

func NewStartSubmissionService(submissionRepo out.SubmissionRepository, clock out.Clock) in.StartSubmissionService {
	return &StartSubmissionService{submissionRepo: submissionRepo, clock: clock}
}

func (s *StartSubmissionService) Execute(ctx context.Context, input in.StartSubmissionInput) error {
	if strings.TrimSpace(input.SubmissionID) == "" {
		return errors.New("submission id is required")
	}

	submission, _, err := s.submissionRepo.FindByID(input.SubmissionID)
	if err != nil {
		return err
	}
	if !submission.Status.IsStartable() {
		return errors.New("submission must be in CREATED state")
	}

	return s.submissionRepo.UpdateStartTime(input.SubmissionID, s.clock.Now(), domain.SubmissionStatusInProgress)
}

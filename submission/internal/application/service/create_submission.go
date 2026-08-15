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

type CreateSubmissionService struct {
	submissionRepo out.SubmissionRepository
	idGen          out.IDGenerator
	clock          out.Clock
}

func NewCreateSubmissionService(submissionRepo out.SubmissionRepository, idGen out.IDGenerator, clock out.Clock) in.CreateSubmissionService {
	return &CreateSubmissionService{submissionRepo: submissionRepo, idGen: idGen, clock: clock}
}

func (s *CreateSubmissionService) Execute(ctx context.Context, input in.CreateSubmissionInput) (in.CreateSubmissionOutput, error) {
	if strings.TrimSpace(input.AssessmentID) == "" {
		return in.CreateSubmissionOutput{}, errors.New("assessment id is required")
	}
	if strings.TrimSpace(input.UserID) == "" {
		return in.CreateSubmissionOutput{}, errors.New("user id is required")
	}

	now := s.clock.Now()
	submission := domain.Submission{
		ID:           s.idGen.NewID(),
		AssessmentID: input.AssessmentID,
		UserID:       input.UserID,
		Status:       "IN_PROGRESS",
		StartedAt:    now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.submissionRepo.Save(submission); err != nil {
		return in.CreateSubmissionOutput{}, err
	}

	return in.CreateSubmissionOutput{SubmissionID: submission.ID, StartedAt: now.Format(time.RFC3339)}, nil
}

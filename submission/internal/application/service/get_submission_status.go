package submission

import (
	"context"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

type GetSubmissionStatusService struct {
	submissionRepo out.SubmissionRepository
}

func NewGetSubmissionStatusService(submissionRepo out.SubmissionRepository) in.GetSubmissionStatusService {
	return &GetSubmissionStatusService{submissionRepo: submissionRepo}
}

func (s *GetSubmissionStatusService) Execute(ctx context.Context, input in.GetSubmissionStatusInput) (in.GetSubmissionStatusOutput, error) {
	status, expiresAt, err := s.submissionRepo.FindStatus(input.SubmissionID)
	if err != nil {
		return in.GetSubmissionStatusOutput{}, err
	}
	return in.GetSubmissionStatusOutput{SubmissionID: input.SubmissionID, Status: status, ExpiresAt: expiresAt}, nil
}

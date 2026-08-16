package submission

import (
	"context"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

const expirationBatchSize = 100

type ExpireSubmissionsService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
}

func NewExpireSubmissionsService(submissionRepo out.SubmissionRepository, clock out.Clock) in.ExpireSubmissionsService {
	return &ExpireSubmissionsService{submissionRepo: submissionRepo, clock: clock}
}

func (s *ExpireSubmissionsService) Execute(ctx context.Context) (int, error) {
	return s.submissionRepo.ExpireSubmissions(s.clock.Now(), expirationBatchSize)
}

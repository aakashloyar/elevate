package out

import (
	"context"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type SubmissionClient interface {
	UpdateSubmissionStatus(ctx context.Context, submissionID string, status domain.SubmissionStatus) error
}

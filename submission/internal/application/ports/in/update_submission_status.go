package submission

import (
	"context"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type UpdateSubmissionStatusInput struct {
	SubmissionID string
	Status       domain.SubmissionStatus
}

type UpdateSubmissionStatusOutput struct {
	SubmissionID string
	Status       domain.SubmissionStatus
}

type UpdateSubmissionStatusService interface {
	Execute(ctx context.Context, input UpdateSubmissionStatusInput) (UpdateSubmissionStatusOutput, error)
}

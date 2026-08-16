package submission

import (
	"context"
	"time"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type GetSubmissionStatusInput struct {
	SubmissionID string
}

type GetSubmissionStatusOutput struct {
	SubmissionID string
	Status       domain.SubmissionStatus
	ExpiresAt    *time.Time
}

type GetSubmissionStatusService interface {
	Execute(context.Context, GetSubmissionStatusInput) (GetSubmissionStatusOutput, error)
}

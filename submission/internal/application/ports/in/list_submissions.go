package submission

import (
	"context"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type ListSubmissionsInput struct {
	UserID string
}

type ListSubmissionsOutput struct {
	Submissions []domain.Submission
}

type ListSubmissionsService interface {
	Execute(context.Context, ListSubmissionsInput) (ListSubmissionsOutput, error)
}

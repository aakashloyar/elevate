package submission

import (
	"context"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type ListSubmissionsInput struct {
	UserID string
	Limit  *int
	Offset *int
}

type ListSubmissionsOutput struct {
	Submissions []ListSubmissionItem
}

type ListSubmissionItem struct {
	Submission      domain.Submission
	UserName        string
	AssessmentTitle string
}

type ListSubmissionsService interface {
	Execute(context.Context, ListSubmissionsInput) (ListSubmissionsOutput, error)
}

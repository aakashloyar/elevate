package out

import "context"

type SubmissionClient interface {
	GetAttemptAssessmentID(ctx context.Context, attemptID string) (string, error)
}

package out

import "context"

type AssessmentClient interface {
	GetAssessmentProblemIDs(ctx context.Context, assessmentID string) ([]string, error)
}

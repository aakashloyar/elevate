package out

import (
	"context"
	"github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/in"
)

type SubmissionGateway interface {
	GetAttemptAssessmentID(context.Context, string) (string, error)
}
type AssessmentGateway interface {
	GetAssessmentProblemIDs(context.Context, string) ([]string, error)
}
type ProblemGateway interface {
	GetProblemByID(context.Context, string) (in.ProblemView, error)
}

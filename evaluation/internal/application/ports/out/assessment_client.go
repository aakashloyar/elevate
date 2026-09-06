package out

import (
	"context"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type AssessmentClient interface {
	GetAssessmentMarkingScheme(ctx context.Context, assessmentID string) (domain.MarkingScheme, error)
	GetAssessmentProblemIDs(ctx context.Context, assessmentID string) ([]string, error)
}

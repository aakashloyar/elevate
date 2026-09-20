package out

import (
	"context"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type AssessmentClient interface {
	GetAssessmentTitle(ctx context.Context, assessmentID string) (string, error)
	GetAssessmentMarkingScheme(ctx context.Context, assessmentID string) (domain.MarkingScheme, error)
	GetAssessmentProblemIDs(ctx context.Context, assessmentID string) ([]string, error)
}

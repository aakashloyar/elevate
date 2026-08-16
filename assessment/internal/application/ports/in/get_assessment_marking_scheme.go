package assessment

import (
	"context"

	"github.com/aakashloyar/elevate/assessment/internal/domain"
)

type GetAssessmentMarkingSchemeInput struct {
	AssessmentID string
}

type GetAssessmentMarkingSchemeService interface {
	Execute(ctx context.Context, input GetAssessmentMarkingSchemeInput) (domain.AssessmentMarkingScheme, error)
}

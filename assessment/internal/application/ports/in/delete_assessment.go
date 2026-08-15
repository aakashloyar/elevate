package assessment

import "context"

type DeleteAssessmentInput struct {
	AssessmentID string
}

type DeleteAssessmentOutput struct{}

type DeleteAssessmentService interface {
	Execute(ctx context.Context, input DeleteAssessmentInput) error
}

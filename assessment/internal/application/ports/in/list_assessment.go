package assessment

import (
	"context"
)

type ListAssessmentsInput struct {
	UserID      *string
	Title       *string
	Description *string
	Limit       *int
	Offset      *int
}

type ListAssessmentsOutput struct {
	Assessments []GetAssessmentOutput
}

type ListAssessmentsService interface {
	Execute(ctx context.Context, input ListAssessmentsInput) (ListAssessmentsOutput, error)
}

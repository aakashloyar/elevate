package assessment

import "context"

type UpsertAssessmentMarkingSchemeInput struct {
	AssessmentID string
	Single       MarksInput
	Multiple     MarksInput
	Numerical    MarksInput
}

type MarksInput struct {
	Correct   float64
	Incorrect float64
	Skipped   float64
}

type UpsertAssessmentMarkingSchemeService interface {
	Execute(ctx context.Context, input UpsertAssessmentMarkingSchemeInput) error
}

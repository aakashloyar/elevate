package assessment

import "context"

type GetAssessmentsBatchInput struct{ AssessmentIDs []string }
type AssessmentSummary struct {
	ID    string
	Title string
}
type GetAssessmentsBatchOutput struct{ Assessments []AssessmentSummary }
type GetAssessmentsBatchService interface {
	Execute(context.Context, GetAssessmentsBatchInput) (GetAssessmentsBatchOutput, error)
}

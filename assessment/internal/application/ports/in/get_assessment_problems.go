package assessment

import "context"

type GetAssessmentProblemsInput struct {
	AssessmentID string
}

type GetAssessmentProblemsOutput struct {
	ProblemIDs []string
}

type GetAssessmentProblemsService interface {
	Execute(context.Context, GetAssessmentProblemsInput) (GetAssessmentProblemsOutput, error)
}

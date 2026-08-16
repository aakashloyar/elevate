package assessment

import "context"

type CreateAssessmentProblemOptionInput struct {
	Text      string
	IsCorrect bool
}

type CreateAssessmentProblemInput struct {
	AssessmentID string
	CreatedBy    string
	Title        string
	Statement    string
	Type         string
	Difficulty   string
	SourceType   string
	Options      []CreateAssessmentProblemOptionInput
	Tags         []string
}

type CreateAssessmentProblemOutput struct {
	ProblemID string
}

type CreateAssessmentProblemService interface {
	Execute(ctx context.Context, input CreateAssessmentProblemInput) (CreateAssessmentProblemOutput, error)
}

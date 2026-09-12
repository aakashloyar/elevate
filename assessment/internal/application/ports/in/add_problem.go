package assessment

import "context"

type AddProblemOptionInput struct {
	Text      string
	IsCorrect bool
}

type AddProblemInput struct {
	AssessmentID string
	CreatedBy    string
	Title        string
	Statement    string
	Type         string
	Difficulty   string
	SourceType   string
	Options      []AddProblemOptionInput
	Tags         []string
}

type AddProblemOutput struct {
	ProblemID string
}

type AddProblemService interface {
	Execute(ctx context.Context, input AddProblemInput) (AddProblemOutput, error)
}

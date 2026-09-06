package in

import "context"

type ProblemType string

const (
	ProblemTypeSingle    ProblemType = "single"
	ProblemTypeMultiple  ProblemType = "multiple"
	ProblemTypeNumerical ProblemType = "numerical"
)

type GetAttemptProblemsInput struct {
	AttemptID string
	Offset    int
	Limit     int
}

type GetAttemptProblemsOutput struct {
	Problems []ProblemView `json:"problems"`
}

type ProblemView struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Statement   string       `json:"statement"`
	Type        ProblemType  `json:"type"`
	Difficulty  string       `json:"difficulty"`
	Options     []OptionView `json:"options"`
	DraftAnswer []string     `json:"draft_answer"`
}

type OptionView struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type GetAttemptProblemsService interface {
	Execute(context.Context, GetAttemptProblemsInput) (GetAttemptProblemsOutput, error)
}

package in

import "context"

type GetAttemptProblemsInput struct {
	AttemptID string
	Offset    int
	Limit     int
}

type GetAttemptProblemsOutput struct {
	Problems []ProblemView `json:"problems"`
}

type ProblemView struct {
	ID         string       `json:"id"`
	Title      string       `json:"title"`
	Statement  string       `json:"statement"`
	Type       string       `json:"type"`
	Difficulty string       `json:"difficulty"`
	Options    []OptionView `json:"options"`
}

type OptionView struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type GetAttemptProblemsService interface {
	Execute(context.Context, GetAttemptProblemsInput) (GetAttemptProblemsOutput, error)
}

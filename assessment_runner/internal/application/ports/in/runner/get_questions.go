package runner

import "context"

type GetQuestionsInput struct {
	SessionID string
	Offset    int
	Limit     int
}

type GetQuestionsOutput struct {
	Questions []QuestionView
}

type QuestionView struct {
	ID         string
	Title      string
	Statement  string
	Type       string
	Difficulty string
	Options    []QuestionOptionView
}

type QuestionOptionView struct {
	ID   string
	Text string
}

type GetQuestionsService interface {
	Execute(ctx context.Context, input GetQuestionsInput) (GetQuestionsOutput, error)
}

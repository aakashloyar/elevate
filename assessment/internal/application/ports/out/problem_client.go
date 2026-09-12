package out

import (
	"context"
	"fmt"
)

type CreateProblemOptionInput struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type CreateProblemInput struct {
	CreatedBy  string                     `json:"created_by"`
	Title      string                     `json:"title"`
	Statement  string                     `json:"statement"`
	Type       string                     `json:"type"`
	Difficulty string                     `json:"difficulty"`
	SourceType string                     `json:"source_type"`
	Options    []CreateProblemOptionInput `json:"options"`
	Tags       []string                   `json:"tags"`
}

type CreateProblemOutput struct {
	ProblemID string `json:"problem_id"`
}

type ProblemClientError struct {
	StatusCode int
	Message    string
}

func (e *ProblemClientError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("problem service request failed with status %d", e.StatusCode)
	}
	return fmt.Sprintf("problem service request failed with status %d: %s", e.StatusCode, e.Message)
}

type ProblemClient interface {
	CreateProblem(ctx context.Context, input CreateProblemInput) (CreateProblemOutput, error)
}

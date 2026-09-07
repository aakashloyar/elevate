package out

import (
	"context"

	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type GeneratedProblemBatchMessage struct {
	Topic string
	Event GeneratedProblemBatchEvent
}

type GeneratedProblemBatchEvent struct {
	AssessmentID string             `json:"assessment_id"`
	Problems     []GeneratedProblem `json:"problems"`
}

type GeneratedProblem struct {
	CreatedBy  string                      `json:"created_by"`
	Title      string                      `json:"title"`
	Statement  string                      `json:"statement"`
	Type       domain.GeneratedProblemType `json:"type"`
	Difficulty domain.GenerationLevel      `json:"difficulty"`
	SourceType string                      `json:"source_type"`
	Options    []GeneratedProblemOption    `json:"options"`
	Tags       []string                    `json:"tags"`
}

type GeneratedProblemOption struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type GeneratedProblemPublisher interface {
	PublishGeneratedProblems(ctx context.Context, message GeneratedProblemBatchMessage) error
	Close()
}

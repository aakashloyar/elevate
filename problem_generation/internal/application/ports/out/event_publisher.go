package out

import (
	"context"
)

type GenerationRequestedEvent struct {
	JobID              string   `json:"job_id"`
	UserID             string   `json:"user_id"`
	AssessmentID       *string  `json:"assessment_id,omitempty"`
	SingleCorrectCount int      `json:"single_correct_count"`
	MultiCorrectCount  int      `json:"multi_correct_count"`
	NumericalCount     int      `json:"numerical_count"`
	DocumentID         *string  `json:"document_id,omitempty"`
	Level              string   `json:"level"`
	Description        string   `json:"description"`
	TopicIDs           []string `json:"topic_ids"`
	Prompt             string   `json:"prompt"`
}

type GenerationRequestedMessage struct {
	Event GenerationRequestedEvent
}

type EventPublisher interface {
	PublishGenerationRequested(ctx context.Context, message GenerationRequestedMessage) error
	Close()
}

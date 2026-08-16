package out

import "context"

type CreatedBatchEvent struct {
	AssessmentID string   `json:"assessment_id"`
	ProblemIDs   []string `json:"problem_ids"`
}

type CreatedBatchMessage struct {
	Topic string
	Event CreatedBatchEvent
}

type EventPublisher interface {
	PublishCreatedBatch(ctx context.Context, message CreatedBatchMessage) error
	Close()
}

package out

import "context"

type GenerationRequestedEvent struct {
	JobID string `json:"job_id"`
}

type GenerationRequestedMessage struct {
	Event GenerationRequestedEvent
}

type EventPublisher interface {
	PublishGenerationRequested(ctx context.Context, message GenerationRequestedMessage) error
	Close()
}

package out

import (
	"context"
)	

type CreatedBatchEvent struct {
	Topic        string
	AssessmentID string
	ProblemIDs   []string
}

type EventPublisher interface {
	PublishCreatedBatch(ctx context.Context, event CreatedBatchEvent) error
	Close()
}
package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aakashloyar/elevate/document/internal/domain"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (p *Producer) PublishDocumentUploaded(ctx context.Context, event domain.DocumentUploadedEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	p.client.Produce(ctx, &kgo.Record{
		Topic: p.topic, Key: []byte(event.DocumentID), Value: value, Timestamp: time.Now().UTC(),
	}, func(_ *kgo.Record, err error) {
		done <- err
	})

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

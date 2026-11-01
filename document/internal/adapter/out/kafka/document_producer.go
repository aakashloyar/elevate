package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aakashloyar/elevate/document/internal/application/ports/out"
	"github.com/aakashloyar/elevate/document/internal/observe"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (p *Producer) PublishDocumentUploaded(ctx context.Context, message out.DocumentUploadedMessage) error {
	defer observe.Dependency("kafka.publish_document_uploaded")()
	value, err := json.Marshal(message.Event)
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	p.client.Produce(ctx, &kgo.Record{
		Topic: p.topic, Key: []byte(message.Event.DocumentID), Value: value, Timestamp: time.Now().UTC(),
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

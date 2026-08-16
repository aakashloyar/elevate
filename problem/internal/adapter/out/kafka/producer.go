package kafkaproducer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
)

func (p *Producer) PublishCreatedBatch(ctx context.Context, message out.CreatedBatchMessage) error {

	bytes, err := json.Marshal(message.Event)
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Topic:     message.Topic,
		Value:     bytes,
		Timestamp: time.Now(),
	}

	done := make(chan error, 1)
	p.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
		done <- err
	})

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Producer) Close() {
	p.client.Close()
}

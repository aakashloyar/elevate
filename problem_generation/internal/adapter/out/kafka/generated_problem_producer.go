package kafkaproducer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
	"github.com/twmb/franz-go/pkg/kgo"
)

func (p *Producer) PublishGeneratedProblems(ctx context.Context, message out.GeneratedProblemBatchMessage) error {
	value, err := json.Marshal(message.Event)
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	p.client.Produce(ctx, &kgo.Record{
		Topic:     message.Topic,
		Key:       []byte(message.Event.AssessmentID),
		Value:     value,
		Timestamp: time.Now().UTC(),
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

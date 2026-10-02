package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

type EmailEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Text      string    `json:"text,omitempty"`
	HTML      string    `json:"html,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Producer) SendEmail(ctx context.Context, to, subject, html string) error {
	event, err := json.Marshal(EmailEvent{
		EventID: uuid.NewString(), EventType: "email.send", To: to,
		Subject: subject, HTML: html, CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	done := make(chan error, 1)
	startedAt := time.Now()
	p.client.Produce(ctx, &kgo.Record{
		Topic: p.topic, Key: []byte(to), Value: event, Timestamp: time.Now().UTC(),
	}, func(_ *kgo.Record, err error) {
		done <- err
	})

	select {
	case err := <-done:
		log.Printf("service=auth kafka_publish topic=%s event_type=email.send duration_ms=%.3f error=%t", p.topic, float64(time.Since(startedAt).Microseconds())/1000, err != nil)
		return err
	case <-ctx.Done():
		log.Printf("service=auth kafka_publish topic=%s event_type=email.send duration_ms=%.3f error=true", p.topic, float64(time.Since(startedAt).Microseconds())/1000)
		return ctx.Err()
	}
}

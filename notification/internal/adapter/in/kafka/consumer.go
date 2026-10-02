package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

type EmailEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Text      string    `json:"text"`
	HTML      string    `json:"html"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Consumer) Start(ctx context.Context) error {
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("notification kafka fetch error: %v", err)
			}
			continue
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var event EmailEvent
			err := json.Unmarshal(record.Value, &event)
			if err != nil || event.EventID == "" || event.To == "" {
				log.Printf("invalid notification event: %v", err)
				if commitErr := c.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit invalid notification event: %v", commitErr)
				}
				continue
			}
			if err := c.sender.Send(event.To, event.Subject, event.Text, event.HTML); err != nil {
				log.Printf("send notification email %s failed: %v; will retry", event.EventID, err)
				select {
				case <-time.After(2 * time.Second):
				case <-ctx.Done():
					return ctx.Err()
				}
				continue
			}
			log.Printf("service=notification kafka_consume event_type=%s event_id=%s status=sent", event.EventType, event.EventID)
			if err := c.client.CommitRecords(ctx, record); err != nil {
				log.Printf("commit notification event: %v", err)
			}
		}
	}
	return ctx.Err()
}

func (c *Consumer) Close() {
	c.client.Close()
}

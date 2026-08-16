package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

func (c *Consumer) Start(ctx context.Context) error {
	for ctx.Err() == nil {
		fetches := c.kafkaClient.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("evaluation kafka fetch error: %v", err)
			}
			continue
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var event domain.SubmissionSubmitted
			if err := json.Unmarshal(record.Value, &event); err != nil {
				log.Printf("invalid submission-submitted event: %v", err)
				continue
			}
			if _, err := c.service.Evaluate(ctx, event); err != nil {
				log.Printf("evaluation failed for submission %s: %v", event.SubmissionID, err)
				continue
			}
			if err := c.kafkaClient.client.CommitRecords(ctx, record); err != nil {
				log.Printf("commit evaluation event: %v", err)
			}
		}
	}
	return ctx.Err()
}
func (c *Consumer) Close() { c.kafkaClient.client.Close() }

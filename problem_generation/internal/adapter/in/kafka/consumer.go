package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out"
)

func (c *Consumer) Start(ctx context.Context) error {
	for ctx.Err() == nil {
		fetches := c.kafkaClient.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("problem generation kafka fetch error: %v", err)
			}
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var event out.GenerationRequestedEvent
			if err := json.Unmarshal(record.Value, &event); err != nil {
				log.Printf("invalid generation-requested event: %v", err)
				if commitErr := c.kafkaClient.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected generation event failed: %v", commitErr)
				}
				continue
			}
			if event.JobID == "" {
				log.Printf("generation-requested event missing job_id; skipping")
				if commitErr := c.kafkaClient.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected generation event failed: %v", commitErr)
				}
				continue
			}

			if err := c.service.Execute(ctx, event.JobID); err != nil {
				log.Printf("generation job failed for job %s: %v", event.JobID, err)
				continue
			}
			if err := c.kafkaClient.client.CommitRecords(ctx, record); err != nil {
				log.Printf("commit generation event: %v", err)
			}
		}
	}
	return ctx.Err()
}

func (c *Consumer) Close() {
	c.kafkaClient.client.Close()
}

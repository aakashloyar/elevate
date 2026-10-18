package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aakashloyar/elevate/document/internal/domain"
)

func (c *Consumer) Start(ctx context.Context) error {
	for ctx.Err() == nil {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("document kafka fetch error: %v", err)
			}
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var event domain.DocumentUploadedEvent
			if err := json.Unmarshal(record.Value, &event); err != nil || event.DocumentID == "" || event.StorageKey == "" {
				log.Printf("invalid document uploaded event: %v", err)
				if commitErr := c.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit invalid document event: %v", commitErr)
				}
				continue
			}

			if err := c.service.Execute(ctx, event); err != nil {
				log.Printf("document processing failed for %s: %v", event.DocumentID, err)
				continue
			}
			if err := c.client.CommitRecords(ctx, record); err != nil {
				log.Printf("commit document event: %v", err)
			}
		}
	}
	return ctx.Err()
}

func (c *Consumer) Close() {
	c.client.Close()
}

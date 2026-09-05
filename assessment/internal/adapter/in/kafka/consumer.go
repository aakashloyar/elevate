package consumer

import (
	"context"
	"encoding/json"
	"log"

	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
)

func (c *Consumer) Start(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		fetches := c.kafkaClient.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("fetch error: %v", err)
			}
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			var batch problemCreatedBatch
			if err := json.Unmarshal(record.Value, &batch); err != nil {
				log.Printf("invalid kafka message payload: %v", err)
				if commitErr := c.kafkaClient.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected message failed: %v", commitErr)
				}
				continue
			}

			if batch.AssessmentID == "" || len(batch.ProblemIDs) == 0 {
				log.Printf("empty or invalid batch: %+v", batch)
				if commitErr := c.kafkaClient.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected message failed: %v", commitErr)
				}
				continue
			}

			if len(batch.ProblemIDs) > c.maxPerBatch {
				log.Printf("batch size %d exceeds max %d; rejecting", len(batch.ProblemIDs), c.maxPerBatch)
				if commitErr := c.kafkaClient.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected message failed: %v", commitErr)
				}
				continue
			}

			if _, err := c.addProblemsSvc.Execute(ctx, in.AddProblemsBatchInput{AssessmentID: batch.AssessmentID, ProblemIDs: batch.ProblemIDs}); err != nil {
				log.Printf("failed to add problems to assessment %s: %v", batch.AssessmentID, err)
				continue
			}

			if err := c.kafkaClient.client.CommitRecords(ctx, record); err != nil {
				log.Printf("commit failed: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.kafkaClient.client.Close()
}

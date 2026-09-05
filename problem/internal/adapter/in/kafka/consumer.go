package consumer

import (
	"context"
	"encoding/json"
	"log"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem/internal/domain"
)

type generationProblemInput struct {
	AssessmentID string                 `json:"assessment_id"`
	Problems     []createProblemPayload `json:"problems"`
}

type createProblemPayload struct {
	CreatedBy  string             `json:"created_by"`
	Title      string             `json:"title"`
	Statement  string             `json:"statement"`
	Type       domain.ProblemType `json:"type"`
	Difficulty domain.Difficulty  `json:"difficulty"`
	SourceType domain.SourceType  `json:"source_type"`
	Options    []Options          `json:"options"`
	Tags       []string           `json:"tags"`
}

type Options struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Printf("fetch error: %v", err)
			}
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			var batch generationProblemInput
			if err := json.Unmarshal(record.Value, &batch); err != nil {
				log.Printf("invalid kafka message payload: %v", err)
				if commitErr := c.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected message failed: %v", commitErr)
				}
				continue
			}

			if batch.AssessmentID == "" {
				log.Printf("message missing assessment_id; skipping")
				if commitErr := c.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected message failed: %v", commitErr)
				}
				continue
			}

			if len(batch.Problems) == 0 {
				log.Printf("empty problem batch for assessment %s; skipping", batch.AssessmentID)
				if commitErr := c.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected message failed: %v", commitErr)
				}
				continue
			}

			if len(batch.Problems) > c.maxPerBatch {
				log.Printf("batch size %d exceeds max %d; rejecting", len(batch.Problems), c.maxPerBatch)
				if commitErr := c.client.CommitRecords(ctx, record); commitErr != nil {
					log.Printf("commit rejected message failed: %v", commitErr)
				}
				continue
			}

			problems := []in.CreateProblemInput{}
			for _, p := range batch.Problems {
				options := make([]in.CreateProblemOptionInput, 0, len(p.Options))
				for _, o := range p.Options {
					options = append(options, in.CreateProblemOptionInput{Text: o.Text, IsCorrect: o.IsCorrect})
				}
				problems = append(problems, in.CreateProblemInput{
					CreatedBy:  p.CreatedBy,
					Title:      p.Title,
					Statement:  p.Statement,
					Type:       p.Type,
					Difficulty: p.Difficulty,
					SourceType: p.SourceType,
					Options:    options,
					Tags:       p.Tags,
				})
			}

			// call the create batch service to persist problems and publish created IDs
			if _, err := c.createSvc.Execute(ctx, in.CreateProblemBatchInput{Problems: problems, AssessmentID: batch.AssessmentID}); err != nil {
				log.Printf("failed to create problem batch for assessment %s: %v", batch.AssessmentID, err)
				continue
			}

			if err := c.client.CommitRecords(ctx, record); err != nil {
				log.Printf("commit failed: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}

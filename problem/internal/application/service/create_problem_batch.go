package problem

import (
	"context"
	"fmt"
	"log"

	"github.com/aakashloyar/elevate/problem/config"
	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out"
)

type CreateProblemBatchService struct {
	createProblemSvc in.CreateProblemService
	producer         out.EventPublisher
	idGen            out.IDGenerator
	clock            out.Clock
}

func NewCreateProblemBatchService(createProblemSvc in.CreateProblemService, producer out.EventPublisher, idGen out.IDGenerator, clock out.Clock) in.CreateProblemBatchService {
	return &CreateProblemBatchService{createProblemSvc: createProblemSvc, producer: producer, idGen: idGen, clock: clock}
}

func (s *CreateProblemBatchService) Execute(ctx context.Context, input in.CreateProblemBatchInput) (in.CreateProblemBatchOutput, error) {
	createdIDs := []string{}
	for _, problemInput := range input.Problems {
		out, err := s.createProblemSvc.Execute(ctx, in.CreateProblemInput{
			CreatedBy:  problemInput.CreatedBy,
			Title:      problemInput.Title,
			Statement:  problemInput.Statement,
			Type:       problemInput.Type,
			Difficulty: problemInput.Difficulty,
			SourceType: problemInput.SourceType,
			Options:    problemInput.Options,
			Tags:       problemInput.Tags,
		})
		if err != nil {
			log.Printf("failed to create problem for assessment %s: %v", input.AssessmentID, err)
			continue
		}
		createdIDs = append(createdIDs, out.ProblemID)
	}
	if len(createdIDs) == 0 {
		return in.CreateProblemBatchOutput{}, nil
	}
	if input.AssessmentID != "" {
		if err := s.producer.PublishCreatedBatch(ctx, out.CreatedBatchMessage{
			Topic: config.CreatedProblemBatchTopic,
			Event: out.CreatedBatchEvent{
				AssessmentID: input.AssessmentID,
				ProblemIDs:   createdIDs,
			},
		}); err != nil {
			return in.CreateProblemBatchOutput{}, fmt.Errorf("publish created problem batch: %w", err)
		}
	}
	return in.CreateProblemBatchOutput{}, nil
}

package submission

import (
	"context"
	"log"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out"
)

const expirationBatchSize = 100

type ExpireSubmissionsService struct {
	submissionRepo out.SubmissionRepository
	clock          out.Clock
	publisher      out.EventPublisher
	topic          string
}

func NewExpireSubmissionsService(submissionRepo out.SubmissionRepository, clock out.Clock, publisher out.EventPublisher, topic string) in.ExpireSubmissionsService {
	return &ExpireSubmissionsService{submissionRepo: submissionRepo, clock: clock, publisher: publisher, topic: topic}
}

func (s *ExpireSubmissionsService) Execute(ctx context.Context) (int, error) {
	submissionIDs, err := s.submissionRepo.ExpireSubmissions(s.clock.Now(), expirationBatchSize)
	if err != nil {
		return 0, err
	}

	for _, submissionID := range submissionIDs {
		submission, answers, err := s.submissionRepo.FindByID(submissionID)
		if err != nil {
			log.Printf("failed to load expired submission %s for publishing: %v", submissionID, err)
			continue
		}
		if err := s.publisher.PublishSubmissionSubmitted(ctx, newSubmissionSubmittedMessage(s.topic, submission, answers)); err != nil {
			log.Printf("failed to publish expired submission %s: %v", submission.ID, err)
		}
	}

	return len(submissionIDs), nil
}

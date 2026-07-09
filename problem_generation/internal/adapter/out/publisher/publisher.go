package publisher

import (
	"context"
	"log"

	"github.com/aakashloyar/elevate/problem_generation/internal/domain"
)

type Publisher struct {
	topic string
}

func NewPublisher(topic string) *Publisher {
	return &Publisher{topic: topic}
}

func (p *Publisher) PublishGenerationRequested(ctx context.Context, job domain.GenerationJob, prompt string) error {
	log.Printf("generation request topic=%s job_id=%s user_id=%s prompt=%q", p.topic, job.ID, job.UserID, prompt)
	return nil
}

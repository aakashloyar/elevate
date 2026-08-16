package out

import (
	"context"
	"time"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type SubmissionSubmittedEvent struct {
	SubmissionID    string                  `json:"submission_id"`
	AssessmentID    string                  `json:"assessment_id"`
	UserID          string                  `json:"user_id"`
	Status          domain.SubmissionStatus `json:"status"`
	StartedAt       time.Time               `json:"started_at"`
	DurationSeconds int                     `json:"duration_seconds"`
	ExpiresAt       *time.Time              `json:"expires_at,omitempty"`
	SubmittedAt     *time.Time              `json:"submitted_at,omitempty"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	Answers         []SubmissionAnswerEvent `json:"answers"`
}

type SubmissionSubmittedMessage struct {
	Topic string
	Event SubmissionSubmittedEvent
}

type SubmissionAnswerEvent struct {
	ProblemID string    `json:"problem_id"`
	Answer    []string  `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventPublisher interface {
	PublishSubmissionSubmitted(ctx context.Context, message SubmissionSubmittedMessage) error
	Close()
}

package submission

import (
	"context"
	"time"
)

type GetSubmissionInput struct {
	SubmissionID string
}

type GetSubmissionOutput struct {
	ID           string
	AssessmentID string
	UserID       string
	Status       string
	StartedAt    time.Time
	SubmittedAt  *time.Time
	Answers      []SubmissionAnswerOutput
}

type SubmissionAnswerOutput struct {
	ProblemID  string
	Answer     []string
	AnsweredAt time.Time
}

type GetSubmissionService interface {
	Execute(ctx context.Context, input GetSubmissionInput) (GetSubmissionOutput, error)
}

package submission

import (
	"context"
	"time"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type GetSubmissionInput struct {
	SubmissionID string
}

type GetSubmissionOutput struct {
	ID           string
	AssessmentID string
	UserID       string
	Status       domain.SubmissionStatus
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

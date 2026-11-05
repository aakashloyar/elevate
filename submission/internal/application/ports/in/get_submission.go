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
	ID              string
	AssessmentID    string
	UserID          string
	Status          domain.SubmissionStatus
	DurationSeconds int
	StartedAt       *time.Time
	ExpiresAt       *time.Time
	SubmittedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	MarkingScheme   *domain.MarkingScheme
	Problems        []SubmissionProblemOutput
}

type SubmissionProblemOutput struct {
	ProblemID        string
	ProblemType      domain.ProblemType
	OptionIDs        []string
	OptionTexts      []string
	Answer           []string
	TimeSpentSeconds int
	AnswerUpdatedAt  *time.Time
}

type GetSubmissionService interface {
	Execute(ctx context.Context, input GetSubmissionInput) (GetSubmissionOutput, error)
}

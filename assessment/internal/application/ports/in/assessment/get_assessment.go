package assessment

import (
	"context"
	"time"
)

type GetAssessmentInput struct {
	AssessmentID string
}

type GetAssessmentOutput struct {
	ID              string
	Title           string
	Description     string
	Status          string
	DurationSeconds int
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type GetAssessmentService interface {
	Execute(ctx context.Context, input GetAssessmentInput) (GetAssessmentOutput, error)
}

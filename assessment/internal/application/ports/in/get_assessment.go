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
	DurationSeconds int
	CreatedBy       string
	CreatedByName   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DocumentIDs     []string
}

type GetAssessmentService interface {
	Execute(ctx context.Context, input GetAssessmentInput) (GetAssessmentOutput, error)
}

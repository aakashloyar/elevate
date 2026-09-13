package out

import (
	"context"

	"github.com/aakashloyar/elevate/submission/internal/domain"
)

type ProblemSnapshot struct {
	ProblemID   string
	ProblemType domain.ProblemType
	Options     []Option
}

type Option struct {
	ID   string
	Text string
}

type ProblemClient interface {
	GetProblemSnapshots(ctx context.Context, problemIDs []string) ([]ProblemSnapshot, error)
}

type AssessmentClient interface {
	GetProblemIDs(ctx context.Context, assessmentID string) ([]string, error)
}

type UserClient interface {
	Exists(ctx context.Context, userID string) error
}

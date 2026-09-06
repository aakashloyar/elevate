package out

import (
	"context"

	in "github.com/aakashloyar/elevate/assessment_runner/internal/application/ports/in"
)

type ProblemClient interface {
	GetProblemByID(ctx context.Context, problemID string) (in.ProblemView, error)
}

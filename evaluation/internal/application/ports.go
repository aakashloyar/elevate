package application

import (
	"context"
	"github.com/aakashloyar/elevate/evaluation/internal/domain"
)

type AssessmentClient interface {
	MarkingScheme(context.Context, string) (domain.MarkingScheme, error)
	ProblemIDs(context.Context, string) ([]string, error)
}
type ProblemClient interface {
	Problem(context.Context, string) (domain.Problem, error)
}
type SubmissionClient interface {
	UpdateStatus(context.Context, string, domain.SubmissionStatus) error
}
type Repository interface {
	Save(context.Context, domain.Evaluation) error
	FindBySubmissionID(context.Context, string) (domain.Evaluation, error)
}

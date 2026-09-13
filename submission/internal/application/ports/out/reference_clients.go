package out

import "context"

type AssessmentClient interface {
	Exists(ctx context.Context, assessmentID string) error
}

type UserClient interface {
	Exists(ctx context.Context, userID string) error
}

type ProblemClient interface {
	Exists(ctx context.Context, problemID string) error
}

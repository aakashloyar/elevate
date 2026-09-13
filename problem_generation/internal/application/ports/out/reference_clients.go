package out

import "context"

type UserClient interface {
	Exists(ctx context.Context, userID string) error
}

type AssessmentClient interface {
	Exists(ctx context.Context, assessmentID string) error
}

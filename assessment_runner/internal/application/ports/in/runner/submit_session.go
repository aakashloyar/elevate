package runner

import "context"

type SubmitSessionInput struct {
	SessionID string
}

type SubmitSessionOutput struct{}

type SubmitSessionService interface {
	Execute(ctx context.Context, input SubmitSessionInput) error
}

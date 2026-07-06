package runner

import "context"

type GetSessionInput struct {
	SessionID string
}

type GetSessionOutput struct {
	SessionID      string
	RemainingTime  int
	TotalQuestions int
	Status         string
}

type GetSessionService interface {
	Execute(ctx context.Context, input GetSessionInput) (GetSessionOutput, error)
}

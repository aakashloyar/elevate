package runner

import "context"

type StartSessionInput struct {
	AssessmentID string
	UserID       string
}

type StartSessionOutput struct {
	SessionID      string
	SubmissionID   string
	RemainingTime  int
	TotalQuestions int
}

type StartSessionService interface {
	Execute(ctx context.Context, input StartSessionInput) (StartSessionOutput, error)
}

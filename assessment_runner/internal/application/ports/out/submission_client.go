package out

import "context"

type AttemptDraft struct {
	AssessmentID string
	Answers      []DraftAnswer
}

type DraftAnswer struct {
	ProblemID string
	Answer    []string
}

type SubmissionClient interface {
	GetAttemptDraft(ctx context.Context, attemptID string) (AttemptDraft, error)
}

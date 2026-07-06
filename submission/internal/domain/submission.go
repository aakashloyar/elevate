package domain

import "time"

type Submission struct {
	ID           string
	AssessmentID string
	UserID       string
	Status       string
	StartedAt    time.Time
	SubmittedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SubmissionAnswer struct {
	ID           string
	SubmissionID string
	ProblemID    string
	Answer       []string
	AnsweredAt   time.Time
}

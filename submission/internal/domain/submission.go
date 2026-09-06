package domain

import "time"

type SubmissionStatus string

const (
	SubmissionStatusCreated          SubmissionStatus = "CREATED"
	SubmissionStatusInProgress       SubmissionStatus = "IN_PROGRESS"
	SubmissionStatusSubmitted        SubmissionStatus = "SUBMITTED"
	SubmissionStatusUnderEvaluation  SubmissionStatus = "UNDER_EVALUATION"
	SubmissionStatusEvaluated        SubmissionStatus = "EVALUATED"
	SubmissionStatusEvaluationFailed SubmissionStatus = "EVALUATION_FAILED"
)

func (s SubmissionStatus) IsStartable() bool {
	return s == SubmissionStatusCreated
}

func (s SubmissionStatus) IsSubmittable() bool {
	return s == SubmissionStatusInProgress
}

type Submission struct {
	ID              string
	AssessmentID    string
	UserID          string
	Status          SubmissionStatus
	StartedAt       time.Time
	DurationSeconds int
	ExpiresAt       *time.Time
	SubmittedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type SubmissionAnswer struct {
	ID           string
	SubmissionID string
	ProblemID    string
	Answer       []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

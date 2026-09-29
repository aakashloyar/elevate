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
	StartedAt       *time.Time
	DurationSeconds int
	ExpiresAt       *time.Time
	SubmittedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	MarkingScheme   *MarkingScheme
}

type MarkingScheme struct {
	AssessmentID string `json:"assessment_id"`
	Single       Marks  `json:"single"`
	Multiple     Marks  `json:"multiple"`
	Numerical    Marks  `json:"numerical"`
}

type Marks struct {
	Correct   float64 `json:"correct"`
	Incorrect float64 `json:"incorrect"`
	Skipped   float64 `json:"skipped"`
}

type SubmissionAnswerDraft struct {
	SubmissionID    string
	ProblemID       string
	ProblemType     ProblemType
	Options         []Option
	Answer          []string
	AnswerUpdatedAt *time.Time
}

type Option struct {
	ID   string
	Text string
}

type ProblemType string

const (
	ProblemTypeSingle    ProblemType = "single"
	ProblemTypeMultiple  ProblemType = "multiple"
	ProblemTypeNumerical ProblemType = "numerical"
)

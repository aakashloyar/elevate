package domain

import "time"

type ProblemType string
type SubmissionStatus string
type QuestionResultStatus string

const (
	ProblemTypeSingle    ProblemType = "single"
	ProblemTypeMultiple  ProblemType = "multiple"
	ProblemTypeNumerical ProblemType = "numerical"

	SubmissionStatusUnderEvaluation  SubmissionStatus = "UNDER_EVALUATION"
	SubmissionStatusEvaluated        SubmissionStatus = "EVALUATED"
	SubmissionStatusEvaluationFailed SubmissionStatus = "EVALUATION_FAILED"

	QuestionResultStatusCorrect          QuestionResultStatus = "correct"
	QuestionResultStatusPartiallyCorrect QuestionResultStatus = "partially_correct"
	QuestionResultStatusIncorrect        QuestionResultStatus = "incorrect"
	QuestionResultStatusSkipped          QuestionResultStatus = "skipped"
)

// submission submitted event
type SubmissionSubmitted struct {
	SubmissionID    string     `json:"submission_id"`
	AssessmentID    string     `json:"assessment_id"`
	UserID          string     `json:"user_id"`
	StartedAt       time.Time  `json:"started_at"`
	DurationSeconds int        `json:"duration_seconds"`
	SubmittedAt     *time.Time `json:"submitted_at,omitempty"`
	Answers         []Answer   `json:"answers"`
}
type Answer struct {
	ProblemID string    `json:"problem_id"`
	Answer    []string  `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Problem struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Statement  string      `json:"statement"`
	Type       ProblemType `json:"type"`
	Difficulty string      `json:"difficulty"`
	Tags       []string    `json:"tags"`
	Options    []Option    `json:"options"`
}
type Option struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// evaluation result
type Evaluation struct {
	SubmissionID    string           `json:"submission_id"`
	AssessmentID    string           `json:"assessment_id"`
	AssessmentTitle string           `json:"assessment_title"`
	UserID          string           `json:"user_id"`
	UserName        string           `json:"user_name"`
	StartedAt       time.Time        `json:"started_at"`
	DurationSeconds int              `json:"duration_seconds"`
	SubmittedAt     *time.Time       `json:"submitted_at,omitempty"`
	Score           float64          `json:"score"`
	TotalMarks      float64          `json:"total_marks"`
	Questions       []QuestionResult `json:"questions"`
	EvaluatedAt     time.Time        `json:"evaluated_at"`
}

type QuestionResult struct {
	ProblemID       string               `json:"problem_id"`
	Title           string               `json:"title"`
	Statement       string               `json:"statement"`
	Type            ProblemType          `json:"type"`
	Difficulty      string               `json:"difficulty"`
	Tags            []string             `json:"tags"`
	Status          QuestionResultStatus `json:"status"`
	Marks           float64              `json:"marks"`
	Options         []SelectedOption     `json:"options"`
	SelectedOptions []SelectedOption     `json:"selected_options"`
}

type SelectedOption struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// marking scheme for an assessment
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

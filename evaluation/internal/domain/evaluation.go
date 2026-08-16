package domain

import "time"

type ProblemType string

const (
	ProblemTypeSingle    ProblemType = "single"
	ProblemTypeMultiple  ProblemType = "multiple"
	ProblemTypeNumerical ProblemType = "numerical"
)

type SubmissionSubmitted struct {
	SubmissionID string   `json:"submission_id"`
	AssessmentID string   `json:"assessment_id"`
	UserID       string   `json:"user_id"`
	Status       string   `json:"status"`
	Answers      []Answer `json:"answers"`
}
type Answer struct {
	ProblemID string   `json:"problem_id"`
	Answer    []string `json:"answer"`
}
type Problem struct {
	ID      string      `json:"id"`
	Type    ProblemType `json:"type"`
	Options []Option    `json:"options"`
}
type Option struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}
type MarkingScheme struct {
	AssessmentID            string  `json:"assessment_id"`
	SingleCorrectMarks      float64 `json:"single_correct_marks"`
	SingleIncorrectMarks    float64 `json:"single_incorrect_marks"`
	SingleSkippedMarks      float64 `json:"single_skipped_marks"`
	MultipleCorrectMarks    float64 `json:"multiple_correct_marks"`
	MultipleIncorrectMarks  float64 `json:"multiple_incorrect_marks"`
	MultipleSkippedMarks    float64 `json:"multiple_skipped_marks"`
	NumericalCorrectMarks   float64 `json:"numerical_correct_marks"`
	NumericalIncorrectMarks float64 `json:"numerical_incorrect_marks"`
	NumericalSkippedMarks   float64 `json:"numerical_skipped_marks"`
}
type QuestionResult struct {
	ProblemID string      `json:"problem_id"`
	Type      ProblemType `json:"type"`
	Status    string      `json:"status"`
	Marks     float64     `json:"marks"`
}
type Evaluation struct {
	SubmissionID string           `json:"submission_id"`
	AssessmentID string           `json:"assessment_id"`
	UserID       string           `json:"user_id"`
	Score        float64          `json:"score"`
	Questions    []QuestionResult `json:"questions"`
	EvaluatedAt  time.Time        `json:"evaluated_at"`
}

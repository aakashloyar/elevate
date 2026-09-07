package domain

import "time"

type GenerationJobStatus string
type GenerationLevel string
type GeneratedProblemType string

const (
	GenerationJobStatusPending    GenerationJobStatus = "pending"
	GenerationJobStatusProcessing GenerationJobStatus = "processing"
	GenerationJobStatusCompleted  GenerationJobStatus = "completed"
	GenerationJobStatusFailed     GenerationJobStatus = "failed"

	GenerationLevelEasy   GenerationLevel = "easy"
	GenerationLevelMedium GenerationLevel = "medium"
	GenerationLevelHard   GenerationLevel = "hard"

	GeneratedProblemTypeSingle    GeneratedProblemType = "single"
	GeneratedProblemTypeMultiple  GeneratedProblemType = "multiple"
	GeneratedProblemTypeNumerical GeneratedProblemType = "numerical"
)

func (l GenerationLevel) IsValid() bool {
	return l == GenerationLevelEasy || l == GenerationLevelMedium || l == GenerationLevelHard
}

func (t GeneratedProblemType) IsValid() bool {
	return t == GeneratedProblemTypeSingle || t == GeneratedProblemTypeMultiple || t == GeneratedProblemTypeNumerical
}

type GenerationJob struct {
	ID                 string
	UserID             string
	SingleCorrectCount int
	MultiCorrectCount  int
	NumericalCount     int
	DocumentID         *string
	AssessmentID       *string
	Level              GenerationLevel
	Description        string
	Status             GenerationJobStatus
	TopicIDs           []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Topic struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

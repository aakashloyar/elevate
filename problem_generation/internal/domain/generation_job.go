package domain

import "time"

type GenerationJobStatus string
type GenerationLevel string

const (
	GenerationJobStatusPending    GenerationJobStatus = "pending"
	GenerationJobStatusProcessing GenerationJobStatus = "processing"
	GenerationJobStatusCompleted  GenerationJobStatus = "completed"
	GenerationJobStatusFailed     GenerationJobStatus = "failed"

	GenerationLevelEasy   GenerationLevel = "easy"
	GenerationLevelMedium GenerationLevel = "medium"
	GenerationLevelHard   GenerationLevel = "hard"
)

func (l GenerationLevel) IsValid() bool {
	return l == GenerationLevelEasy || l == GenerationLevelMedium || l == GenerationLevelHard
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

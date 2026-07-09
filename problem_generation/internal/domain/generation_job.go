package domain

import "time"

type GenerationJob struct {
	ID                 string
	UserID             string
	SingleCorrectCount int
	MultiCorrectCount  int
	NumericalCount     int
	DocumentID         *string
	AssessmentID       *string
	Level              string
	Description        string
	Status             string
	TopicIDs           []string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Topic struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

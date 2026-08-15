package domain

import (
	"time"
)

type Problem struct {
	ID         string
	Title      string
	Statement  string
	Type       ProblemType
	Difficulty Difficulty
	SourceType SourceType
	CreatedBy  string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ProblemOption struct {
	ID        string
	ProblemID string
	Text      string
	IsCorrect bool
}

type ProblemTag struct {
	ProblemID string
	Tag       string
}

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type ProblemType string

const (
	ProblemTypeSingle    ProblemType = "single"
	ProblemTypeMultiple  ProblemType = "multiple"
	ProblemTypeNumerical ProblemType = "numerical"
)

type SourceType string

const (
	SourceTypeManual SourceType = "manual"
	SourceTypeAI     SourceType = "ai"
)

package domain

import "time"

type Problem struct {
	ID         string
	CreatedBy  string
	Title      string
	Statement  string
	Type       string
	Difficulty string
	SourceType string
	Status     string
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

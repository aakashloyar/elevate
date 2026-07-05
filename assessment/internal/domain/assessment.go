package domain

import "time"

type Assessment struct {
	ID              string
	Title           string
	Description     string
	Status          string
	DurationSeconds int
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

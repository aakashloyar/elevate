package domain

type AssessmentMarkingScheme struct {
	AssessmentID string

	SingleCorrectMarks   float64
	SingleIncorrectMarks float64
	SingleSkippedMarks   float64

	MultipleCorrectMarks   float64
	MultipleIncorrectMarks float64
	MultipleSkippedMarks   float64

	NumericalCorrectMarks   float64
	NumericalIncorrectMarks float64
	NumericalSkippedMarks   float64
}

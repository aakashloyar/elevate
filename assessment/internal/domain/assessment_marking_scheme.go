package domain

type AssessmentMarkingScheme struct {
	AssessmentID string
	Single       Marks
	Multiple     Marks
	Numerical    Marks
}

type Marks struct {
	Correct   float64
	Incorrect float64
	Skipped   float64
}

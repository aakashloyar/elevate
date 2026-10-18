package domain

type IDGenerator interface {
	NewID() string
}

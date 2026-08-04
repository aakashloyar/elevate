package in

import "context"

type DeleteUserInput struct {
	UserID string
}

type DeleteUserOutput struct{}

type DeleteUserService interface {
	Execute(ctx context.Context, input DeleteUserInput) (DeleteUserOutput, error)
}

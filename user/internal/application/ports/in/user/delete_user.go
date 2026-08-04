package in

import "context"

type DeleteUserInput struct {
	UserID string
}

type DeleteUserOutput struct {
	Deleted bool
}

type DeleteUserService interface {
	Execute(ctx context.Context, input DeleteUserInput) (DeleteUserOutput, error)
}

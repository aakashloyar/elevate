package in

import "context"

type CreateUserInput struct {
	Username string
	Email    string
}

type CreateUserOutput struct {
	UserID string
}

type CreateUserService interface {
	Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error)
}

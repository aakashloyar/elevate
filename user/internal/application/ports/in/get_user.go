package in

import (
	"context"
	"time"
)

type GetUserInput struct {
	UserID string
}

type GetUserOutput struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
}

type GetUserService interface {
	Execute(ctx context.Context, input GetUserInput) (GetUserOutput, error)
}

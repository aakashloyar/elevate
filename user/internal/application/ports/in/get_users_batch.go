package in

import "context"

type GetUsersBatchInput struct {
	UserIDs []string
}

type UserSummary struct {
	ID       string
	Username string
}

type GetUsersBatchOutput struct {
	Users []UserSummary
}

type GetUsersBatchService interface {
	Execute(context.Context, GetUsersBatchInput) (GetUsersBatchOutput, error)
}

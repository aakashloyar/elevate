package out

import "context"

type UserClient interface {
	GetUserName(ctx context.Context, userID string) (string, error)
}

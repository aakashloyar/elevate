package submission

import "context"

type ExpireSubmissionsService interface {
	Execute(ctx context.Context) (int, error)
}

package in

import "context"

type ProcessGenerationJobService interface {
	Execute(ctx context.Context, jobID string) error
}

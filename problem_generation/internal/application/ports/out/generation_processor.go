package out

import "context"

type GenerationProcessor interface {
	ProcessGeneration(ctx context.Context, event GenerationRequestedEvent) error
}

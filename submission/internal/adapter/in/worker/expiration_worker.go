package worker

import (
	"context"
	"log"
	"time"

	in "github.com/aakashloyar/elevate/submission/internal/application/ports/in"
)

const expirationCheckInterval = time.Second

type ExpirationWorker struct {
	expireSubmissionsService in.ExpireSubmissionsService
}

func NewExpirationWorker(expireSubmissionsService in.ExpireSubmissionsService) *ExpirationWorker {
	return &ExpirationWorker{expireSubmissionsService: expireSubmissionsService}
}

func (w *ExpirationWorker) Start(ctx context.Context) {
	w.expire(ctx)

	ticker := time.NewTicker(expirationCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.expire(ctx)
		}
	}
}

func (w *ExpirationWorker) expire(ctx context.Context) {
	count, err := w.expireSubmissionsService.Execute(ctx)
	if err != nil {
		log.Printf("failed to expire submissions: %v", err)
		return
	}
	if count > 0 {
		log.Printf("expired %d submission(s)", count)
	}
}

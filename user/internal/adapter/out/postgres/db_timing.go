package postgres

import (
	"log"
	"time"
)

func observeDB(operation string) func() {
	startedAt := time.Now()
	return func() {
		log.Printf("service=user db_operation=%s duration_ms=%.3f", operation, float64(time.Since(startedAt).Microseconds())/1000)
	}
}

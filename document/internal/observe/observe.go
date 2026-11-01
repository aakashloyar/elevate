package observe

import (
	"log"
	"time"
)

func Dependency(operation string) func() {
	startedAt := time.Now()
	return func() {
		log.Printf("service=document dependency=%s duration_ms=%.3f", operation, float64(time.Since(startedAt).Microseconds())/1000)
	}
}

package httpobserve

import (
	"log"
	"net/http"
	"time"
)

func Middleware(service string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startedAt := time.Now()
		statusWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(statusWriter, r)
		log.Printf("service=%s http_request method=%s path=%s status=%d duration_ms=%.3f", service, r.Method, r.URL.Path, statusWriter.status, float64(time.Since(startedAt).Microseconds())/1000)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

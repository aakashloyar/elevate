package httpobserve

import (
	"log"
	"net/http"
	"time"
)

func Do(client *http.Client, request *http.Request, dependency string) (*http.Response, error) {
	startedAt := time.Now()
	response, err := client.Do(request)
	status := 0
	if response != nil {
		status = response.StatusCode
	}
	log.Printf("service=auth http_dependency=%s method=%s path=%s status=%d duration_ms=%.3f error=%t", dependency, request.Method, request.URL.Path, status, float64(time.Since(startedAt).Microseconds())/1000, err != nil)
	return response, err
}

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

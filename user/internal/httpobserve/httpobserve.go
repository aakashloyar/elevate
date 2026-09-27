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
	log.Printf("service=user http_dependency=%s method=%s path=%s status=%d duration_ms=%.3f error=%t", dependency, request.Method, request.URL.Path, status, float64(time.Since(startedAt).Microseconds())/1000, err != nil)
	return response, err
}

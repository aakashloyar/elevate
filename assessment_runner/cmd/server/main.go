package main

import (
	"log"
	"net/http"

	"github.com/aakashloyar/elevate/assessment_runner/config"
	httpapi "github.com/aakashloyar/elevate/assessment_runner/internal/adapter/in/http"
	assessmenthttp "github.com/aakashloyar/elevate/assessment_runner/internal/adapter/out/assessmenthttp"
	problemhttp "github.com/aakashloyar/elevate/assessment_runner/internal/adapter/out/problemhttp"
	submissionhttp "github.com/aakashloyar/elevate/assessment_runner/internal/adapter/out/submissionhttp"
	"github.com/aakashloyar/elevate/assessment_runner/internal/application/service"
)

func main() {
	cfg := config.Load()
	service := service.NewGetAttemptProblemsService(submissionhttp.NewClient(cfg.SubmissionServiceURL), assessmenthttp.NewClient(cfg.AssessmentServiceURL), problemhttp.NewClient(cfg.ProblemServiceURL))
	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, httpapi.NewHandler(service))
	log.Printf("assessment runner starting on :%s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, withCORS(mux)))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

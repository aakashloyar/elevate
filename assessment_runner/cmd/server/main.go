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
	httpapi.NewHandler(service).Register(mux)
	log.Printf("assessment runner starting on :%s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, mux))
}

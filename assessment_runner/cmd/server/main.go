package main

import (
	"log"
	"net/http"

	"github.com/aakashloyar/elevate/assessment_runner/config"
	httpRunner "github.com/aakashloyar/elevate/assessment_runner/internal/adapter/in/http/runner"
	client "github.com/aakashloyar/elevate/assessment_runner/internal/adapter/out/http"
	service "github.com/aakashloyar/elevate/assessment_runner/internal/application/service/runner"
)

func main() {
	clientFactory := client.NewClientFactory(config.App.Services)
	orchestrator := service.NewRunnerService(clientFactory)
	handler := httpRunner.NewHandler(orchestrator)

	mux := http.NewServeMux()
	httpRunner.RegisterRoutes(mux, handler)

	log.Printf("assessment runner starting on :%s", config.App.Server.Port)
	if err := http.ListenAndServe(":"+config.App.Server.Port, mux); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/submission/config"
	httpsubmission "github.com/aakashloyar/elevate/submission/internal/adapter/in/http/submission"
	postgres "github.com/aakashloyar/elevate/submission/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out/system"
	submissionservice "github.com/aakashloyar/elevate/submission/internal/application/service/submission"
)

func main() {
	port, err := strconv.Atoi(config.App.Postgres.Port)
	if err != nil {
		log.Fatalf("invalid POSTGRES_PORT: %v", err)
	}

	dbConfig := postgres.Config{
		Host:     config.App.Postgres.Host,
		Port:     port,
		User:     config.App.Postgres.User,
		Password: config.App.Postgres.Password,
		DBName:   config.App.Postgres.DBName,
		SSLMode:  config.App.Postgres.SSLMode,
	}

	db, err := dbConfig.NewDB()
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	submissionRepo := postgres.NewSubmissionRepository(db)
	if err := submissionRepo.(*postgres.SubmissionRepository).Migrate(); err != nil {
		log.Fatalf("failed to migrate submission tables: %v", err)
	}

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	createSubmissionService := submissionservice.NewCreateSubmissionService(submissionRepo, idGen, clock)
	submitAnswerService := submissionservice.NewSubmitAnswerService(submissionRepo, clock)
	getSubmissionService := submissionservice.NewGetSubmissionService(submissionRepo)

	handler := httpsubmission.NewHandler(createSubmissionService, submitAnswerService, getSubmissionService)

	mux := http.NewServeMux()
	httpsubmission.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("submission service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

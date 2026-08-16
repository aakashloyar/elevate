package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/submission/config"
	httpsubmission "github.com/aakashloyar/elevate/submission/internal/adapter/in/http"
	postgres "github.com/aakashloyar/elevate/submission/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/submission/internal/application/ports/out/system"
	submissionservice "github.com/aakashloyar/elevate/submission/internal/application/service"
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
	startSubmissionService := submissionservice.NewStartSubmissionService(submissionRepo, clock)
	saveAnswerService := submissionservice.NewSaveAnswerService(submissionRepo, clock)
	saveAnswerBatchService := submissionservice.NewSaveAnswerBatchService(saveAnswerService)
	getSubmissionService := submissionservice.NewGetSubmissionService(submissionRepo)
	submitSubmissionService := submissionservice.NewSubmitSubmissionService(submissionRepo, clock)

	handler := httpsubmission.NewHandler(createSubmissionService, startSubmissionService, saveAnswerService, saveAnswerBatchService, getSubmissionService, submitSubmissionService)

	mux := http.NewServeMux()
	httpsubmission.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("submission service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

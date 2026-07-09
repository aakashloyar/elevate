package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/problem_generation/config"
	httpgenerationjob "github.com/aakashloyar/elevate/problem_generation/internal/adapter/in/http/generation_job"
	postgres "github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/publisher"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out/system"
	generationjobsvc "github.com/aakashloyar/elevate/problem_generation/internal/application/service/generation_job"
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

	jobRepo := postgres.NewGenerationJobRepository(db)
	if err := jobRepo.(*postgres.GenerationJobRepository).Migrate(); err != nil {
		log.Fatalf("failed to migrate generation jobs tables: %v", err)
	}

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}
	eventPublisher := publisher.NewPublisher(config.App.Kafka.Topic)

	createJobService := generationjobsvc.NewCreateGenerationJobService(jobRepo, eventPublisher, idGen, clock)
	getJobService := generationjobsvc.NewGetGenerationJobService(jobRepo)

	handler := httpgenerationjob.NewHandler(createJobService, getJobService)

	mux := http.NewServeMux()
	httpgenerationjob.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("problem_generation service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

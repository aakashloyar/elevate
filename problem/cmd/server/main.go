package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/problem/config"
	httpproblem "github.com/aakashloyar/elevate/problem/internal/adapter/in/http/problem"
	postgres "github.com/aakashloyar/elevate/problem/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out/system"
	problemservice "github.com/aakashloyar/elevate/problem/internal/application/service/problem"
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

	problemRepo := postgres.NewProblemRepository(db)
	if err := problemRepo.(*postgres.ProblemRepository).Migrate(); err != nil {
		log.Fatalf("failed to migrate problem tables: %v", err)
	}

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	createProblemService := problemservice.NewCreateProblemService(problemRepo, idGen, clock)
	getProblemService := problemservice.NewGetProblemService(problemRepo)
	listProblemsService := problemservice.NewListProblemsService(problemRepo)
	updateProblemService := problemservice.NewUpdateProblemService(problemRepo, clock)
	deleteProblemService := problemservice.NewDeleteProblemService(problemRepo)

	handler := httpproblem.NewHandler(createProblemService, getProblemService, listProblemsService, updateProblemService, deleteProblemService)

	mux := http.NewServeMux()
	httpproblem.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("problem service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/assessment/config"
	httpassessment "github.com/aakashloyar/elevate/assessment/internal/adapter/in/http/assessment"
	postgres "github.com/aakashloyar/elevate/assessment/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out/system"
	assessmentsvc "github.com/aakashloyar/elevate/assessment/internal/application/service/assessment"
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

	assessmentRepo := postgres.NewAssessmentRepository(db)
	if err := assessmentRepo.(*postgres.AssessmentRepository).Migrate(); err != nil {
		log.Fatalf("failed to migrate assessments table: %v", err)
	}

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	createAssessmentService := assessmentsvc.NewCreateAssessmentService(assessmentRepo, idGen, clock)
	getAssessmentService := assessmentsvc.NewGetAssessmentService(assessmentRepo)
	deleteAssessmentService := assessmentsvc.NewDeleteAssessmentService(assessmentRepo)

	handler := httpassessment.NewHandler(createAssessmentService, getAssessmentService, deleteAssessmentService)

	mux := http.NewServeMux()
	httpassessment.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("assessment service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

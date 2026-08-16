package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/assessment/config"
	httpassessment "github.com/aakashloyar/elevate/assessment/internal/adapter/in/http"
	kafkaconsumer "github.com/aakashloyar/elevate/assessment/internal/adapter/in/kafka"
	postgres "github.com/aakashloyar/elevate/assessment/internal/adapter/out/postgres"
	problemhttp "github.com/aakashloyar/elevate/assessment/internal/adapter/out/problemhttp"
	"github.com/aakashloyar/elevate/assessment/internal/application/ports/out/system"
	assessmentsvc "github.com/aakashloyar/elevate/assessment/internal/application/service"
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

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	createAssessmentService := assessmentsvc.NewCreateAssessmentService(assessmentRepo, idGen, clock)
	listAssessmentsService := assessmentsvc.NewListAssessmentsService(assessmentRepo)
	getAssessmentService := assessmentsvc.NewGetAssessmentService(assessmentRepo)
	deleteAssessmentService := assessmentsvc.NewDeleteAssessmentService(assessmentRepo)
	problemClient := problemhttp.NewClient(config.App.Services.ProblemServiceURL)
	createAssessmentProblemService := assessmentsvc.NewCreateAssessmentProblemService(assessmentRepo, problemClient)
	addProblemsBatchService := assessmentsvc.NewAddProblemsBatchService(assessmentRepo)
	getMarkingSchemeService := assessmentsvc.NewGetAssessmentMarkingSchemeService(assessmentRepo)
	upsertMarkingSchemeService := assessmentsvc.NewUpsertAssessmentMarkingSchemeService(assessmentRepo)
	createMarkingSchemeService := assessmentsvc.NewCreateAssessmentMarkingSchemeService(assessmentRepo)

	handler := httpassessment.NewHandler(createAssessmentService, listAssessmentsService, getAssessmentService, deleteAssessmentService, createAssessmentProblemService, getMarkingSchemeService, upsertMarkingSchemeService, createMarkingSchemeService)

	mux := http.NewServeMux()
	httpassessment.RegisterRoutes(mux, handler)

	// Start Kafka consumer for problem-created events
	consumerConfig := kafkaconsumer.Config{
		Brokers:   config.App.Kafka.Brokers,
		Topics:    config.App.Kafka.Topics,
		ClientID:  config.App.Kafka.ClientID,
		GroupID:   config.App.Kafka.GroupID,
		APIKey:    config.App.Kafka.APIKey,
		APISecret: config.App.Kafka.APISecret,
	}
	consumer, err := consumerConfig.NewConsumer(addProblemsBatchService)
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
	defer consumer.Close()

	ctx := context.Background()
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("assessment kafka consumer stopped: %v", err)
		}
	}()

	serverPort := config.App.Server.Port
	log.Printf("assessment service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

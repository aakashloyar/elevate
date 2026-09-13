package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/aakashloyar/elevate/submission/config"
	httpsubmission "github.com/aakashloyar/elevate/submission/internal/adapter/in/http"
	"github.com/aakashloyar/elevate/submission/internal/adapter/in/worker"
	assessmenthttp "github.com/aakashloyar/elevate/submission/internal/adapter/out/assessmenthttp"
	kafkaproducer "github.com/aakashloyar/elevate/submission/internal/adapter/out/kafka"
	postgres "github.com/aakashloyar/elevate/submission/internal/adapter/out/postgres"
	problemhttp "github.com/aakashloyar/elevate/submission/internal/adapter/out/problemhttp"
	userhttp "github.com/aakashloyar/elevate/submission/internal/adapter/out/userhttp"
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

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}
	producer, err := kafkaproducer.NewProducer(kafkaproducer.Config{
		Brokers:   config.App.Kafka.Brokers,
		ClientID:  config.App.Kafka.ClientID,
		APIKey:    config.App.Kafka.APIKey,
		APISecret: config.App.Kafka.APISecret,
	})
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	assessmentClient := assessmenthttp.NewClient(config.App.Services.AssessmentServiceURL)
	userClient := userhttp.NewClient(config.App.Services.UserServiceURL)
	problemClient := problemhttp.NewClient(config.App.Services.ProblemServiceURL)
	createSubmissionService := submissionservice.NewCreateSubmissionService(submissionRepo, assessmentClient, userClient, idGen, clock)
	startSubmissionService := submissionservice.NewStartSubmissionService(submissionRepo, clock)
	saveAnswerService := submissionservice.NewSaveAnswerService(submissionRepo, problemClient, clock)
	saveAnswerBatchService := submissionservice.NewSaveAnswerBatchService(saveAnswerService)
	getSubmissionService := submissionservice.NewGetSubmissionService(submissionRepo)
	getSubmissionStatusService := submissionservice.NewGetSubmissionStatusService(submissionRepo)
	submitSubmissionService := submissionservice.NewSubmitSubmissionService(submissionRepo, clock, producer, config.App.Kafka.SubmissionSubmittedTopic)
	updateSubmissionStatusService := submissionservice.NewUpdateSubmissionStatusService(submissionRepo)
	expireSubmissionsService := submissionservice.NewExpireSubmissionsService(submissionRepo, clock, producer, config.App.Kafka.SubmissionSubmittedTopic)

	handler := httpsubmission.NewHandler(createSubmissionService, startSubmissionService, saveAnswerService, saveAnswerBatchService, getSubmissionService, getSubmissionStatusService, submitSubmissionService, updateSubmissionStatusService)
	expirationWorker := worker.NewExpirationWorker(expireSubmissionsService)
	workerContext, stopWorker := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopWorker()
	go expirationWorker.Start(workerContext)

	mux := http.NewServeMux()
	httpsubmission.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("submission service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

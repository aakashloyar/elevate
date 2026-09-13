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

	problemClient := problemhttp.NewClient(config.App.Services.ProblemServiceURL)
	userClient := userhttp.NewClient(config.App.Services.UserServiceURL)
	assessmentClient := assessmenthttp.NewClient(config.App.Services.AssessmentServiceURL)
	createSubmissionService := submissionservice.NewCreateSubmissionService(submissionRepo, problemClient, assessmentClient, userClient, idGen, clock)
	startSubmissionService := submissionservice.NewStartSubmissionService(submissionRepo, clock)
	saveAnswerService := submissionservice.NewSaveAnswerService(submissionRepo, clock)
	saveAnswerBatchService := submissionservice.NewSaveAnswerBatchService(submissionRepo, clock)
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
	if err := http.ListenAndServe(":"+serverPort, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
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

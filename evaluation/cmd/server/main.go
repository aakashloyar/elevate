package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aakashloyar/elevate/evaluation/config"
	httpapi "github.com/aakashloyar/elevate/evaluation/internal/adapter/in/http"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/in/kafka"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/out/assessmenthttp"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/out/problemhttp"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/out/submissionhttp"
	evaluationservice "github.com/aakashloyar/elevate/evaluation/internal/application/service"
)

func main() {
	cfg := config.Load()
	db, err := postgres.Open(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repository := postgres.New(db)
	evaluateSubmissionService := evaluationservice.NewEvaluateSubmissionService(assessmenthttp.NewClient(cfg.AssessmentServiceURL), problemhttp.NewClient(cfg.ProblemServiceURL), submissionhttp.NewClient(cfg.SubmissionServiceURL), repository)
	getEvaluationService := evaluationservice.NewGetEvaluationService(repository)
	consumerConfig := kafka.Config{Brokers: cfg.Kafka.Brokers, Topics: []string{cfg.Kafka.Topic}, ClientID: cfg.Kafka.ClientID, GroupID: cfg.Kafka.GroupID, APIKey: cfg.Kafka.APIKey, APISecret: cfg.Kafka.APISecret}
	consumer, err := consumerConfig.NewConsumer(evaluateSubmissionService)
	if err != nil {
		log.Fatalf("create kafka consumer: %v", err)
	}
	defer consumer.Close()
	go func() {
		if err := consumer.Start(context.Background()); err != nil {
			log.Printf("evaluation consumer stopped: %v", err)
		}
	}()
	mux := http.NewServeMux()
	httpapi.RegisterRoutes(mux, httpapi.NewHandler(getEvaluationService))
	log.Printf("evaluation service listening on :%s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, mux))
}

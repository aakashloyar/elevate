package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aakashloyar/elevate/evaluation/config"
	httpapi "github.com/aakashloyar/elevate/evaluation/internal/adapter/in/http"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/in/kafka"
	httpclient "github.com/aakashloyar/elevate/evaluation/internal/adapter/out/http"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/evaluation/internal/adapter/out/submissionhttp"
	"github.com/aakashloyar/elevate/evaluation/internal/application"
)

func main() {
	cfg := config.Load()
	db, err := postgres.Open(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	service := application.NewService(httpclient.New(cfg.AssessmentServiceURL), httpclient.New(cfg.ProblemServiceURL), submissionhttp.NewClient(cfg.SubmissionServiceURL), postgres.New(db))
	consumerConfig := kafka.Config{Brokers: cfg.Kafka.Brokers, Topics: []string{cfg.Kafka.Topic}, ClientID: cfg.Kafka.ClientID, GroupID: cfg.Kafka.GroupID, APIKey: cfg.Kafka.APIKey, APISecret: cfg.Kafka.APISecret}
	consumer, err := consumerConfig.NewConsumer(service)
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
	httpapi.New(service).Register(mux)
	log.Printf("evaluation service listening on :%s", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, mux))
}

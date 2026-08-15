package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/problem/config"
	httpproblem "github.com/aakashloyar/elevate/problem/internal/adapter/in/http"
	kafkaconsumer "github.com/aakashloyar/elevate/problem/internal/adapter/in/kafka"
	kafkaproducer "github.com/aakashloyar/elevate/problem/internal/adapter/out/kafka"
	postgres "github.com/aakashloyar/elevate/problem/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/problem/internal/application/ports/out/system"
	problemservice "github.com/aakashloyar/elevate/problem/internal/application/service"
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

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	createProblemService := problemservice.NewCreateProblemService(problemRepo, idGen, clock)
	getProblemService := problemservice.NewGetProblemService(problemRepo)
	listProblemsService := problemservice.NewListProblemsService(problemRepo)
	updateProblemService := problemservice.NewUpdateProblemService(problemRepo, idGen, clock)
	deleteProblemService := problemservice.NewDeleteProblemService(problemRepo)

	handler := httpproblem.NewHandler(createProblemService, getProblemService, listProblemsService, updateProblemService, deleteProblemService)

	mux := http.NewServeMux()
	httpproblem.RegisterRoutes(mux, handler)

	// Start Kafka consumer to receive generated-problem batches and publish created ids
	brokers := config.App.Kafka.Brokers

	producer, err := kafkaproducer.NewProducer(kafkaproducer.Config{
		Brokers:   brokers,
		ClientID:  config.App.Kafka.ClientID,
		APIKey:    config.App.Kafka.APIKey,
		APISecret: config.App.Kafka.APISecret,
	})

	createProblemBatchService := problemservice.NewCreateProblemBatchService(createProblemService, producer, idGen, clock)

	if err != nil {
		log.Fatalf("failed to create kafka producer: %v", err)
	}
	defer producer.Close()

	consumerConfig := kafkaconsumer.Config{
		Brokers:   config.App.Kafka.Brokers,
		Topics:    config.App.Kafka.Topics,
		ClientID:  config.App.Kafka.ClientID,
		GroupID:   config.App.Kafka.GroupID,
		APIKey:    config.App.Kafka.APIKey,
		APISecret: config.App.Kafka.APISecret,
	}
	consumer, err := consumerConfig.NewConsumer(createProblemBatchService)
	if err != nil {
		log.Fatalf("failed to create kafka consumer: %v", err)
	}
	defer consumer.Close()

	ctx := context.Background()
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("problem kafka consumer stopped: %v", err)
		}
	}()

	serverPort := config.App.Server.Port
	log.Printf("problem service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

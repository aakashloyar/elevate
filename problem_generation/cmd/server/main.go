package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/problem_generation/config"
	httpgenerationjob "github.com/aakashloyar/elevate/problem_generation/internal/adapter/in/http"
	kafkaconsumer "github.com/aakashloyar/elevate/problem_generation/internal/adapter/in/kafka"
	"github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/ai/gemini"
	assessmenthttp "github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/assessmenthttp"
	kafkaproducer "github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/kafka"
	postgres "github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/processor"
	userhttp "github.com/aakashloyar/elevate/problem_generation/internal/adapter/out/userhttp"
	"github.com/aakashloyar/elevate/problem_generation/internal/application/ports/out/system"
	generationjobsvc "github.com/aakashloyar/elevate/problem_generation/internal/application/service"
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
	if err := jobRepo.Migrate(); err != nil {
		log.Fatalf("failed to migrate generation jobs tables: %v", err)
	}

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}
	eventPublisher, err := kafkaproducer.NewProducer(kafkaproducer.Config{
		Brokers:                 config.App.Kafka.Brokers,
		GenerationRequestsTopic: config.App.Kafka.GenerationRequestsTopic,
		ClientID:                config.App.Kafka.ClientID,
		APIKey:                  config.App.Kafka.APIKey,
		APISecret:               config.App.Kafka.APISecret,
	})
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer eventPublisher.Close()

	userClient := userhttp.NewClient(config.App.Services.UserServiceURL)
	assessmentClient := assessmenthttp.NewClient(config.App.Services.AssessmentServiceURL)
	createJobService := generationjobsvc.NewCreateGenerationJobService(jobRepo, eventPublisher, userClient, assessmentClient, idGen, clock)
	getJobService := generationjobsvc.NewGetGenerationJobService(jobRepo)

	if config.App.Worker.Enabled {
		aiClient := gemini.Config{
			BaseURL: config.App.AI.BaseURL,
			APIKey:  config.App.AI.APIKey,
			Model:   config.App.AI.Model,
		}.NewClient()
		generationProcessor := processor.NewAIProblemGenerator(aiClient, eventPublisher, config.App.Kafka.GeneratedProblemsTopic)
		processJobService := generationjobsvc.NewProcessGenerationJobService(jobRepo, generationProcessor)

		consumer, err := kafkaconsumer.Config{
			Brokers:                 config.App.Kafka.Brokers,
			GenerationRequestsTopic: config.App.Kafka.GenerationRequestsTopic,
			ClientID:                config.App.Kafka.ClientID,
			GroupID:                 config.App.Kafka.GroupID,
			APIKey:                  config.App.Kafka.APIKey,
			APISecret:               config.App.Kafka.APISecret,
		}.NewConsumer(processJobService)
		if err != nil {
			log.Fatalf("failed to create Kafka consumer: %v", err)
		}
		defer consumer.Close()

		go func() {
			if err := consumer.Start(context.Background()); err != nil {
				log.Printf("generation worker stopped: %v", err)
			}
		}()
		log.Printf("generation worker enabled on topic %q with group %q", config.App.Kafka.GenerationRequestsTopic, config.App.Kafka.GroupID)
	} else {
		log.Printf("generation worker disabled; set GENERATION_WORKER_ENABLED=true to consume generation requests")
	}

	handler := httpgenerationjob.NewHandler(createJobService, getJobService)

	mux := http.NewServeMux()
	httpgenerationjob.RegisterRoutes(mux, handler)

	serverPort := config.App.Server.Port
	log.Printf("problem_generation service starting on :%s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

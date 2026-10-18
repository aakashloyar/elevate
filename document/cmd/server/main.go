package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	config "github.com/aakashloyar/elevate/document/config"
	chunkhttp "github.com/aakashloyar/elevate/document/internal/adapter/in/http/chunk"
	uploadhttp "github.com/aakashloyar/elevate/document/internal/adapter/in/http/upload"
	kafkaproducer "github.com/aakashloyar/elevate/document/internal/adapter/out/kafka"
	"github.com/aakashloyar/elevate/document/internal/adapter/out/postgres"
	s3adapter "github.com/aakashloyar/elevate/document/internal/adapter/out/s3"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out/system"
	chunkSvc "github.com/aakashloyar/elevate/document/internal/application/service/chunk"
	uploadsvc "github.com/aakashloyar/elevate/document/internal/application/service/upload"
)

func main() {
	ctx := context.Background()
	appConfig := config.App

	s3Config := s3adapter.Config{
		Region: appConfig.S3.Region,
		Bucket: appConfig.S3.Bucket,
	}
	s3Client, err := s3Config.NewS3Client(ctx)
	if err != nil {
		log.Fatalf("initialize S3 client: %v", err)
	}

	storage := s3adapter.NewS3Storage(s3Client.Client, s3Client.Bucket)
	eventPublisher, err := kafkaproducer.NewProducer(kafkaproducer.Config{
		Brokers:   config.App.Kafka.Brokers,
		Topic:     config.App.Kafka.Topic,
		ClientID:  config.App.Kafka.ClientID,
		APIKey:    config.App.Kafka.APIKey,
		APISecret: config.App.Kafka.APISecret,
	})
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer eventPublisher.Close()

	port, err := strconv.Atoi(config.App.Postgres.Port)
	if err != nil {
		log.Fatalf("invalid POSTGRES_PORT: %v", err)
	}
	postgresConfig := postgres.Config{
		Host:     config.App.Postgres.Host,
		Port:     port,
		User:     config.App.Postgres.User,
		Password: config.App.Postgres.Password,
		DBName:   config.App.Postgres.DBName,
		SSLMode:  config.App.Postgres.SSLMode,
	}

	db, err := postgresConfig.NewDB()

	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	clock := system.SystemClock{}
	idGen := system.UUIDGenerator{}

	uploadRepo := postgres.NewUploadRepository(db)
	chunkRepo := postgres.NewChunkRepository(db)

	initUploadService := uploadsvc.NewInitUploadService(uploadRepo, storage, idGen, clock)
	completeUploadService := uploadsvc.NewCompleteUploadService(uploadRepo, storage, eventPublisher)
	markChunkService := chunkSvc.NewMarkChunkService(chunkRepo, idGen, clock)

	uploadHandler := uploadhttp.NewHandler(initUploadService, completeUploadService)
	chunkHandler := chunkhttp.NewHandler(markChunkService)

	mux := http.NewServeMux()

	uploadhttp.RegisterRoutes(mux, uploadHandler)
	chunkhttp.RegisterRoutes(mux, chunkHandler)

	log.Println("Routes registered successfully. Starting server on :8080...")
	http.ListenAndServe(":8080", mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

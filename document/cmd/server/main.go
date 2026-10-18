package main

import (
	"strconv"
	"context"
	"log"
	"net/http"

	config "github.com/aakashloyar/elevate/document/config"
	chunkhttp "github.com/aakashloyar/elevate/document/internal/adapter/in/http/chunk"
	uploadhttp "github.com/aakashloyar/elevate/document/internal/adapter/in/http/upload"
	s3adapter "github.com/aakashloyar/elevate/document/internal/adapter/out/s3"
	"github.com/aakashloyar/elevate/document/internal/application/ports/out/system"
	uploadsvc "github.com/aakashloyar/elevate/document/internal/application/service/upload"
	chunkSvc "github.com/aakashloyar/elevate/document/internal/application/service/chunk"
	"github.com/aakashloyar/elevate/document/internal/adapter/out/postgres"
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
	completeUploadService := uploadsvc.NewCompleteUploadService(uploadRepo, storage)
	markChunkService := chunkSvc.NewMarkChunkService(chunkRepo, idGen, clock)


	uploadHandler := uploadhttp.NewHandler(initUploadService, completeUploadService)
	chunkHandler := chunkhttp.NewHandler(markChunkService)

	mux:= http.NewServeMux()

	uploadhttp.RegisterRoutes(mux, uploadHandler)
	chunkhttp.RegisterRoutes(mux, chunkHandler)
	
	log.Println("Routes registered successfully. Starting server on :8080...")
	http.ListenAndServe(":8080",mux)
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

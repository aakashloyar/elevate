package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/aakashloyar/elevate/auth/config"
	api "github.com/aakashloyar/elevate/auth/internal/adapter/in/http"
	kafkaproducer "github.com/aakashloyar/elevate/auth/internal/adapter/out/kafka"
	"github.com/aakashloyar/elevate/auth/internal/adapter/out/postgres"
	"github.com/aakashloyar/elevate/auth/internal/adapter/out/userhttp"
	"github.com/aakashloyar/elevate/auth/internal/httpobserve"
	"github.com/aakashloyar/elevate/auth/internal/service"
)

func main() {
	port, _ := strconv.Atoi(config.App.Postgres.Port)
	db, err := postgres.Open(config.App.Postgres.Host, strconv.Itoa(port), config.App.Postgres.User, config.App.Postgres.Password, config.App.Postgres.DBName, config.App.Postgres.SSLMode)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repo := postgres.NewRepository(db)
	producer, err := kafkaproducer.NewProducer(kafkaproducer.Config{
		Brokers: config.App.KafkaBrokers, ClientID: config.App.KafkaClientID,
		APIKey: config.App.KafkaAPIKey, APISecret: config.App.KafkaAPISecret,
		Topic: config.App.NotificationEmailTopic,
	})
	if err != nil {
		log.Fatalf("failed to create notification Kafka producer: %v", err)
	}
	defer producer.Close()
	svc := service.New(repo, userhttp.New(config.App.UserServiceURL), producer, config.App.JWTSecret, config.App.TokenTTL)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, api.NewHandler(svc))
	log.Printf("auth service starting on :%s", config.App.HTTPPort)
	log.Fatal(http.ListenAndServe(":"+config.App.HTTPPort, withCORS(httpobserve.Middleware("auth", mux))))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

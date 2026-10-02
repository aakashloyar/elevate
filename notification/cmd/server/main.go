package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aakashloyar/elevate/notification/config"
	api "github.com/aakashloyar/elevate/notification/internal/adapter/in/http"
	kafkaconsumer "github.com/aakashloyar/elevate/notification/internal/adapter/in/kafka"
	"github.com/aakashloyar/elevate/notification/internal/email"
	"github.com/aakashloyar/elevate/notification/internal/httpobserve"
)

func main() {
	sender := email.New(config.App.SMTPHost, config.App.SMTPPort, config.App.SMTPUsername, config.App.SMTPPassword, config.App.From)
	var consumer *kafkaconsumer.Consumer
	if config.App.KafkaBrokers[0] != "" {
		var err error
		consumer, err = kafkaconsumer.NewConsumer(kafkaconsumer.Config{
			Brokers: config.App.KafkaBrokers, Topic: config.App.NotificationEmailTopic,
			GroupID: config.App.KafkaGroupID, ClientID: config.App.KafkaClientID,
			APIKey: config.App.KafkaAPIKey, APISecret: config.App.KafkaAPISecret,
		}, sender)
		if err != nil {
			log.Fatalf("failed to create notification Kafka consumer: %v", err)
		}
		defer consumer.Close()
		go func() {
			if err := consumer.Start(context.Background()); err != nil {
				log.Printf("notification Kafka consumer stopped: %v", err)
			}
		}()
	}
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, api.NewHandler(sender))
	log.Printf("notification service starting on :%s", config.App.Port)
	log.Fatal(http.ListenAndServe(":"+config.App.Port, withCORS(httpobserve.Middleware("notification", mux))))
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

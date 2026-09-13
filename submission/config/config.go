package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port string
}

type ServiceConfig struct {
	AssessmentServiceURL string
	UserServiceURL       string
	ProblemServiceURL    string
}

type KafkaConfig struct {
	Brokers                  []string
	ClientID                 string
	APIKey                   string
	APISecret                string
	SubmissionSubmittedTopic string
}

type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
	Services ServiceConfig
	Kafka    KafkaConfig
}

func load() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env: %v", err)
	}

	postgres := PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	server := ServerConfig{Port: os.Getenv("HTTP_PORT")}
	if server.Port == "" {
		server.Port = "8080"
	}

	kafka := KafkaConfig{
		Brokers:                  strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		ClientID:                 os.Getenv("KAFKA_CLIENT_ID"),
		APIKey:                   os.Getenv("KAFKA_API_KEY"),
		APISecret:                os.Getenv("KAFKA_API_SECRET"),
		SubmissionSubmittedTopic: os.Getenv("KAFKA_SUBMISSION_SUBMITTED_TOPIC"),
	}

	services := ServiceConfig{
		AssessmentServiceURL: os.Getenv("ASSESSMENT_SERVICE_URL"),
		ProblemServiceURL:    os.Getenv("PROBLEM_SERVICE_URL"),
		UserServiceURL:       os.Getenv("USER_SERVICE_URL"),
	}
	if services.ProblemServiceURL == "" {
		services.ProblemServiceURL = "http://localhost:8083"
	}
	if services.UserServiceURL == "" {
		services.UserServiceURL = "http://localhost:8081"
	}
	if services.AssessmentServiceURL == "" {
		services.AssessmentServiceURL = "http://localhost:8082"
	}

	return Config{Postgres: postgres, Server: server, Services: services, Kafka: kafka}
}

var App = load()

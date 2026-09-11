package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	MaxProblemsPerBatch = 50
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
	ProblemServiceURL string
}

type KafkaConfig struct {
	Brokers                  []string
	ProblemCreatedBatchTopic string
	ClientID                 string
	GroupID                  string
	APIKey                   string
	APISecret                string
}

type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
	Services ServiceConfig
	Kafka    KafkaConfig
}

func load() Config {
	err := godotenv.Load()
	if err != nil {
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

	services := ServiceConfig{
		ProblemServiceURL: os.Getenv("PROBLEM_SERVICE_URL"),
	}
	if services.ProblemServiceURL == "" {
		services.ProblemServiceURL = "http://localhost:8080"
	}

	kafka := KafkaConfig{
		Brokers:                  splitList(os.Getenv("KAFKA_BROKERS")),
		ProblemCreatedBatchTopic: os.Getenv("KAFKA_PROBLEM_CREATED_BATCH_TOPIC"),
		ClientID:                 os.Getenv("KAFKA_CLIENT_ID"),
		GroupID:                  os.Getenv("KAFKA_GROUP_ID"),
		APIKey:                   os.Getenv("KAFKA_API_KEY"),
		APISecret:                os.Getenv("KAFKA_API_SECRET"),
	}
	return Config{Postgres: postgres, Server: server, Services: services, Kafka: kafka}
}

var App = load()

func splitList(value string) []string {
	values := strings.Split(value, ",")
	result := make([]string, 0, len(values))
	for _, item := range values {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

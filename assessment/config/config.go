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
	Brokers   []string
	Topics    []string
	ClientID  string
	GroupID   string
	APIKey    string
	APISecret string
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
		Brokers:   strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		Topics:    strings.Split(os.Getenv("KAFKA_TOPICS"), ","),
		ClientID:  os.Getenv("KAFKA_CLIENT_ID"),
		GroupID:   os.Getenv("KAFKA_GROUP_ID"),
		APIKey:    os.Getenv("KAFKA_API_KEY"),
		APISecret: os.Getenv("KAFKA_API_SECRET"),
	}
	if len(kafka.Topics) == 0 {
		log.Fatal("Did not find any kafka topics")
	}

	return Config{Postgres: postgres, Server: server, Services: services, Kafka: kafka}
}

var App = load()

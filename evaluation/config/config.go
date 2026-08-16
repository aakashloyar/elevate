package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort             string
	AssessmentServiceURL string
	ProblemServiceURL    string
	Postgres             PostgresConfig
	Kafka                KafkaConfig
}

type PostgresConfig struct{ Host, Port, User, Password, DBName, SSLMode string }
type KafkaConfig struct {
	Brokers                                     []string
	Topic, ClientID, GroupID, APIKey, APISecret string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env: %v", err)
	}
	cfg := Config{
		HTTPPort:             os.Getenv("HTTP_PORT"),
		AssessmentServiceURL: os.Getenv("ASSESSMENT_SERVICE_URL"),
		ProblemServiceURL:    os.Getenv("PROBLEM_SERVICE_URL"),
		Postgres:             PostgresConfig{Host: os.Getenv("POSTGRES_HOST"), Port: os.Getenv("POSTGRES_PORT"), User: os.Getenv("POSTGRES_USER"), Password: os.Getenv("POSTGRES_PASSWORD"), DBName: os.Getenv("POSTGRES_DB"), SSLMode: os.Getenv("POSTGRES_SSLMODE")},
		Kafka:                KafkaConfig{Brokers: split(os.Getenv("KAFKA_BROKERS")), Topic: os.Getenv("KAFKA_SUBMISSION_SUBMITTED_TOPIC"), ClientID: os.Getenv("KAFKA_CLIENT_ID"), GroupID: os.Getenv("KAFKA_GROUP_ID"), APIKey: os.Getenv("KAFKA_API_KEY"), APISecret: os.Getenv("KAFKA_API_SECRET")},
	}
	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8084"
	}
	if cfg.AssessmentServiceURL == "" {
		cfg.AssessmentServiceURL = "http://localhost:8082"
	}
	if cfg.ProblemServiceURL == "" {
		cfg.ProblemServiceURL = "http://localhost:8081"
	}
	if cfg.Kafka.Topic == "" {
		cfg.Kafka.Topic = "submission-submitted"
	}
	if cfg.Kafka.GroupID == "" {
		cfg.Kafka.GroupID = "evaluation-service"
	}
	return cfg
}

func split(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, ",")
}

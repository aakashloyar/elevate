package config

import (
	"log"
	"os"
	"strconv"
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

type KafkaConfig struct {
	Brokers                []string
	Topic                  string
	GeneratedProblemsTopic string
	ClientID               string
	GroupID                string
	APIKey                 string
	APISecret              string
}

type WorkerConfig struct {
	Enabled bool
}

type AIConfig struct {
	Provider string
	BaseURL  string
	APIKey   string
	Model    string
}

type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
	Kafka    KafkaConfig
	Worker   WorkerConfig
	AI       AIConfig
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

	kafka := KafkaConfig{
		Brokers:                strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		Topic:                  os.Getenv("KAFKA_GENERATION_REQUESTS_TOPIC"),
		GeneratedProblemsTopic: os.Getenv("KAFKA_GENERATED_PROBLEMS_TOPIC"),
		ClientID:               os.Getenv("KAFKA_CLIENT_ID"),
		GroupID:                os.Getenv("KAFKA_GROUP_ID"),
		APIKey:                 os.Getenv("KAFKA_API_KEY"),
		APISecret:              os.Getenv("KAFKA_API_SECRET"),
	}
	if kafka.Topic == "" {
		kafka.Topic = "generation-requests"
	}
	if kafka.GeneratedProblemsTopic == "" {
		kafka.GeneratedProblemsTopic = "generated-problems"
	}
	if kafka.GroupID == "" {
		kafka.GroupID = "problem-generation-service"
	}

	worker := WorkerConfig{Enabled: boolEnv("GENERATION_WORKER_ENABLED", false)}
	ai := AIConfig{
		Provider: strings.ToLower(strings.TrimSpace(os.Getenv("AI_PROVIDER"))),
		BaseURL:  os.Getenv("AI_BASE_URL"),
		APIKey:   firstNonEmpty(os.Getenv("GEMINI_API_KEY"), os.Getenv("Elevate_Gemini_API_Key")),
		Model:    os.Getenv("AI_MODEL"),
	}
	if ai.Provider == "" {
		ai.Provider = "gemini"
	}
	if ai.Model == "" {
		ai.Model = "gemini-flash-latest"
	}

	return Config{Postgres: postgres, Server: server, Kafka: kafka, Worker: worker, AI: ai}
}

var App = load()

func boolEnv(name string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("warning: invalid %s value %q, using %v", name, value, fallback)
		return fallback
	}
	return parsed
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

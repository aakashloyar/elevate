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
		Brokers:   strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		Topics:    strings.Split(os.Getenv("KAFKA_TOPICS"), ","),
		ClientID:  os.Getenv("KAFKA_CLIENT_ID"),
		GroupID:   os.Getenv("KAFKA_GROUP_ID"),
		APIKey:    os.Getenv("KAFKA_API_KEY"),
		APISecret: os.Getenv("KAFKA_API_SECRET"),
	}
	return Config{Postgres: postgres, Server: server, Kafka: kafka}
}

var (
	App = load()
	CreatedProblemBatchTopic = strings.Split(os.Getenv("KAFKA_TOPICS"), ",")[0]
)

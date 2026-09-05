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
	Topic     string
	ClientID  string
	APIKey    string
	APISecret string
}

type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
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

	kafka := KafkaConfig{
		Brokers:   strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		Topic:     os.Getenv("KAFKA_GENERATION_REQUESTS_TOPIC"),
		ClientID:  os.Getenv("KAFKA_CLIENT_ID"),
		APIKey:    os.Getenv("KAFKA_API_KEY"),
		APISecret: os.Getenv("KAFKA_API_SECRET"),
	}
	if kafka.Topic == "" {
		kafka.Topic = "generation-requests"
	}

	return Config{Postgres: postgres, Server: server, Kafka: kafka}
}

var App = load()

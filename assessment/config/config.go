package config

import (
	"log"
	"os"

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

type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
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

	return Config{Postgres: postgres, Server: server}
}

var App = load()

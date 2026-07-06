package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

type ServiceConfig struct {
	AssessmentServiceURL string
	ProblemServiceURL    string
	SubmissionServiceURL string
}

type Config struct {
	Server   ServerConfig
	Services ServiceConfig
}

func load() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: could not load .env: %v", err)
	}

	server := ServerConfig{Port: os.Getenv("HTTP_PORT")}
	if server.Port == "" {
		server.Port = "8080"
	}

	services := ServiceConfig{
		AssessmentServiceURL: os.Getenv("ASSESSMENT_SERVICE_URL"),
		ProblemServiceURL:    os.Getenv("PROBLEM_SERVICE_URL"),
		SubmissionServiceURL: os.Getenv("SUBMISSION_SERVICE_URL"),
	}
	if services.AssessmentServiceURL == "" {
		services.AssessmentServiceURL = "http://localhost:8081"
	}
	if services.ProblemServiceURL == "" {
		services.ProblemServiceURL = "http://localhost:8082"
	}
	if services.SubmissionServiceURL == "" {
		services.SubmissionServiceURL = "http://localhost:8083"
	}

	return Config{Server: server, Services: services}
}

var App = load()

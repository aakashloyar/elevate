package config

import "os"

type Config struct {
	HTTPPort             string
	AssessmentServiceURL string
	ProblemServiceURL    string
	SubmissionServiceURL string
}

func Load() Config {
	cfg := Config{
		HTTPPort:             os.Getenv("HTTP_PORT"),
		AssessmentServiceURL: os.Getenv("ASSESSMENT_SERVICE_URL"),
		ProblemServiceURL:    os.Getenv("PROBLEM_SERVICE_URL"),
		SubmissionServiceURL: os.Getenv("SUBMISSION_SERVICE_URL"),
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
	if cfg.SubmissionServiceURL == "" {
		cfg.SubmissionServiceURL = "http://localhost:8083"
	}
	return cfg
}

package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type UploadConfig struct {
	PresignExpirty int64
	MaxChunkSize   int64
}

type S3Config struct {
	Region string
	Bucket string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers               []string
	DocumentUploadedTopic string
	ClientID              string
	GroupID               string
	APIKey                string
	APISecret             string
}

type GeminiConfig struct {
	BaseURL          string
	APIKey           string
	EmbeddingModel   string
	OutputDimensions int
}

type Config struct {
	S3       S3Config
	Upload   UploadConfig
	Postgres PostgresConfig
	Kafka    KafkaConfig
	Gemini   GeminiConfig
	HTTPPort string
}

func load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env", err)
	}

	upload := UploadConfig{
		MaxChunkSize:   5 * 1024 * 1024, // 5MB
		PresignExpirty: 15,              //15min
	}

	s3 := S3Config{
		Region: os.Getenv("AWS_REGION"),
		Bucket: os.Getenv("AWS_BUCKET"),
	}

	postgres := PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	kafka := KafkaConfig{
		Brokers:               strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		DocumentUploadedTopic: os.Getenv("KAFKA_DOCUMENT_UPLOADED_TOPIC"),
		ClientID:              os.Getenv("KAFKA_CLIENT_ID"),
		GroupID:               os.Getenv("KAFKA_DOCUMENT_GROUP_ID"),
		APIKey:                os.Getenv("KAFKA_API_KEY"),
		APISecret:             os.Getenv("KAFKA_API_SECRET"),
	}
	outputDimensions := 768
	if value, err := strconv.Atoi(os.Getenv("GEMINI_EMBEDDING_DIMENSIONS")); err == nil && value > 0 {
		outputDimensions = value
	}
	gemini := GeminiConfig{
		BaseURL:          os.Getenv("GEMINI_API_BASE_URL"),
		APIKey:           os.Getenv("GEMINI_API_KEY"),
		EmbeddingModel:   os.Getenv("GEMINI_EMBEDDING_MODEL"),
		OutputDimensions: outputDimensions,
	}
	if len(kafka.Brokers) == 0 || kafka.Brokers[0] == "" {
		kafka.Brokers = []string{os.Getenv("KAFKA_BROKER")}
	}
	if kafka.DocumentUploadedTopic == "" {
		kafka.DocumentUploadedTopic = os.Getenv("KAFKA_TOPIC")
	}
	if kafka.DocumentUploadedTopic == "" {
		kafka.DocumentUploadedTopic = "document.uploaded"
	}
	if kafka.GroupID == "" {
		kafka.GroupID = "document-processing-service"
	}

	return Config{
		S3:       s3,
		Upload:   upload,
		Postgres: postgres,
		Kafka:    kafka,
		Gemini:   gemini,
		HTTPPort: os.Getenv("HTTP_PORT"),
	}
}

var App = load()

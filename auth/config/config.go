package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort               string
	UserServiceURL         string
	KafkaBrokers           []string
	KafkaClientID          string
	KafkaAPIKey            string
	KafkaAPISecret         string
	NotificationEmailTopic string
	JWTSecret              string
	TokenTTL               time.Duration
	Postgres               struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
}

func load() Config {
	_ = godotenv.Load()
	c := Config{
		HTTPPort:               os.Getenv("HTTP_PORT"),
		UserServiceURL:         os.Getenv("USER_SERVICE_URL"),
		KafkaBrokers:           strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaClientID:          os.Getenv("KAFKA_CLIENT_ID"),
		KafkaAPIKey:            os.Getenv("KAFKA_API_KEY"),
		KafkaAPISecret:         os.Getenv("KAFKA_API_SECRET"),
		NotificationEmailTopic: os.Getenv("KAFKA_NOTIFICATION_EMAIL_TOPIC"),
		JWTSecret:              os.Getenv("JWT_SECRET"),
	}
	if c.HTTPPort == "" {
		c.HTTPPort = "8088"
	}
	if c.UserServiceURL == "" {
		c.UserServiceURL = "http://localhost:8081"
	}
	if c.NotificationEmailTopic == "" {
		c.NotificationEmailTopic = "notification.email"
	}
	ttl := 3 * time.Hour
	if raw := os.Getenv("ACCESS_TOKEN_TTL_HOURS"); raw != "" {
		if hours, err := strconv.Atoi(raw); err == nil && hours > 0 {
			ttl = time.Duration(hours) * time.Hour
		}
	}
	c.TokenTTL = ttl
	c.Postgres.Host = os.Getenv("POSTGRES_HOST")
	c.Postgres.Port = os.Getenv("POSTGRES_PORT")
	c.Postgres.User = os.Getenv("POSTGRES_USER")
	c.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	c.Postgres.DBName = os.Getenv("POSTGRES_DB")
	c.Postgres.SSLMode = os.Getenv("POSTGRES_SSLMODE")
	if c.JWTSecret == "" {
		log.Println("WARNING: JWT_SECRET is empty; auth service cannot safely issue tokens")
	}
	return c
}

var App = load()

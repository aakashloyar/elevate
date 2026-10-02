package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	SMTPHost               string
	SMTPPort               string
	SMTPUsername           string
	SMTPPassword           string
	From                   string
	KafkaBrokers           []string
	KafkaClientID          string
	KafkaAPIKey            string
	KafkaAPISecret         string
	NotificationEmailTopic string
	KafkaGroupID           string
}

func load() Config {
	_ = godotenv.Load()
	c := Config{
		Port: os.Getenv("HTTP_PORT"), SMTPHost: os.Getenv("SMTP_HOST"), SMTPPort: os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"), SMTPPassword: os.Getenv("SMTP_PASSWORD"), From: os.Getenv("SMTP_FROM"),
		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","), KafkaClientID: os.Getenv("KAFKA_CLIENT_ID"),
		KafkaAPIKey: os.Getenv("KAFKA_API_KEY"), KafkaAPISecret: os.Getenv("KAFKA_API_SECRET"),
		NotificationEmailTopic: os.Getenv("KAFKA_NOTIFICATION_EMAIL_TOPIC"), KafkaGroupID: os.Getenv("KAFKA_NOTIFICATION_GROUP_ID"),
	}
	if c.Port == "" {
		c.Port = "8089"
	}
	if c.SMTPPort == "" {
		c.SMTPPort = "587"
	}
	if c.From == "" {
		c.From = c.SMTPUsername
	}
	if c.NotificationEmailTopic == "" {
		c.NotificationEmailTopic = os.Getenv("KAFKA_NOTIFICATION_TOPIC")
	}
	if c.NotificationEmailTopic == "" {
		c.NotificationEmailTopic = "notification.email"
	}
	if c.KafkaGroupID == "" {
		c.KafkaGroupID = os.Getenv("KAFKA_GROUP_ID")
	}
	if c.KafkaGroupID == "" {
		c.KafkaGroupID = "notification-service"
	}
	return c
}

var App = load()

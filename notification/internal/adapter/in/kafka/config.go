package kafka

import (
	"crypto/tls"

	"github.com/aakashloyar/elevate/notification/internal/email"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Config struct {
	Brokers   []string
	Topic     string
	GroupID   string
	ClientID  string
	APIKey    string
	APISecret string
}

type Consumer struct {
	client *kgo.Client
	sender *email.Sender
}

func NewConsumer(cfg Config, sender *email.Sender) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...), kgo.ConsumerGroup(cfg.GroupID), kgo.ConsumeTopics(cfg.Topic),
		kgo.ClientID(cfg.ClientID), kgo.DialTLSConfig(&tls.Config{}),
		kgo.SASL(plain.Auth{User: cfg.APIKey, Pass: cfg.APISecret}.AsMechanism()),
	)
	if err != nil {
		return nil, err
	}
	return &Consumer{client: client, sender: sender}, nil
}

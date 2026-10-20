package kafka

import (
	"crypto/tls"

	in "github.com/aakashloyar/elevate/document/internal/application/ports/in"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Config struct {
	Brokers               []string
	DocumentUploadedTopic string
	ClientID              string
	GroupID               string
	APIKey                string
	APISecret             string
}

type Consumer struct {
	client  *kgo.Client
	service in.ProcessDocumentService
}

func NewConsumer(cfg Config, service in.ProcessDocumentService) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.DocumentUploadedTopic),
		kgo.ClientID(cfg.ClientID),
		kgo.DialTLSConfig(&tls.Config{}),
		kgo.SASL(plain.Auth{User: cfg.APIKey, Pass: cfg.APISecret}.AsMechanism()),
	)
	if err != nil {
		return nil, err
	}
	return &Consumer{client: client, service: service}, nil
}

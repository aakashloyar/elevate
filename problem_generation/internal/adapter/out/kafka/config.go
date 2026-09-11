package kafkaproducer

import (
	"crypto/tls"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Config struct {
	Brokers                 []string
	GenerationRequestsTopic string
	ClientID                string
	APIKey                  string
	APISecret               string
}

type Producer struct {
	client *kgo.Client
	topic  string
}

func NewProducer(cfg Config) (*Producer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),
		kgo.DialTLSConfig(&tls.Config{}),
		kgo.SASL(plain.Auth{User: cfg.APIKey, Pass: cfg.APISecret}.AsMechanism()),
	)
	if err != nil {
		return nil, err
	}

	return &Producer{client: client, topic: cfg.GenerationRequestsTopic}, nil
}

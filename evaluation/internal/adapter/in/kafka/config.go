package kafka

import (
	"crypto/tls"

	"github.com/aakashloyar/elevate/evaluation/internal/application"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Config struct {
	Brokers   []string
	Topics    []string
	ClientID  string
	GroupID   string
	APIKey    string
	APISecret string
}

type Consumer struct {
	kafkaClient *KafkaConsumerClient
	service     *application.Service
}

type KafkaConsumerClient struct {
	client *kgo.Client
	topics []string
}

func (cfg Config) NewConsumer(service *application.Service) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.ClientID(cfg.ClientID),
		kgo.DialTLSConfig(&tls.Config{}),
		kgo.SASL(plain.Auth{User: cfg.APIKey, Pass: cfg.APISecret}.AsMechanism()),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		kafkaClient: &KafkaConsumerClient{client: client, topics: cfg.Topics},
		service:     service,
	}, nil
}

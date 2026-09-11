package kafka

import (
	"crypto/tls"

	in "github.com/aakashloyar/elevate/problem_generation/internal/application/ports/in"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

type Config struct {
	Brokers                 []string
	GenerationRequestsTopic string
	ClientID                string
	GroupID                 string
	APIKey                  string
	APISecret               string
}

type Consumer struct {
	kafkaClient *KafkaConsumerClient
	service     in.ProcessGenerationJobService
}

type KafkaConsumerClient struct {
	client *kgo.Client
	topic  string
}

func (cfg Config) NewConsumer(service in.ProcessGenerationJobService) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.GenerationRequestsTopic),
		kgo.ClientID(cfg.ClientID),
		kgo.DialTLSConfig(&tls.Config{}),
		kgo.SASL(plain.Auth{User: cfg.APIKey, Pass: cfg.APISecret}.AsMechanism()),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		kafkaClient: &KafkaConsumerClient{client: client, topic: cfg.GenerationRequestsTopic},
		service:     service,
	}, nil
}

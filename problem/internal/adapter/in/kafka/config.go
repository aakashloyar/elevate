package consumer

import (
	"crypto/tls"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"

	in "github.com/aakashloyar/elevate/problem/internal/application/ports/in"
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
	client      *kgo.Client
	topic       []string 
	createSvc   in.CreateProblemBatchService
	maxPerBatch int
}

type KafkaConsumerClient struct {
	client *kgo.Client
	topic  []string
}

func (cfg Config) NewConsumer(createSvc in.CreateProblemBatchService) (*Consumer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.ClientID(cfg.ClientID),

		kgo.DialTLSConfig(&tls.Config{}),

		kgo.SASL(
			plain.Auth{
				User: cfg.APIKey,
				Pass: cfg.APISecret,
			}.AsMechanism(),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:      client,
		topic:       cfg.Topics,
		createSvc:   createSvc,
		maxPerBatch: 50,
	}, nil
}

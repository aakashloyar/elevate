package consumer

import (
	"crypto/tls"

	"github.com/twmb/franz-go/pkg/kgo"
    "github.com/twmb/franz-go/pkg/sasl/plain"

	"github.com/aakashloyar/elevate/assessment/config"
	in "github.com/aakashloyar/elevate/assessment/internal/application/ports/in"
)

type Config struct {
	Brokers   []string
	Topics    []string
	ClientID  string
	GroupID   string
	APIKey    string
	APISecret string
}

type problemCreatedBatch struct {
	AssessmentID string   `json:"assessment_id"`
	ProblemIDs   []string `json:"problem_ids"`
}

type Consumer struct {
	kafkaClient    *KafkaConsumerClient
	addProblemsSvc in.AddProblemsBatchService
	maxPerBatch    int
}

type KafkaConsumerClient struct {
	client *kgo.Client
	topic  []string
}

func (cfg Config) NewConsumer(addProblemsSvc in.AddProblemsBatchService) (*Consumer, error) {
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
		kafkaClient:    &KafkaConsumerClient{client: client, topic: cfg.Topics},
		addProblemsSvc: addProblemsSvc,
		maxPerBatch:    config.MaxProblemsPerBatch,
	}, nil
}

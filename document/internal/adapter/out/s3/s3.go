package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

type S3Client struct {
	Client *s3.Client
	Bucket string
}

func (cfg Config) NewS3Client(ctx context.Context) (*S3Client, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}

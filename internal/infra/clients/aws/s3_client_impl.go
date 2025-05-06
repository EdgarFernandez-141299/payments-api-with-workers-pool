package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
)

type S3ClientImpl struct {
	client *s3.Client
}

func NewS3WithClient(client *s3.Client) S3Client {
	return &S3ClientImpl{client: client}
}

func NewS3Client() (S3Client, error) {
	if config.IsLocal() {
		// Return a mock implementation for local development
		return &S3ClientImpl{}, nil
	}

	s3Client, err := newS3Client(config.Config())
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %w", err)
	}

	return &S3ClientImpl{
		client: s3Client,
	}, nil
}

func (s *S3ClientImpl) PutObject(ctx context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return s.client.PutObject(ctx, params)
}

func newS3Client(cfg config.Configuration) (*s3.Client, error) {
	awsCfg, err := awsConf(cfg.Aws)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	return s3Client, nil
}

func awsConf(awsConfig config.Aws) (aws.Config, error) {
	ctx := context.Background()

	baseCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(awsConfig.Region),
	)

	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	return baseCfg, nil
}

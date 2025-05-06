package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type LocalstackContainer struct {
	testcontainers.Container
	Endpoint string
	Region   string
}

func SetupLocalstackContainer(t *testing.T) (*LocalstackContainer, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack:3.5.0",
		ExposedPorts: []string{"4566/tcp"},
		Env: map[string]string{
			"SERVICES": "s3",
		},
		WaitingFor: wait.ForListeningPort("4566/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "4566")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("http://%s:%s", hostIP, mappedPort.Port())

	return &LocalstackContainer{
		Container: container,
		Endpoint:  endpoint,
		Region:    "us-east-1",
	}, nil
}

func CreateS3Client(container *LocalstackContainer) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               container.Endpoint,
			SigningRegion:     container.Region,
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(container.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "test")),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return s3Client, nil
}

func CreateS3Bucket(ctx context.Context, client *s3.Client, bucketName string) error {
	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

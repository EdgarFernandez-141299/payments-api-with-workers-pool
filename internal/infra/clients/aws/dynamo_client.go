package aws

import (
	"context"
	"fmt"
	"log"
	"os"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"

	"gitlab.com/clubhub.ai1/gommon/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsDynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"gitlab.com/clubhub.ai1/gommon/nosql/dynamodb"
)

func GetLocalConfig() (aws.Config, error) {
	awsConfig := config.Config().Aws
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == awsDynamodb.ServiceID && region == awsConfig.Region {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               awsConfig.EndpointURL,
					SigningRegion:     awsConfig.Region,
					HostnameImmutable: true,
				}, nil
			}

			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(awsConfig.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsConfig.AccessKeyID,
			awsConfig.SecretAccessKey,
			"")),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return awsCfg, fmt.Errorf("failed to load LOCAL configuration, %w", err)
	}

	return awsCfg, nil
}

func NewDynamoClient(cfg config.Configuration, logger logger.LoggerInterface) *dynamodb.Client {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(), awsconfig.WithRegion(cfg.Aws.Region))
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	if e := os.Getenv("ENVIRONMENT"); e == "local" {
		awsCfg, err = GetLocalConfig()
		if err != nil {
			logger.Fatal(err)
		}
	}

	awsDynamoClient := awsDynamodb.NewFromConfig(awsCfg)

	if err != nil {
		logger.Fatal(err)
	}

	dynamodbClient, err := dynamodb.New(awsDynamoClient)

	if err != nil {
		logger.Fatal(err)
	}

	return dynamodbClient
}

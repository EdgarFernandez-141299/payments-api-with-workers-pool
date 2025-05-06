package storage

import (
	"context"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	aws2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/aws"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/clubhub.ai1/gommon/logger"
)

const (
	errBucketNameNotSet = "s3 bucket name not set in configuration"
	errUploadFailed     = "failed to upload file to S3: %w"
)

type S3StorageAdapter struct {
	client         aws2.S3Client
	bucketName     string
	logger         logger.LoggerInterface
	cdnURLProvider CDNURLProvider
}

func NewS3StorageAdapter(s3Client aws2.S3Client, bucket *PaymentReceiptBucket, logger logger.LoggerInterface, cdnURLProvider CDNURLProvider) (*S3StorageAdapter, error) {
	return &S3StorageAdapter{
		client:         s3Client,
		bucketName:     bucket.name,
		logger:         logger,
		cdnURLProvider: cdnURLProvider,
	}, nil
}

func (s *S3StorageAdapter) Store(oldCtx context.Context, path string, reader io.Reader) (string, error) {
	var fileURL string
	err := decorators.TraceDecoratorNoReturn(oldCtx, "S3StorageAdapter.Store", func(ctx context.Context, span decorators.Span) error {
		_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(path),
			Body:   reader,
		})
		if err != nil {
			return fmt.Errorf(errUploadFailed, err)
		}

		// Construct the URL using CDN URL
		fileURL = fmt.Sprintf("%s/%s", s.cdnURLProvider.GetCDNURL(), path)

		return nil
	})
	return fileURL, err
}

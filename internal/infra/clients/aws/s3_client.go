package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

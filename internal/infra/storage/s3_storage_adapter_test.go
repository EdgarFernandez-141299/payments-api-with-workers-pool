package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	aws2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/clients/aws"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/common/testutils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/common"
)

func createTestAdapter(client aws2.S3Client, bucketName string) *S3StorageAdapter {
	mockCDNURLProvider := NewMockCDNURLProvider("royal-cms-img.dev.clubhub.ai")
	adapter, _ := NewS3StorageAdapter(client, &PaymentReceiptBucket{name: bucketName}, nil, mockCDNURLProvider)
	return adapter
}

func TestS3StorageAdapter_Store(t *testing.T) {
	localstack, err := testutils.SetupLocalstackContainer(t)
	require.NoError(t, err)
	defer func() {
		if err := localstack.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	s3Client, err := testutils.CreateS3Client(localstack)
	require.NoError(t, err)

	bucketName := "test-bucket"
	err = testutils.CreateS3Bucket(context.Background(), s3Client, bucketName)
	require.NoError(t, err)

	client := aws2.NewS3WithClient(s3Client)

	tests := []struct {
		name          string
		bucketName    string
		path          string
		content       string
		client        aws2.S3Client
		expectedError bool
	}{
		{
			name:          "Success - File uploaded successfully (real client)",
			bucketName:    bucketName,
			path:          "test/path.txt",
			content:       "test content",
			client:        client,
			expectedError: false,
		},
		{
			name:       "Error - Failed to upload file (mock client)",
			bucketName: "test-bucket",
			path:       "test/path.txt",
			content:    "test content",
			client: func() aws2.S3Client {
				mockClient := common.NewS3Client(t)
				mockClient.On("PutObject", mock.Anything, mock.Anything).
					Return(nil, errors.New("upload failed"))
				return mockClient
			}(),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter := createTestAdapter(tt.client, tt.bucketName)

			reader := bytes.NewReader([]byte(tt.content))
			fileURL, err := adapter.Store(context.Background(), tt.path, reader)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to upload file to S3")
			} else {
				assert.NoError(t, err)
				// Verify the fileURL is constructed correctly using the CDN URL
				expectedURL := fmt.Sprintf("%s/%s", "royal-cms-img.dev.clubhub.ai", tt.path)
				assert.Equal(t, expectedURL, fileURL)

				if _, ok := tt.client.(*common.S3Client); !ok {
					output, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
						Bucket: aws.String(tt.bucketName),
						Key:    aws.String(tt.path),
					})
					require.NoError(t, err)
					defer output.Body.Close()

					content, err := io.ReadAll(output.Body)
					require.NoError(t, err)
					assert.Equal(t, tt.content, string(content))
				}
			}

			if mockClient, ok := tt.client.(*common.S3Client); ok {
				mockClient.AssertExpectations(t)
			}
		})
	}
}

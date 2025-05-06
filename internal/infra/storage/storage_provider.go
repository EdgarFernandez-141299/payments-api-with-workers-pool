package storage

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters"
)

func NewProvideStorageAdapter(s3Storage *S3StorageAdapter, diskStorage *DiskStorageAdapter) adapters.StorageAdapter {
	if config.IsLocal() {
		return diskStorage
	}
	return s3Storage
}

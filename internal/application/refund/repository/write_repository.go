package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
)

type RefundWriteRepositoryIF interface {
	Create(ctx context.Context, refund entities.RefundEntity) error
}

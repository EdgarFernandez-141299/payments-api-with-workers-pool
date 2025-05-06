package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
)

type RefundReadRepositoryIF interface {
	GetRefundsByReferenceOrderID(ctx context.Context, referenceOrderID, paymentID string, enterpriseID string) ([]entities.RefundEntity, error)
}

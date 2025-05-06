package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
)

type OrderReadRepositoryIF interface {
	GetOrderByReferenceID(ctx context.Context, referenceOrderID, enterpriseID string) (entities.OrderEntity, error)
	GetOrderPayments(ctx context.Context, referenceOrderID, enterpriseID string) ([]projections.OrderPaymentsProjection, error)
}

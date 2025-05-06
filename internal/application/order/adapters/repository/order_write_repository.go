package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
)

type OrderWriteRepositoryIF interface {
	CreateOrder(ctx context.Context, entity entities.OrderEntity) error
	UpdateOrder(ctx context.Context, entity entities.OrderEntity) error
}

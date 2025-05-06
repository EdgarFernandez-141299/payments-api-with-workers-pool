package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
)

type PaymentOrderRepositoryIF interface {
	CreatePaymentOrder(ctx context.Context, entity entities.PaymentOrderEntity) error
	UpdatePaymentOrder(ctx context.Context, entity entities.PaymentOrderEntity) error
}

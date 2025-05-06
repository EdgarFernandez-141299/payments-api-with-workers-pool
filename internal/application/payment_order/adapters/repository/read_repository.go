package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
)

type GetPaymentOrderByReferenceIF interface {
	GetPaymentOrderByReference(ctx context.Context, referenceID, orderID, enterpriseID string) (entities.PaymentOrderEntity, error)
}

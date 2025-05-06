package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type TransactionsRepositoryIF interface {
	UpdatePaymentOrderStatus(ctx context.Context, orderReferenceId, paymentId, enterpriseID string, status enums.PaymentStatus) error
}

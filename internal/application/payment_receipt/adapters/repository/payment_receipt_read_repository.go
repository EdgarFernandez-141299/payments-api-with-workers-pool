package repository

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
)

type PaymentReceiptRepository interface {
	GetByPaymentID(ctx context.Context, paymentID string) (entities.PaymentReceipt, error)
	CreatePaymentReceipt(ctx context.Context, cmd entities.PaymentReceipt) error
}

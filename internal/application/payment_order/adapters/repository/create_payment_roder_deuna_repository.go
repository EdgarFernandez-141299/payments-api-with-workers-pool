package repository

import (
	"context"
)

type DeunaOrderRepository interface {
	CreatePaymentOrderDeuna(ctx context.Context, paymentID, orderID, deunaOrderToken string) error
	GetTokenByOrderAndPaymentID(ctx context.Context, orderID string, paymentID string) (token string, err error)
}

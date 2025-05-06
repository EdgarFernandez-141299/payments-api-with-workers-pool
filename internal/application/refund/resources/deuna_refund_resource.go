package resources

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"
)

type DeunaRefundResourceIF interface {
	MakeTotalRefund(
		ctx context.Context,
		body utils.DeunaTotalRefundRequest,
		orderToken string,
	) (response.DeunaRefundPaymentResponse, error)

	MakePartialRefund(
		ctx context.Context,
		body utils.DeunaPartialRefundRequest,
		orderToken string,
	) (response.DeunaRefundPaymentResponse, error)
}

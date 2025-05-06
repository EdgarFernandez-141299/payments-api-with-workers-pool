package resources

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"
)

type DeunaPaymentResourceIF interface {
	MakeOrderPayment(
		ctx context.Context,
		body request.DeunaOrderPaymentRequest,
		token string,
	) error

	MakeOrderPaymentV2(
		ctx context.Context,
		body request.DeunaOrderPaymentRequestV2,
	) (response.DeunaOrderPaymentResponseV2, error)

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

package http

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"

	paymentResources "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"
	refundResources "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/resources"
)

type DeUnaRefundHTTPClient struct {
	paymentClient paymentResources.DeunaPaymentResourceIF
}

func NewDeUnaRefundHTTPClient(paymentClient paymentResources.DeunaPaymentResourceIF) refundResources.DeunaRefundResourceIF {
	return &DeUnaRefundHTTPClient{
		paymentClient: paymentClient,
	}
}

func (d *DeUnaRefundHTTPClient) MakeTotalRefund(
	ctx context.Context,
	body utils.DeunaTotalRefundRequest,
	orderToken string,
) (response.DeunaRefundPaymentResponse, error) {
	return d.paymentClient.MakeTotalRefund(ctx, body, orderToken)
}

func (d *DeUnaRefundHTTPClient) MakePartialRefund(
	ctx context.Context,
	body utils.DeunaPartialRefundRequest,
	orderToken string,
) (response.DeunaRefundPaymentResponse, error) {
	return d.paymentClient.MakePartialRefund(ctx, body, orderToken)
}

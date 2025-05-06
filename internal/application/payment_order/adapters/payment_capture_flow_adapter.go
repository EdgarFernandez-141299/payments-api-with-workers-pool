package adapter

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
)

type PaymentCaptureFlowAdapterIF interface {
	CapturePayment(ctx context.Context, orderID, paymentID string, paymentTotal decimal.Decimal) error
	ReleasePayment(ctx context.Context, orderID, paymentID, reason string) error
}

type PaymentCaptureFlowAdapter struct {
	deunaCaptureFlowResource resources.DeunaCaptureFlowResourceIF
	deunaOrderRepository     repository.DeunaOrderRepository
}

func NewPaymentCaptureFlowAdapter(
	deunaCaptureFlowResource resources.DeunaCaptureFlowResourceIF,
	deunaOrderRepository repository.DeunaOrderRepository,
) PaymentCaptureFlowAdapterIF {
	return &PaymentCaptureFlowAdapter{
		deunaCaptureFlowResource: deunaCaptureFlowResource,
		deunaOrderRepository:     deunaOrderRepository,
	}
}

func (a *PaymentCaptureFlowAdapter) CapturePayment(ctx context.Context, orderID, paymentID string, paymentTotal decimal.Decimal) error {
	return decorators.TraceDecoratorNoReturn(
		ctx,
		"PaymentCaptureFlowAdapter.CapturePayment",
		func(ctx context.Context, span decorators.Span) error {
			token, err := a.deunaOrderRepository.GetTokenByOrderAndPaymentID(ctx, orderID, paymentID)

			if err != nil {
				return err
			}

			captureSuccess, err := a.deunaCaptureFlowResource.Capture(ctx, token, utils.NewDeunaAmount(paymentTotal))

			if err != nil {
				return err
			}

			if !captureSuccess {
				return errors.New("capture failed")
			}

			return nil
		},
	)
}

func (a *PaymentCaptureFlowAdapter) ReleasePayment(ctx context.Context, orderID, paymentID, reason string) error {
	return decorators.TraceDecoratorNoReturn(
		ctx,
		"PaymentCaptureFlowAdapter.ReleasePayment",
		func(ctx context.Context, span decorators.Span) error {
			token, err := a.deunaOrderRepository.GetTokenByOrderAndPaymentID(ctx, orderID, paymentID)

			if err != nil {
				return err
			}

			releaseSuccess, err := a.deunaCaptureFlowResource.Release(ctx, token, reason)

			if err != nil {
				return err
			}

			if !releaseSuccess {
				return errors.New("release failed")
			}

			return nil
		},
	)
}

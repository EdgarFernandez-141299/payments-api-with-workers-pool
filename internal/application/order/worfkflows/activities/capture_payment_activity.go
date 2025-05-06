package activities

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/saga/errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/capture_flow"
)

const CapturePaymentActivityName = "CapturePaymentActivity"

type CapturePaymentActivity struct {
	capturePaymentUseCase capture_flow.PaymentCaptureUseCaseIF
}

func NewCapturePaymentActivity(capturePaymentUseCase capture_flow.PaymentCaptureUseCaseIF) *CapturePaymentActivity {
	return &CapturePaymentActivity{
		capturePaymentUseCase: capturePaymentUseCase,
	}
}

func (c *CapturePaymentActivity) CapturePayment(ctx context.Context, OrderID, PaymentID string) error {
	return errors.WrapActivityError(c.capturePaymentUseCase.CapturePayment(ctx, OrderID, PaymentID))
}

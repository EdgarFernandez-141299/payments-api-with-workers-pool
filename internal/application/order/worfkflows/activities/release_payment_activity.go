package activities

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/saga/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/capture_flow"
)

const ReleasePaymentActivityName = "ReleasePaymentActivity"

type ReleasePaymentActivity struct {
	releasePaymentUseCase capture_flow.PaymentReleaseUseCaseIF
}

func NewReleasePaymentActivity(releasePaymentUseCase capture_flow.PaymentReleaseUseCaseIF) *ReleasePaymentActivity {
	return &ReleasePaymentActivity{
		releasePaymentUseCase: releasePaymentUseCase,
	}
}

func (r *ReleasePaymentActivity) ReleasePayment(ctx context.Context, orderID, paymentID, reason string) error {
	return errors.WrapActivityError(r.releasePaymentUseCase.ReleasePayment(ctx, orderID, paymentID, reason))
}

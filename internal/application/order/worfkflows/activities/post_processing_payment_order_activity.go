package activities

import (
	"context"

	postpayment "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/post_payment"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

const PostProcessingPaymentOrderActivityName = "PostProcessingPaymentOrderActivity"

type PostProcessingPaymentOrderActivity struct {
	useCase postpayment.PostProcessingPaymentOrderUseCaseIF
}

func NewPostProcessingPaymentOrderActivity(useCase postpayment.PostProcessingPaymentOrderUseCaseIF) *PostProcessingPaymentOrderActivity {
	return &PostProcessingPaymentOrderActivity{
		useCase: useCase,
	}
}

func (a *PostProcessingPaymentOrderActivity) PostProcessingPaymentOrder(ctx context.Context, cmd postpayment.PostProcessingPaymentOrderCommand) (enums.PaymentFlowEnum, error) {
	return a.useCase.PostProcessPaymentOrder(ctx, cmd)
}

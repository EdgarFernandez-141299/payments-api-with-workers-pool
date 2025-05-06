package activities

import (
	"context"
	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/go-libraries/saga/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/use_cases/create"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

const CreatePaymentOrderActivityName = "CreatePaymentOrderActivity"

type CreatePaymentOrderActivity struct {
	useCase create.CreatePaymentOrderUseCaseIF
}

func NewCreatePaymentOrderActivity(useCase create.CreatePaymentOrderUseCaseIF) *CreatePaymentOrderActivity {
	return &CreatePaymentOrderActivity{
		useCase: useCase,
	}
}

func (a *CreatePaymentOrderActivity) CreatePaymentOrder(ctx context.Context, cmd command.CreatePaymentOrderCommand) (response.PaymentOrderResponseDTO, error) {
	return decorators.TraceDecorator[response.PaymentOrderResponseDTO](
		ctx,
		"CreatePaymentOrderActivity.CreatePaymentOrder",
		func(ctx context.Context, span decorators.Span) (response.PaymentOrderResponseDTO, error) {
			payment, err := a.useCase.CreatePaymentOrder(ctx, cmd)
			return payment, errors.WrapActivityError(err)
		})
}

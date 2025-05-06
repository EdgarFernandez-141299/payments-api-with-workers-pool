package capture_flow

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	adapter "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type PaymentReleaseUseCaseIF interface {
	ReleasePayment(ctx context.Context, orderID, paymentID, reason string) error
}

type PaymentReleaseUseCase struct {
	repository                event_store.OrderEventRepository
	paymentCaptureFlowAdapter adapter.PaymentCaptureFlowAdapterIF
	transactionsRepository    repository.TransactionsRepositoryIF
}

func NewPaymentReleaseUseCase(
	repository event_store.OrderEventRepository,
	paymentCaptureFlowAdapter adapter.PaymentCaptureFlowAdapterIF,
	transactionsRepository repository.TransactionsRepositoryIF,
) PaymentReleaseUseCaseIF {
	return &PaymentReleaseUseCase{
		repository:                repository,
		paymentCaptureFlowAdapter: paymentCaptureFlowAdapter,
		transactionsRepository:    transactionsRepository,
	}
}

func (u *PaymentReleaseUseCase) ReleasePayment(ctx context.Context, orderID, paymentID, reason string) error {
	return decorators.TraceDecoratorNoReturn(
		ctx,
		"PaymentReleaseUseCase.ReleasePayment",
		func(ctx context.Context, span decorators.Span) error {
			order := new(aggregate.Order)
			err := u.repository.Get(ctx, orderID, order)

			if err != nil {
				return err
			}

			err = u.paymentCaptureFlowAdapter.ReleasePayment(ctx, orderID, paymentID, reason)

			if err != nil {
				return err
			}

			order.OrderPaymentReleased(command.NewPaymentOrderReleasedCommand(
				orderID,
				paymentID,
				reason,
			))

			err = u.transactionsRepository.UpdatePaymentOrderStatus(ctx, orderID, paymentID, order.EnterpriseID, enums.PaymentStatus(order.Status.Get()))
			if err != nil {
				return err
			}

			err = u.repository.Save(ctx, order)

			if err != nil {
				return err
			}

			return nil
		},
	)
}

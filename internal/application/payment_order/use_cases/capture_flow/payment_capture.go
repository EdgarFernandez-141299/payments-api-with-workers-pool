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

type PaymentCaptureUseCaseIF interface {
	CapturePayment(ctx context.Context, orderID, paymentID string) error
}

type PaymentCaptureUseCase struct {
	repository                event_store.OrderEventRepository
	deunaOrderRepository      repository.DeunaOrderRepository
	paymentCaptureFlowAdapter adapter.PaymentCaptureFlowAdapterIF
	transactionsRepository    repository.TransactionsRepositoryIF
}

func NewPaymentCaptureUseCase(
	repository event_store.OrderEventRepository,
	deunaOrderRepository repository.DeunaOrderRepository,
	paymentCaptureFlowAdapter adapter.PaymentCaptureFlowAdapterIF,
	transactionsRepository repository.TransactionsRepositoryIF,
) PaymentCaptureUseCaseIF {
	return &PaymentCaptureUseCase{
		repository:                repository,
		deunaOrderRepository:      deunaOrderRepository,
		paymentCaptureFlowAdapter: paymentCaptureFlowAdapter,
		transactionsRepository:    transactionsRepository,
	}
}

func (u *PaymentCaptureUseCase) CapturePayment(ctx context.Context, orderID, paymentID string) error {
	return decorators.TraceDecoratorNoReturn(
		ctx,
		"PaymentCaptureUseCase.CapturePayment",
		func(ctx context.Context, span decorators.Span) error {
			order := new(aggregate.Order)
			err := u.repository.Get(ctx, orderID, order)
			if err != nil {
				return err
			}

			payment, err := order.FindPaymentByID(paymentID)

			if err != nil {
				return err
			}

			err = u.paymentCaptureFlowAdapter.CapturePayment(ctx, orderID, paymentID, payment.Total.Value)

			if err != nil {
				return err
			}

			order.OrderPaymentCaptured(command.NewPaymentOrderCapturedCommand(
				orderID,
				paymentID,
				payment.Total.Code.Code,
				payment.Total.Value,
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

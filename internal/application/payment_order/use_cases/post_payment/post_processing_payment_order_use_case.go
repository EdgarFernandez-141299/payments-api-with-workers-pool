package use_cases

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	paymentOrderRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type PostProcessingPaymentOrderUseCaseIF interface {
	PostProcessPaymentOrder(ctx context.Context, cmd PostProcessingPaymentOrderCommand) (enums.PaymentFlowEnum, error)
}

type PostProcessingPaymentOrderUseCase struct {
	orderRepository             event_store.OrderEventRepository
	orderReadRepository         repository.OrderReadRepositoryIF
	orderWriteRepository        repository.OrderWriteRepositoryIF
	paymentOrderWriteRepository paymentOrderRepository.PaymentOrderRepositoryIF
	paymentOrderReadRepository  paymentOrderRepository.GetPaymentOrderByReferenceIF
}

func NewPostProcessingPaymentOrderUseCase(
	orderRepository event_store.OrderEventRepository,
	orderReadRepository repository.OrderReadRepositoryIF,
	orderWriteRepository repository.OrderWriteRepositoryIF,
	paymentOrderWriteRepository paymentOrderRepository.PaymentOrderRepositoryIF,
	paymentOrderReadRepository paymentOrderRepository.GetPaymentOrderByReferenceIF,
) PostProcessingPaymentOrderUseCaseIF {
	return &PostProcessingPaymentOrderUseCase{
		orderRepository:             orderRepository,
		orderReadRepository:         orderReadRepository,
		orderWriteRepository:        orderWriteRepository,
		paymentOrderWriteRepository: paymentOrderWriteRepository,
		paymentOrderReadRepository:  paymentOrderReadRepository,
	}
}

type CardData struct {
	CardBrand string
	CardLast4 string
	CardType  string
}
type PostProcessingPaymentOrderCommand struct {
	ReferenceOrderID  string
	PaymentID         string
	Status            enums.PaymentStatus
	AuthorizationCode string
	OrderStatusString string
	PaymentReason     string
	PaymentCard       CardData
}

func (c *PostProcessingPaymentOrderUseCase) PostProcessPaymentOrder(ctx context.Context,
	cmd PostProcessingPaymentOrderCommand) (enums.PaymentFlowEnum, error) {
	order := new(aggregate.Order)
	err := c.orderRepository.Get(ctx, cmd.ReferenceOrderID, order)

	if err != nil {
		return enums.Autocapture, err
	}

	if cmd.Status.IsFailure() {
		order.OrderPaymentFailed(command.NewCreatePaymentOrderFailCommand(
			cmd.ReferenceOrderID,
			cmd.PaymentID,
			cmd.PaymentReason,
			cmd.Status.String(),
		))
	} else if cmd.Status.IsAuthorized() {
		order.OrderPaymentAuthorized(command.NewCreatePaymentOrderAuthorizedCommand(
			cmd.ReferenceOrderID,
			cmd.PaymentID,
			cmd.AuthorizationCode,
			cmd.Status.String(),
			command.CardData{
				CardBrand: cmd.PaymentCard.CardBrand,
				CardLast4: cmd.PaymentCard.CardLast4,
				CardType:  cmd.PaymentCard.CardType,
			},
		))
	} else {
		order.OrderPaymentProcessed(command.NewCreatePaymentOrderProcessedCommand(
			cmd.ReferenceOrderID,
			cmd.PaymentID,
			cmd.AuthorizationCode,
			cmd.Status.String(),
			command.CardData{
				CardBrand: cmd.PaymentCard.CardBrand,
				CardLast4: cmd.PaymentCard.CardLast4,
				CardType:  cmd.PaymentCard.CardType,
			},
		))
	}

	paymentOrder, err := c.paymentOrderReadRepository.GetPaymentOrderByReference(ctx, cmd.ReferenceOrderID,
		cmd.PaymentID, order.EnterpriseID)
	if err != nil {
		return enums.Autocapture, err
	}

	paymentOrder.SetStatus(enums.PaymentStatus(order.Status.Get()))
	paymentOrder.SetAuthorizationCode(cmd.AuthorizationCode)

	err = c.paymentOrderWriteRepository.UpdatePaymentOrder(ctx, paymentOrder)
	if err != nil {
		return enums.Autocapture, err
	}

	orderEntity, err := c.orderReadRepository.GetOrderByReferenceID(ctx, cmd.ReferenceOrderID, order.EnterpriseID)
	if err != nil {
		return enums.Autocapture, err
	}

	orderEntity.SetStatus(order.Status.Get())

	err = c.orderWriteRepository.UpdateOrder(ctx, orderEntity)

	if err != nil {
		return enums.Autocapture, err
	}

	err = c.orderRepository.Save(ctx, order)

	if err != nil {
		return enums.Autocapture, err
	}

	return enums.PaymentFlowEnum(paymentOrder.PaymentFlow), nil
}

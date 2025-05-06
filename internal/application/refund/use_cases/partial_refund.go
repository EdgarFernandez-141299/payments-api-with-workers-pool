package use_cases

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	paymentOrderRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	adapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/adapters/models"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

type PartialRefundUseCaseIF interface {
	PartialRefund(ctx context.Context, refund command.CreatePartialPaymentRefundCommand, enterpriseID string) (response.RefundResponseDTO, error)
}

type PartialRefundUse struct {
	repository                  event_store.OrderEventRepository
	partialRefundAdapter        adapters.PartialRefundAdapterIF
	paymentOrderReadRepository  paymentOrderRepository.GetPaymentOrderByReferenceIF
	paymentOrderWriteRepository paymentOrderRepository.PaymentOrderRepositoryIF
	orderReadRepository         repository.OrderReadRepositoryIF
	orderWriteRepository        repository.OrderWriteRepositoryIF
}

func NewPartialRefundUse(repository event_store.OrderEventRepository,
	partialRefundAdapter adapters.PartialRefundAdapterIF,
	paymentOrderReadRepository paymentOrderRepository.GetPaymentOrderByReferenceIF,
	paymentOrderWriteRepository paymentOrderRepository.PaymentOrderRepositoryIF,
	orderReadRepository repository.OrderReadRepositoryIF,
	orderWriteRepository repository.OrderWriteRepositoryIF,
) PartialRefundUseCaseIF {
	return &PartialRefundUse{
		repository:                  repository,
		partialRefundAdapter:        partialRefundAdapter,
		paymentOrderReadRepository:  paymentOrderReadRepository,
		paymentOrderWriteRepository: paymentOrderWriteRepository,
		orderReadRepository:         orderReadRepository,
		orderWriteRepository:        orderWriteRepository,
	}
}

func (u *PartialRefundUse) PartialRefund(ctx context.Context, refund command.CreatePartialPaymentRefundCommand, enterpriseID string) (response.RefundResponseDTO, error) {
	order := new(aggregate.Order)

	err := u.repository.Get(ctx, refund.ReferenceOrderID, order)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	validationErr := order.RefundPartialPayment(refund.PaymentOrderID, refund.Reason, refund.Amount)

	if validationErr != nil {
		return response.RefundResponseDTO{}, validationErr
	}

	refundModel := models.NewRefundModel(refund.Amount, refund.PaymentOrderID, order.ID)

	refundResponse, err := u.partialRefundAdapter.PartialRefund(ctx, refundModel, enterpriseID)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	paymentOrder, err := u.paymentOrderReadRepository.GetPaymentOrderByReference(ctx, refund.ReferenceOrderID,
		refund.PaymentOrderID, enterpriseID)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	paymentOrder.SetStatus(enums.PaymentStatus(order.Status.Get()))

	err = u.paymentOrderWriteRepository.UpdatePaymentOrder(ctx, paymentOrder)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	orderEntity, err := u.orderReadRepository.GetOrderByReferenceID(ctx, refund.ReferenceOrderID, enterpriseID)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	orderEntity.SetStatus(order.Status.Get())

	err = u.orderWriteRepository.UpdateOrder(ctx, orderEntity)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	err = u.repository.Save(ctx, order)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	return refundResponse, nil
}

package use_cases

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/event_store"
	paymentOrderRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
)

type RefundTotalUseCaseIF interface {
	Refund(ctx context.Context, refund command.RefundTotalCommand, enterpriseID string) (response.RefundResponseDTO, error)
}

type RefundTotalUse struct {
	eventRepository             event_store.OrderEventRepository
	refundAdapter               adapters.RefundAdapterIF
	paymentOrderReadRepository  paymentOrderRepository.GetPaymentOrderByReferenceIF
	paymentOrderWriteRepository paymentOrderRepository.PaymentOrderRepositoryIF
	orderReadRepository         repository.OrderReadRepositoryIF
	orderWriteRepository        repository.OrderWriteRepositoryIF
}

func NewRefundTotalUseCase(
	repository event_store.OrderEventRepository,
	refundAdapter adapters.RefundAdapterIF,
	paymentOrderReadRepository paymentOrderRepository.GetPaymentOrderByReferenceIF,
	paymentOrderWriteRepository paymentOrderRepository.PaymentOrderRepositoryIF,
	orderReadRepository repository.OrderReadRepositoryIF,
	orderWriteRepository repository.OrderWriteRepositoryIF,
) RefundTotalUseCaseIF {
	return &RefundTotalUse{
		eventRepository:             repository,
		refundAdapter:               refundAdapter,
		paymentOrderReadRepository:  paymentOrderReadRepository,
		paymentOrderWriteRepository: paymentOrderWriteRepository,
		orderReadRepository:         orderReadRepository,
		orderWriteRepository:        orderWriteRepository,
	}
}

func (r *RefundTotalUse) Refund(
	ctx context.Context,
	refund command.RefundTotalCommand, enterpriseID string,
) (response.RefundResponseDTO, error) {
	order := new(aggregate.Order)

	err := r.eventRepository.Get(ctx, refund.ReferenceOrderID, order)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	_, err = order.RefundPayment(refund.PaymentOrderID, refund.Reason)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	refundResponse, err := r.refundAdapter.RefundPayment(
		ctx, refund.PaymentOrderID, refund.ReferenceOrderID, enterpriseID, refund.Reason,
	)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	paymentOrder, err := r.paymentOrderReadRepository.GetPaymentOrderByReference(ctx, refund.ReferenceOrderID,
		refund.PaymentOrderID, enterpriseID)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	paymentOrder.SetStatus(enums.PaymentRefunded)

	err = r.paymentOrderWriteRepository.UpdatePaymentOrder(ctx, paymentOrder)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	orderEntity, err := r.orderReadRepository.GetOrderByReferenceID(ctx, refund.ReferenceOrderID, enterpriseID)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	orderEntity.SetStatus(enums.PaymentRefunded.String())

	err = r.orderWriteRepository.UpdateOrder(ctx, orderEntity)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	err = r.eventRepository.Save(ctx, order)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	return refundResponse, nil
}

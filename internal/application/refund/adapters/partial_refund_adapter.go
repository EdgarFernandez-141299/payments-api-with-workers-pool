package adapters

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/adapters/models"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"

	repositoryPaymentOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
)

type PartialRefundAdapterIF interface {
	PartialRefund(ctx context.Context,
		refund models.RefundAdapterModel,
		refundReason string) (response.RefundResponseDTO, error)
}

type PartialRefundAdapter struct {
	deunaRefundResource  resources.DeunaRefundResourceIF
	writeRepository      repository.RefundWriteRepositoryIF
	deunaOrderRepository repositoryPaymentOrder.DeunaOrderRepository
}

func NewPartialRefundAdapter(
	deunaRefundResource resources.DeunaRefundResourceIF,
	writeRepository repository.RefundWriteRepositoryIF,
	deunaOrderRepo repositoryPaymentOrder.DeunaOrderRepository,
) PartialRefundAdapterIF {
	return &PartialRefundAdapter{
		deunaRefundResource:  deunaRefundResource,
		writeRepository:      writeRepository,
		deunaOrderRepository: deunaOrderRepo,
	}
}

func (a *PartialRefundAdapter) PartialRefund(
	ctx context.Context,
	refund models.RefundAdapterModel,
	refundReason string,
) (response.RefundResponseDTO, error) {
	refundAmount := refund.Amount
	payload := utils.DeunaPartialRefundRequest{
		Amount: utils.NewDeunaAmount(refundAmount),
		Reason: refundReason,
	}

	orderToken, err := a.deunaOrderRepository.GetTokenByOrderAndPaymentID(
		ctx, refund.OrderID, refund.PaymentID,
	)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	_, err = a.deunaRefundResource.MakePartialRefund(ctx, payload, orderToken)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	refundEntity := entities.NewRefundEntityBuilder().
		WithPaymentID(refund.PaymentID).
		WithOrderID(refund.OrderID).
		WithAmount(refund.Amount).
		WithReason(refundReason).
		WithStatus(enums.PartiallyRefunded.String()).
		Build()

	err = a.writeRepository.Create(ctx, *refundEntity)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	return response.RefundResponseDTO{
		ReferenceOrderID: refund.OrderID,
		PaymentOrderID:   refund.PaymentID,
		Amount:           refund.Amount,
	}, nil
}

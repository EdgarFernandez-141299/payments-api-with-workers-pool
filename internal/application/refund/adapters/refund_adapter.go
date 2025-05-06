package adapters

import (
	"context"
	"fmt"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"

	"github.com/shopspring/decimal"
	repositoryPaymentOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	refundEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
)

type RefundAdapterIF interface {
	RefundPayment(ctx context.Context,
		paymentID string,
		orderID string,
		enterpriseID string,
		reason string,
	) (response.RefundResponseDTO, error)
}

type RefundAdapter struct {
	deunaRefundResource  resources.DeunaRefundResourceIF
	writeRepository      repository.RefundWriteRepositoryIF
	deunaOrderRepository repositoryPaymentOrder.DeunaOrderRepository
}

func NewRefundAdapter(
	deunaRefundResource resources.DeunaRefundResourceIF,
	writeRepository repository.RefundWriteRepositoryIF,
	deunaOrderRepository repositoryPaymentOrder.DeunaOrderRepository,
) RefundAdapterIF {
	return &RefundAdapter{
		deunaRefundResource:  deunaRefundResource,
		writeRepository:      writeRepository,
		deunaOrderRepository: deunaOrderRepository,
	}
}

func (r *RefundAdapter) RefundPayment(ctx context.Context,
	paymentID string,
	orderID string,
	enterpriseID string,
	reason string,
) (response.RefundResponseDTO, error) {
	orderToken, err := r.deunaOrderRepository.GetTokenByOrderAndPaymentID(
		ctx, orderID, paymentID,
	)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	refundResponse, err := r.deunaRefundResource.MakeTotalRefund(
		ctx,
		utils.DeunaTotalRefundRequest{
			Reason: reason,
		},
		orderToken,
	)

	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	amountRefund, err := decimal.NewFromString(refundResponse.Data.RefundAmount.Amount)

	if err != nil {
		return response.RefundResponseDTO{}, fmt.Errorf("failed to parse refund amount: %w", err)
	}

	// save refund in the database
	refundEntity := refundEntities.NewRefundEntityBuilder().
		WithPaymentID(paymentID).
		WithOrderID(orderID).
		WithEnterpriseID(enterpriseID).
		WithReason(reason).
		WithStatus(refundResponse.Data.Status).
		WithAmount(utils.DeunaAmountToAmount(amountRefund.IntPart())).
		Build()

	err = r.writeRepository.Create(ctx, *refundEntity)
	if err != nil {
		return response.RefundResponseDTO{}, err
	}

	return response.RefundResponseDTO{
		ReferenceOrderID: orderID,
		PaymentOrderID:   paymentID,
		Amount:           utils.DeunaAmountToAmount(amountRefund.IntPart()),
	}, nil
}

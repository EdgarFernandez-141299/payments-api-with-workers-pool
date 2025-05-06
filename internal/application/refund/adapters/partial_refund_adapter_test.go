package adapters

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"
	adapterUtils "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/adapters/models"
	paymentOrderUtils "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

func TestPartialRefundAdapter(t *testing.T) {
	refundModel := models.NewRefundModel(
		decimal.NewFromFloat(50.00),
		"paymentOrderID",
		"referenceOrderID",
	)
	orderToken := "orderToken"
	enterpriseID := "enterpriseID"
	refundReason := "refundReason"

	payload := adapterUtils.DeunaPartialRefundRequest{
		Amount: adapterUtils.NewDeunaAmount(refundModel.Amount),
		Reason: refundReason,
	}

	t.Run("should return an error getting payment order", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			refundModel.OrderID,
			refundModel.PaymentID,
		).
			Return("", assert.AnError)

		partialRefundAdapter := NewPartialRefundAdapter(
			mockDeunaRefundResource,
			mockWriteRepository,
			deunaOrderRepository,
		)

		_, err := partialRefundAdapter.PartialRefund(context.TODO(), refundModel, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		deunaOrderRepository.AssertExpectations(t)
		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
	})

	t.Run("should fail when external deuna integration fails", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			refundModel.OrderID,
			refundModel.PaymentID,
		).
			Return(orderToken, nil)

		mockDeunaRefundResource.On("MakePartialRefund", mock.Anything, payload, orderToken).
			Return(response.DeunaRefundPaymentResponse{}, assert.AnError)

		partialRefundAdapter := NewPartialRefundAdapter(
			mockDeunaRefundResource,
			mockWriteRepository,
			deunaOrderRepository,
		)

		_, err := partialRefundAdapter.PartialRefund(context.TODO(), refundModel, refundReason)

		assert.Error(t, err, assert.AnError.Error())

		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
		deunaOrderRepository.AssertExpectations(t)
	})

	t.Run("should return an error when make partial refund deuna fails", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			refundModel.OrderID,
			refundModel.PaymentID,
		).
			Return(orderToken, nil)

		mockDeunaRefundResource.On("MakePartialRefund", mock.Anything, payload, orderToken).
			Return(response.DeunaRefundPaymentResponse{
				Data: response.RefundData{RefundID: orderToken},
			}, nil)

		partialRefundAdapter := NewPartialRefundAdapter(
			mockDeunaRefundResource,
			mockWriteRepository,
			deunaOrderRepository,
		)

		refundEntity := entities.NewRefundEntityBuilder().
			WithPaymentID(refundModel.PaymentID).
			WithOrderID(refundModel.OrderID).
			WithAmount(refundModel.Amount).
			WithReason(payload.Reason).
			Build()

		mockWriteRepository.On("Create", mock.Anything,
			mock.MatchedBy(func(entity entities.RefundEntity) bool {
				expectedPaymentID := paymentOrderUtils.GeneratePaymentOrderID(refundModel.OrderID, refundModel.PaymentID)
				return entity.PaymentID == expectedPaymentID &&
					entity.Amount.Equal(refundEntity.Amount) &&
					entity.Reason == refundEntity.Reason
			})).Return(assert.AnError)

		_, err := partialRefundAdapter.PartialRefund(context.TODO(), refundModel, refundReason)

		assert.Error(t, err, assert.AnError.Error())
		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
		deunaOrderRepository.AssertExpectations(t)
	})

	t.Run("should success when make partial refund deuna fails", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			refundModel.OrderID,
			refundModel.PaymentID,
		).
			Return(orderToken, nil)

		mockDeunaRefundResource.On("MakePartialRefund", mock.Anything, payload, orderToken).
			Return(response.DeunaRefundPaymentResponse{
				Data: response.RefundData{RefundID: orderToken},
			}, nil)

		partialRefundAdapter := NewPartialRefundAdapter(
			mockDeunaRefundResource,
			mockWriteRepository,
			deunaOrderRepository,
		)

		refundEntity := entities.NewRefundEntityBuilder().
			WithPaymentID(refundModel.PaymentID).
			WithOrderID(refundModel.OrderID).
			WithAmount(refundModel.Amount).
			WithReason(payload.Reason).
			Build()

		mockWriteRepository.On("Create", mock.Anything,
			mock.MatchedBy(func(entity entities.RefundEntity) bool {
				expectedPaymentID := paymentOrderUtils.GeneratePaymentOrderID(refundModel.OrderID, refundModel.PaymentID)
				return entity.PaymentID == expectedPaymentID &&
					entity.Amount.Equal(refundEntity.Amount) &&
					entity.Reason == refundEntity.Reason
			})).Return(nil)

		_, err := partialRefundAdapter.PartialRefund(context.TODO(), refundModel, refundReason)

		assert.NoError(t, err)
		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
		deunaOrderRepository.AssertExpectations(t)
	})
}

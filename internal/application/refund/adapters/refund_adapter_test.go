package adapters

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"
	adapterUtils "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	paymentOrderUtils "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

var (
	ctx   = context.Background()
	id, _ = uid.NewUniqueID(uid.WithID("id"))
)

func TestRefundAdapter(t *testing.T) {
	paymentID := "paymentID"
	orderID := "orderID"
	refundReason := "refundReason"
	enterpriseID := "enterpriseID"

	t.Run("should return an error getting payment order", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			orderID,
			paymentID,
		).
			Return("", assert.AnError)

		refundAdapter := NewRefundAdapter(mockDeunaRefundResource, mockWriteRepository, deunaOrderRepository)

		_, err := refundAdapter.RefundPayment(ctx, paymentID, orderID, enterpriseID, refundReason)

		assert.EqualError(t, err, assert.AnError.Error())
		deunaOrderRepository.AssertExpectations(t)
		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
	})

	t.Run("should return an error when make total refund fails", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		token := "token"
		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			orderID,
			paymentID,
		).Return(token, nil)

		mockDeunaRefundResource.On("MakeTotalRefund", mock.Anything, mock.Anything, token).
			Return(response.DeunaRefundPaymentResponse{}, assert.AnError)

		refundAdapter := NewRefundAdapter(mockDeunaRefundResource, mockWriteRepository, deunaOrderRepository)

		_, err := refundAdapter.RefundPayment(ctx, paymentID, orderID, enterpriseID, refundReason)

		assert.EqualError(t, err, assert.AnError.Error())
		deunaOrderRepository.AssertExpectations(t)
		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
	})

	t.Run("should return refund model when make total refund succeeds", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)
		deunaAmount := "1000"

		token := "token"
		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			orderID,
			paymentID,
		).Return(token, nil)

		refundResponse := response.DeunaRefundPaymentResponse{
			Data: response.RefundData{
				RefundAmount: response.RefundMonetaryValue{Amount: deunaAmount, Currency: "USD"},
				RefundID:     "refundID",
				Status:       "success",
			},
		}

		mockDeunaRefundResource.On("MakeTotalRefund", mock.Anything, mock.Anything, token).
			Return(refundResponse, nil)

		mockWriteRepository.On("Create", mock.Anything,
			mock.MatchedBy(func(refund entities.RefundEntity) bool {
				expectedPaymentID := paymentOrderUtils.GeneratePaymentOrderID(orderID, paymentID)
				return refund.PaymentID == expectedPaymentID &&
					refund.EnterpriseID == enterpriseID &&
					refund.Amount.Equal(adapterUtils.DeunaAmountToAmount(1000)) &&
					refund.Status == refundResponse.Data.Status
			})).Return(nil)

		refundAdapter := NewRefundAdapter(mockDeunaRefundResource, mockWriteRepository, deunaOrderRepository)

		_, err := refundAdapter.RefundPayment(ctx, paymentID, orderID, enterpriseID, refundReason)

		assert.NoError(t, err)
		deunaOrderRepository.AssertExpectations(t)
		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
	})
	t.Run("should return refund model when make total refund fails", func(t *testing.T) {
		mockDeunaRefundResource := resources.NewDeunaRefundResourceIF(t)
		mockWriteRepository := repository.NewRefundWriteRepositoryIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		token := "token"
		deunaOrderRepository.On("GetTokenByOrderAndPaymentID", mock.Anything,
			orderID,
			paymentID,
		).Return(token, nil)

		refundResponse := response.DeunaRefundPaymentResponse{
			Data: response.RefundData{
				RefundAmount: response.RefundMonetaryValue{Amount: "1000", Currency: "USD"},
				RefundID:     "refundID",
				Status:       "failed",
			},
		}

		mockDeunaRefundResource.On("MakeTotalRefund", mock.Anything, mock.Anything, token).
			Return(refundResponse, nil)

		mockWriteRepository.On("Create", mock.Anything,
			mock.MatchedBy(func(refund entities.RefundEntity) bool {
				expectedPaymentID := paymentOrderUtils.GeneratePaymentOrderID(orderID, paymentID)
				return refund.PaymentID == expectedPaymentID &&
					refund.EnterpriseID == enterpriseID &&
					refund.Amount.Equal(adapterUtils.DeunaAmountToAmount(1000)) &&
					refund.Status == refundResponse.Data.Status
			})).Return(assert.AnError)

		refundAdapter := NewRefundAdapter(mockDeunaRefundResource, mockWriteRepository, deunaOrderRepository)

		_, err := refundAdapter.RefundPayment(ctx, paymentID, orderID, enterpriseID, refundReason)

		assert.Error(t, err, assert.AnError.Error())
		deunaOrderRepository.AssertExpectations(t)
		mockDeunaRefundResource.AssertExpectations(t)
		mockWriteRepository.AssertExpectations(t)
	})
}

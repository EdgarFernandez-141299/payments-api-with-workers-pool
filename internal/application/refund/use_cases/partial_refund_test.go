package use_cases

import (
	"context"
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/adapters/models"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	order_entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	payment_order_entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository/fixture"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPartialRefundUseCase(t *testing.T) {
	ctx := context.Background()
	refundCommand := command.CreatePartialPaymentRefundCommand{
		ReferenceOrderID: "referenceOrderID",
		PaymentOrderID:   "paymentOrderID",
		Amount:           decimal.NewFromFloat(50.00),
		Reason:           "Customer request",
	}

	enterpriseID := "enterpriseID"

	t.Run("should return an error getting order", func(t *testing.T) {
		mockPartialRefundAdapter := adapters.NewPartialRefundAdapterIF(t)
		mockPaymentOrderReadRepo := repository.NewGetPaymentOrderByReferenceIF(t)
		mockPaymentOrderWriteRepo := repository.NewPaymentOrderRepositoryIF(t)
		mockOrderReadRepo := repository.NewOrderReadRepositoryIF(t)
		mockOrderWriteRepo := repository.NewOrderWriteRepositoryIF(t)

		currencyCode, _ := value_objects.NewCurrencyCode("USD")
		totalValue := decimal.NewFromFloat(100.01)
		currencyAmount, _ := value_objects.NewCurrencyAmount(currencyCode, totalValue)

		mockRepo := fixture.OrderEventRepositoryGetFixture(
			t,
			refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID,
			enums.PaymentProcessed,
			currencyAmount,
			currencyAmount,
			assert.AnError,
		)

		partialRefundUseCase := NewPartialRefundUse(
			mockRepo,
			mockPartialRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := partialRefundUseCase.PartialRefund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockPartialRefundAdapter.AssertExpectations(t)
	})

	t.Run("should return an error getting payment order", func(t *testing.T) {
		mockPartialRefundAdapter := adapters.NewPartialRefundAdapterIF(t)
		mockPaymentOrderReadRepo := repository.NewGetPaymentOrderByReferenceIF(t)
		mockPaymentOrderWriteRepo := repository.NewPaymentOrderRepositoryIF(t)
		mockOrderReadRepo := repository.NewOrderReadRepositoryIF(t)
		mockOrderWriteRepo := repository.NewOrderWriteRepositoryIF(t)

		currencyCode, _ := value_objects.NewCurrencyCode("USD")
		totalValue := decimal.NewFromFloat(100.01)
		currencyAmount, _ := value_objects.NewCurrencyAmount(currencyCode, totalValue)

		mockRepo := fixture.OrderEventRepositoryGetFixture(
			t,
			refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID,
			enums.PaymentProcessed,
			currencyAmount,
			currencyAmount,
			nil,
		)

		refundModel := models.NewRefundModel(refundCommand.Amount, refundCommand.PaymentOrderID, refundCommand.ReferenceOrderID)

		mockPartialRefundAdapter.On("PartialRefund", mock.Anything, refundModel, enterpriseID).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, assert.AnError)

		partialRefundUseCase := NewPartialRefundUse(
			mockRepo,
			mockPartialRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := partialRefundUseCase.PartialRefund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockPartialRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
	})

	t.Run("should return an error updating payment order", func(t *testing.T) {
		mockPartialRefundAdapter := adapters.NewPartialRefundAdapterIF(t)
		mockPaymentOrderReadRepo := repository.NewGetPaymentOrderByReferenceIF(t)
		mockPaymentOrderWriteRepo := repository.NewPaymentOrderRepositoryIF(t)
		mockOrderReadRepo := repository.NewOrderReadRepositoryIF(t)
		mockOrderWriteRepo := repository.NewOrderWriteRepositoryIF(t)

		currencyCode, _ := value_objects.NewCurrencyCode("USD")
		totalValue := decimal.NewFromFloat(100.01)
		currencyAmount, _ := value_objects.NewCurrencyAmount(currencyCode, totalValue)

		mockRepo := fixture.OrderEventRepositoryGetFixture(
			t,
			refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID,
			enums.PaymentProcessed,
			currencyAmount,
			currencyAmount,
			nil,
		)

		refundModel := models.NewRefundModel(refundCommand.Amount, refundCommand.PaymentOrderID, refundCommand.ReferenceOrderID)

		mockPartialRefundAdapter.On("PartialRefund", mock.Anything, refundModel, enterpriseID).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.Anything).Return(assert.AnError)

		partialRefundUseCase := NewPartialRefundUse(
			mockRepo,
			mockPartialRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := partialRefundUseCase.PartialRefund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockPartialRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
	})

	t.Run("should return an error getting order entity", func(t *testing.T) {
		mockPartialRefundAdapter := adapters.NewPartialRefundAdapterIF(t)
		mockPaymentOrderReadRepo := repository.NewGetPaymentOrderByReferenceIF(t)
		mockPaymentOrderWriteRepo := repository.NewPaymentOrderRepositoryIF(t)
		mockOrderReadRepo := repository.NewOrderReadRepositoryIF(t)
		mockOrderWriteRepo := repository.NewOrderWriteRepositoryIF(t)

		currencyCode, _ := value_objects.NewCurrencyCode("USD")
		totalValue := decimal.NewFromFloat(100.01)
		currencyAmount, _ := value_objects.NewCurrencyAmount(currencyCode, totalValue)

		mockRepo := fixture.OrderEventRepositoryGetFixture(
			t,
			refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID,
			enums.PaymentProcessed,
			currencyAmount,
			currencyAmount,
			nil,
		)

		refundModel := models.NewRefundModel(refundCommand.Amount, refundCommand.PaymentOrderID, refundCommand.ReferenceOrderID)

		mockPartialRefundAdapter.On("PartialRefund", mock.Anything, refundModel, enterpriseID).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.Anything).Return(nil)

		orderEntity := order_entities.OrderEntity{}
		mockOrderReadRepo.On("GetOrderByReferenceID", ctx, refundCommand.ReferenceOrderID, enterpriseID).
			Return(orderEntity, assert.AnError)

		partialRefundUseCase := NewPartialRefundUse(
			mockRepo,
			mockPartialRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := partialRefundUseCase.PartialRefund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockPartialRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
		mockOrderReadRepo.AssertExpectations(t)
	})

	t.Run("should return an error updating order entity", func(t *testing.T) {
		mockPartialRefundAdapter := adapters.NewPartialRefundAdapterIF(t)
		mockPaymentOrderReadRepo := repository.NewGetPaymentOrderByReferenceIF(t)
		mockPaymentOrderWriteRepo := repository.NewPaymentOrderRepositoryIF(t)
		mockOrderReadRepo := repository.NewOrderReadRepositoryIF(t)
		mockOrderWriteRepo := repository.NewOrderWriteRepositoryIF(t)

		currencyCode, _ := value_objects.NewCurrencyCode("USD")
		totalValue := decimal.NewFromFloat(100.01)
		currencyAmount, _ := value_objects.NewCurrencyAmount(currencyCode, totalValue)

		mockRepo := fixture.OrderEventRepositoryGetFixture(
			t,
			refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID,
			enums.PaymentProcessed,
			currencyAmount,
			currencyAmount,
			nil,
		)

		refundModel := models.NewRefundModel(refundCommand.Amount, refundCommand.PaymentOrderID, refundCommand.ReferenceOrderID)

		mockPartialRefundAdapter.On("PartialRefund", mock.Anything, refundModel, enterpriseID).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.Anything).Return(nil)

		orderEntity := order_entities.OrderEntity{}
		mockOrderReadRepo.On("GetOrderByReferenceID", ctx, refundCommand.ReferenceOrderID, enterpriseID).
			Return(orderEntity, nil)

		mockOrderWriteRepo.On("UpdateOrder", ctx, mock.Anything).Return(assert.AnError)

		partialRefundUseCase := NewPartialRefundUse(
			mockRepo,
			mockPartialRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := partialRefundUseCase.PartialRefund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockPartialRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
		mockOrderReadRepo.AssertExpectations(t)
		mockOrderWriteRepo.AssertExpectations(t)
	})

	t.Run("should success partial refund", func(t *testing.T) {
		mockPartialRefundAdapter := adapters.NewPartialRefundAdapterIF(t)
		mockPaymentOrderReadRepo := repository.NewGetPaymentOrderByReferenceIF(t)
		mockPaymentOrderWriteRepo := repository.NewPaymentOrderRepositoryIF(t)
		mockOrderReadRepo := repository.NewOrderReadRepositoryIF(t)
		mockOrderWriteRepo := repository.NewOrderWriteRepositoryIF(t)

		currencyCode, _ := value_objects.NewCurrencyCode("USD")
		totalValue := decimal.NewFromFloat(100.01)
		currencyAmount, _ := value_objects.NewCurrencyAmount(currencyCode, totalValue)

		mockRepo := fixture.OrderEventRepositoryGetFixture(
			t,
			refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID,
			enums.PaymentProcessed,
			currencyAmount,
			currencyAmount,
			nil,
		)

		refundModel := models.NewRefundModel(refundCommand.Amount, refundCommand.PaymentOrderID, refundCommand.ReferenceOrderID)

		mockPartialRefundAdapter.On("PartialRefund", mock.Anything, refundModel, enterpriseID).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.Anything).Return(nil)

		orderEntity := order_entities.OrderEntity{}
		mockOrderReadRepo.On("GetOrderByReferenceID", ctx, refundCommand.ReferenceOrderID, enterpriseID).
			Return(orderEntity, nil)

		mockOrderWriteRepo.On("UpdateOrder", ctx, mock.Anything).Return(nil)

		mockRepo.On("Save", ctx, mock.Anything).Return(nil)

		partialRefundUseCase := NewPartialRefundUse(
			mockRepo,
			mockPartialRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := partialRefundUseCase.PartialRefund(ctx, refundCommand, enterpriseID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockPartialRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
		mockOrderReadRepo.AssertExpectations(t)
		mockOrderWriteRepo.AssertExpectations(t)
	})
}

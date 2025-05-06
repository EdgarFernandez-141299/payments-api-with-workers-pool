package use_cases

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository/fixture"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/order/dto/response"
	order_entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	payment_order_entities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

var (
	ctx = context.Background()
)

func TestRefundTotalUseCase(t *testing.T) {
	refundCommand := *command.NewRefundTotalCommand(
		"referenceOrderID",
		"paymentOrderID",
		"Customer request",
	)

	enterpriseID := "enterpriseID"

	t.Run("should return an error getting order", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
	})

	t.Run("should return an error, cannot refund payment", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		mockRefundAdapter.On(
			"RefundPayment",
			mock.Anything,
			refundCommand.PaymentOrderID,
			refundCommand.ReferenceOrderID,
			enterpriseID,
			refundCommand.Reason,
		).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.AnythingOfType("entities.PaymentOrderEntity")).Return(nil)

		orderEntity := order_entities.OrderEntity{}
		mockOrderReadRepo.On("GetOrderByReferenceID", ctx, refundCommand.ReferenceOrderID, enterpriseID).
			Return(orderEntity, nil)

		mockOrderWriteRepo.On("UpdateOrder", ctx, mock.AnythingOfType("entities.OrderEntity")).Return(nil)

		mockRepo.On("Save", ctx, mock.Anything).Return(nil)

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
		mockOrderReadRepo.AssertExpectations(t)
		mockOrderWriteRepo.AssertExpectations(t)
	})

	t.Run("should return an error refunding payment", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		mockRefundAdapter.On(
			"RefundPayment",
			mock.Anything,
			refundCommand.PaymentOrderID,
			refundCommand.ReferenceOrderID,
			enterpriseID,
			refundCommand.Reason,
		).
			Return(response.RefundResponseDTO{}, errors.New("refund failed"))

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, "refund failed")
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
	})

	t.Run("should return an error getting payment order", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		mockRefundAdapter.On(
			"RefundPayment",
			mock.Anything,
			refundCommand.PaymentOrderID,
			refundCommand.ReferenceOrderID,
			enterpriseID,
			refundCommand.Reason,
		).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, assert.AnError)

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
	})

	t.Run("should return an error updating payment order", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		mockRefundAdapter.On(
			"RefundPayment",
			mock.Anything,
			refundCommand.PaymentOrderID,
			refundCommand.ReferenceOrderID,
			enterpriseID,
			refundCommand.Reason,
		).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.AnythingOfType("entities.PaymentOrderEntity")).Return(assert.AnError)

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
	})

	t.Run("should return an error getting order entity", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		mockRefundAdapter.On(
			"RefundPayment",
			mock.Anything,
			refundCommand.PaymentOrderID,
			refundCommand.ReferenceOrderID,
			enterpriseID,
			refundCommand.Reason,
		).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.AnythingOfType("entities.PaymentOrderEntity")).Return(nil)

		orderEntity := order_entities.OrderEntity{}
		mockOrderReadRepo.On("GetOrderByReferenceID", ctx, refundCommand.ReferenceOrderID, enterpriseID).
			Return(orderEntity, assert.AnError)

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
		mockOrderReadRepo.AssertExpectations(t)
	})

	t.Run("should return an error updating order entity", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		mockRefundAdapter.On(
			"RefundPayment",
			mock.Anything,
			refundCommand.PaymentOrderID,
			refundCommand.ReferenceOrderID,
			enterpriseID,
			refundCommand.Reason,
		).
			Return(response.RefundResponseDTO{}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.AnythingOfType("entities.PaymentOrderEntity")).Return(nil)

		orderEntity := order_entities.OrderEntity{}
		mockOrderReadRepo.On("GetOrderByReferenceID", ctx, refundCommand.ReferenceOrderID, enterpriseID).
			Return(orderEntity, nil)

		mockOrderWriteRepo.On("UpdateOrder", ctx, mock.AnythingOfType("entities.OrderEntity")).Return(assert.AnError)

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		_, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.EqualError(t, err, assert.AnError.Error())
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
		mockOrderReadRepo.AssertExpectations(t)
		mockOrderWriteRepo.AssertExpectations(t)
	})

	t.Run("should process refund successfully", func(t *testing.T) {
		mockRefundAdapter := adapters.NewRefundAdapterIF(t)
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

		mockRefundAdapter.On(
			"RefundPayment",
			mock.Anything,
			refundCommand.PaymentOrderID,
			refundCommand.ReferenceOrderID,
			enterpriseID,
			refundCommand.Reason,
		).
			Return(response.RefundResponseDTO{
				ReferenceOrderID: refundCommand.ReferenceOrderID,
				PaymentOrderID:   refundCommand.PaymentOrderID,
			}, nil)

		paymentOrder := payment_order_entities.PaymentOrderEntity{}
		mockPaymentOrderReadRepo.On("GetPaymentOrderByReference", ctx, refundCommand.ReferenceOrderID,
			refundCommand.PaymentOrderID, enterpriseID).Return(paymentOrder, nil)

		mockPaymentOrderWriteRepo.On("UpdatePaymentOrder", ctx, mock.AnythingOfType("entities.PaymentOrderEntity")).Return(nil)

		orderEntity := order_entities.OrderEntity{}
		mockOrderReadRepo.On("GetOrderByReferenceID", ctx, refundCommand.ReferenceOrderID, enterpriseID).
			Return(orderEntity, nil)

		mockOrderWriteRepo.On("UpdateOrder", ctx, mock.AnythingOfType("entities.OrderEntity")).Return(nil)

		mockRepo.On("Save", ctx, mock.Anything).Return(nil)

		refundUseCase := NewRefundTotalUseCase(
			mockRepo,
			mockRefundAdapter,
			mockPaymentOrderReadRepo,
			mockPaymentOrderWriteRepo,
			mockOrderReadRepo,
			mockOrderWriteRepo,
		)

		got, err := refundUseCase.Refund(ctx, refundCommand, enterpriseID)

		assert.NoError(t, err)
		assert.Equal(t, "referenceOrderID", got.ReferenceOrderID)
		assert.Equal(t, "paymentOrderID", got.PaymentOrderID)
		mockRepo.AssertExpectations(t)
		mockRefundAdapter.AssertExpectations(t)
		mockPaymentOrderReadRepo.AssertExpectations(t)
		mockPaymentOrderWriteRepo.AssertExpectations(t)
		mockOrderReadRepo.AssertExpectations(t)
		mockOrderWriteRepo.AssertExpectations(t)
	})
}

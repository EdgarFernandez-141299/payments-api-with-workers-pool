package queries

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/logger"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

var (
	ctx          = context.Background()
	orderID      = "order-123"
	enterpriseID = "enterprise-456"
)

func TestQueriesOrderUseCaseImpl_GetOrderStatus(t *testing.T) {
	t.Run("should return order status successfully", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewOrderReadRepositoryIF(t)
		log := logger.NewLogger()

		id := uid.GenerateID()
		metadata := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		metadataJSON, _ := json.Marshal(metadata)

		repositoryMock.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).
			Return(entities.OrderEntity{
				ID:               id.String(),
				ReferenceOrderID: orderID,
				Status:           enums.PaymentProcessing.String(),
				TotalAmount:      decimal.NewFromFloat(100.50),
				CurrencyCode:     "USD",
				CountryCode:      "US",
				Metadata:         json.RawMessage(metadataJSON),
			}, nil)

		useCase := NewQueriesOrderUseCase(log, repositoryMock)

		// Act
		result, err := useCase.GetOrderDetail(ctx, orderID, enterpriseID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, orderID, result.ReferenceOrderID)
		assert.Equal(t, enums.PaymentProcessing, result.Status)
		assert.Equal(t, decimal.NewFromFloat(100.50), result.Total)
		assert.Equal(t, "USD", result.Currency)
		assert.Equal(t, "US", result.CountryCode)
		assert.Equal(t, metadata, result.Metadata)

		repositoryMock.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewOrderReadRepositoryIF(t)
		log := logger.NewLogger()

		repositoryMock.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).
			Return(entities.OrderEntity{}, errors.New("database error"))

		useCase := NewQueriesOrderUseCase(log, repositoryMock)

		// Act
		result, err := useCase.GetOrderDetail(ctx, orderID, enterpriseID)

		// Assert
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Empty(t, result.ReferenceOrderID)

		repositoryMock.AssertExpectations(t)
	})

	t.Run("should return error when metadata is invalid JSON", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewOrderReadRepositoryIF(t)
		log := logger.NewLogger()

		id := uid.GenerateID()
		repositoryMock.On("GetOrderByReferenceID", mock.Anything, orderID, enterpriseID).
			Return(entities.OrderEntity{
				ID:               id.String(),
				ReferenceOrderID: orderID,
				Status:           enums.PaymentProcessing.String(),
				TotalAmount:      decimal.NewFromFloat(100.50),
				CurrencyCode:     "USD",
				CountryCode:      "US",
				Metadata:         json.RawMessage("123"),
			}, nil)

		useCase := NewQueriesOrderUseCase(log, repositoryMock)

		// Act
		result, err := useCase.GetOrderDetail(ctx, orderID, enterpriseID)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result.ReferenceOrderID)

		repositoryMock.AssertExpectations(t)
	})
}

func TestGetOrderPayment(t *testing.T) {
	t.Run("should return order with payments successfully", func(t *testing.T) {
		repositoryMock := mockRepository.NewOrderReadRepositoryIF(t)
		log := logger.NewLogger()

		id := uid.GenerateID()
		metadata := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		metadataJSON, _ := json.Marshal(metadata)

		repositoryMock.On("GetOrderPayments", mock.Anything, orderID, enterpriseID).
			Return([]projections.OrderPaymentsProjection{
				{
					ReferenceOrderID:  orderID,
					TotalAmount:       decimal.NewFromFloat(100.50),
					CurrencyCode:      "USD",
					CountryCode:       "US",
					Metadata:          string(metadataJSON),
					PaymentID:         id.String(),
					PaymentStatus:     enums.PaymentProcessing.String(),
					PaymentMethod:     "card",
					CardID:            "card-123",
					AuthorizationCode: "123456",
					PaymentOrderID:    "payment-123",
				},
			}, nil)

		useCase := NewQueriesOrderUseCase(log, repositoryMock)

		result, err := useCase.GetOrderPayments(ctx, orderID, enterpriseID)

		assert.NoError(t, err)
		assert.Equal(t, orderID, result.ReferenceOrderID)
		assert.Equal(t, "PROCESSING", result.Payments[0].Status)
		assert.Equal(t, decimal.NewFromFloat(100.50), result.Total)
		assert.Equal(t, "USD", result.Currency)
		assert.Equal(t, "US", result.CountryCode)
		assert.Equal(t, metadata, result.Metadata)

		repositoryMock.AssertExpectations(t)
		repositoryMock.ExpectedCalls = nil
	})

	t.Run("should return error when order is not found", func(t *testing.T) {
		repositoryMock := mockRepository.NewOrderReadRepositoryIF(t)
		log := logger.NewLogger()

		repositoryMock.On("GetOrderPayments", mock.Anything, orderID, enterpriseID).
			Return([]projections.OrderPaymentsProjection{}, nil)

		useCase := NewQueriesOrderUseCase(log, repositoryMock)

		result, err := useCase.GetOrderPayments(ctx, orderID, enterpriseID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		repositoryMock := mockRepository.NewOrderReadRepositoryIF(t)
		log := logger.NewLogger()

		repositoryMock.On("GetOrderPayments", mock.Anything, orderID, enterpriseID).
			Return([]projections.OrderPaymentsProjection{
				{
					ReferenceOrderID:  orderID,
					TotalAmount:       decimal.NewFromFloat(100.50),
					CurrencyCode:      "USD",
					CountryCode:       "US",
					Metadata:          "invalid json",
					PaymentID:         "id",
					PaymentStatus:     enums.PaymentProcessing.String(),
					PaymentMethod:     "card",
					CardID:            "card-123",
					AuthorizationCode: "123456",
					PaymentOrderID:    "payment-123",
				},
			}, nil)

		useCase := NewQueriesOrderUseCase(log, repositoryMock)

		result, err := useCase.GetOrderPayments(ctx, orderID, enterpriseID)

		assert.Error(t, err)
		assert.Nil(t, result)

		repositoryMock.AssertExpectations(t)
		repositoryMock.ExpectedCalls = nil
	})

	t.Run("should return error when metadata is invalid JSON", func(t *testing.T) {
		repositoryMock := mockRepository.NewOrderReadRepositoryIF(t)
		log := logger.NewLogger()

		repositoryMock.On("GetOrderPayments", mock.Anything, orderID, enterpriseID).
			Return(nil, errors.New("invalid character 'i' looking for beginning of value"))

		useCase := NewQueriesOrderUseCase(log, repositoryMock)

		result, err := useCase.GetOrderPayments(ctx, orderID, enterpriseID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

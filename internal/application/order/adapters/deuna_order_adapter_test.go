package adapters

import (
	"context"
	"os"
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository/fixture"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	entitiesOrder "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

var (
	ctx          = context.TODO()
	enterpriseID = "enterprise-1234"
	referenceID  = "1234"
	id, _        = uid.NewUniqueID(uid.WithID("id"))
)

func TestCreateOrder(t *testing.T) {
	orderID := "123"
	paymentID := "456"
	deunaOrderID := utils.NewDeunaOrderID(orderID, paymentID)
	totalValue := decimal.NewFromFloat(1.0)
	orderToken := "2345"
	currencyCode := "USD"
	paymentFlow := enums.Autocapture

	orderDTO := request.CreateDeunaOrderRequestDTO{
		Order: request.DeunaOrder{
			OrderID:     deunaOrderID.GetID(),
			Currency:    currencyCode,
			TotalAmount: utils.NewDeunaAmount(totalValue),
			SubTotal:    utils.NewDeunaAmount(totalValue),
			StoreCode:   "all",
			WebhooksURL: request.WebhooksURL{
				NotifyOrder: os.Getenv("DEUNA_NOTIFY_ORDER"),
			},
			Metadata: &map[string]interface{}{
				"additional_data.flow_type": "auto-capture",
			},
		},
		OrderType: request.DeunaOrderType(request.DeUnaNow.String()),
	}

	t.Run("should create the order and return no error", func(t *testing.T) {
		readRepositoryMock := repository.NewOrderReadRepositoryIF(t)
		resourceMock := resources.NewDeunaOrderResourceIF(t)

		resourceMock.On("CreateOrder", mock.Anything, orderDTO).Return(response.DeunaOrderResponseDTO{
			Token: orderToken,
		}, nil)

		deunaRepository := fixture.CreateDeunaPaymentOrderRepositoryFixture(t, orderID, paymentID, orderToken, nil)

		deunaRepository.On("GetTokenByOrderAndPaymentID", mock.Anything, orderID, paymentID).Return("", nil)

		adapter := NewOrderAdapter(readRepositoryMock, deunaRepository, resourceMock)

		got, err := adapter.CreateOrder(ctx, orderID, paymentID, currencyCode, totalValue, paymentFlow)
		assert.NoError(t, err)
		assert.Equal(t, orderToken, got.Token)

		deunaRepository.AssertExpectations(t)
		readRepositoryMock.AssertExpectations(t)
		resourceMock.AssertExpectations(t)
	})

	t.Run("should return token when order exist in repository", func(t *testing.T) {
		readRepositoryMock := repository.NewOrderReadRepositoryIF(t)
		resourceMock := resources.NewDeunaOrderResourceIF(t)
		deunaRepository := repository.NewDeunaOrderRepository(t)

		adapter := NewOrderAdapter(readRepositoryMock, deunaRepository, resourceMock)

		deunaRepository.On("GetTokenByOrderAndPaymentID", mock.Anything, orderID, paymentID).Return("nop", nil)

		got, err := adapter.CreateOrder(ctx, orderID, paymentID, currencyCode, totalValue, paymentFlow)
		assert.NoError(t, err)
		assert.Equal(t, "nop", got.Token)

		deunaRepository.AssertExpectations(t)
		readRepositoryMock.AssertExpectations(t)
		resourceMock.AssertExpectations(t)
	})

	t.Run("should return error when payment flow is invalid", func(t *testing.T) {
		readRepositoryMock := repository.NewOrderReadRepositoryIF(t)
		resourceMock := resources.NewDeunaOrderResourceIF(t)
		deunaRepository := repository.NewDeunaOrderRepository(t)

		adapter := NewOrderAdapter(readRepositoryMock, deunaRepository, resourceMock)

		deunaRepository.On("GetTokenByOrderAndPaymentID", mock.Anything, orderID, paymentID).Return("", nil)

		got, err := adapter.CreateOrder(ctx, orderID, paymentID, currencyCode, totalValue, enums.PaymentFlowEnum("INVALID"))
		assert.Error(t, err)
		assert.Equal(t, enums.ErrInvalidPaymentFlow, err)
		assert.Empty(t, got.Token)

		deunaRepository.AssertExpectations(t)
		readRepositoryMock.AssertExpectations(t)
		resourceMock.AssertExpectations(t)
	})
}

func TestGetOrderByReferenceID(t *testing.T) {
	t.Run("should get an error getting order by reference id", func(t *testing.T) {
		readRepositoryMock := repository.NewOrderReadRepositoryIF(t)
		resourceMock := resources.NewDeunaOrderResourceIF(t)
		deunaOrderRepository := repository.NewDeunaOrderRepository(t)

		readRepositoryMock.On("GetOrderByReferenceID", ctx, referenceID, enterpriseID).
			Return(entitiesOrder.OrderEntity{}, assert.AnError)

		adapter := NewOrderAdapter(readRepositoryMock, deunaOrderRepository, resourceMock)
		got, err := adapter.GetOrderByReferenceID(ctx, referenceID, enterpriseID)

		assert.Equal(t, assert.AnError, err)
		assert.Zero(t, got)

		readRepositoryMock.AssertExpectations(t)
		resourceMock.AssertExpectations(t)
		deunaOrderRepository.AssertExpectations(t)
	})

	t.Run("should get order successful", func(t *testing.T) {
		readRepositoryMock := repository.NewOrderReadRepositoryIF(t)
		resourceMock := resources.NewDeunaOrderResourceIF(t)
		deunaRepository := repository.NewDeunaOrderRepository(t)

		readRepositoryMock.On("GetOrderByReferenceID", ctx, referenceID, enterpriseID).
			Return(entitiesOrder.OrderEntity{
				ID: id.String(),
			}, nil)

		adapter := NewOrderAdapter(readRepositoryMock, deunaRepository, resourceMock)
		got, err := adapter.GetOrderByReferenceID(ctx, referenceID, enterpriseID)

		assert.Nil(t, err)
		assert.NotZero(t, got)

		deunaRepository.AssertExpectations(t)
		readRepositoryMock.AssertExpectations(t)
		resourceMock.AssertExpectations(t)
	})
}

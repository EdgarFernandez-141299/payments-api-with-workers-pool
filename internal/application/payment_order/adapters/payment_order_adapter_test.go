package adapter

import (
	"context"
	"fmt"
	"os"
	"testing"

	fixture2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters/fixture"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository/fixture"

	resources "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
	mockResource "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources/fixture"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/aggregate"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

func TestPaymentOrderAdapter_CreatePaymentOrder(t *testing.T) {
	ctx := context.Background()
	enterpriseID := "test-enterprise"
	orderID := "test-order"
	userID := "test-user"
	cardID := "test-card"
	externalCardID := "test-external-card"
	email := "test@example.com"
	cvv := "123"
	countryCode := "MX"

	currencyCode := value_objects.CurrencyCode{
		Code:   "USD",
		Symbol: "$",
	}

	cmd := command.CreatePaymentOrderCommand{
		ID:               "test-payment",
		ReferenceOrderID: orderID,
		User: entities.User{
			ID: userID,
		},
		CurrencyCode: currencyCode,
		CountryCode:  countryCode,
		Payment: entities.PaymentOrder{
			ID: "test-payment",
			Total: value_objects.CurrencyAmount{
				Value: decimal.NewFromFloat(100.0),
				Code:  currencyCode,
			},
			Method: value_objects.PaymentMethod{
				Type: enums.CCMethod,
				CCData: value_objects.PaymentMethodData[value_objects.CardInfo]{
					Data: value_objects.CardInfo{
						CardID: cardID,
						CVV:    cvv,
					},
				},
			},
			Status: enums.PaymentProcessing,
			OriginType: value_objects.AssociatedOrigin{
				Type: enums.Club,
			},
			CollectionAccount: entities.CollectionAccount{
				ID: "test-collection",
			},
		},
		AssociatedOrigin: value_objects.AssociatedOrigin{
			Type: enums.Club,
		},
		CollectionAccount: entities.CollectionAccount{
			ID: "test-collection",
		},
	}

	order := &aggregate.Order{
		ID:           orderID,
		Email:        email,
		EnterpriseID: enterpriseID,
		AllowCapture: true,
	}

	card := entities.Card{
		ID:             cardID,
		ExternalCardID: externalCardID,
		CardType:       enums.CreditCard,
	}

	os.Setenv("DEUNA_NOTIFY_ORDER", "http://test-webhook.com")
	defer os.Unsetenv("DEUNA_NOTIFY_ORDER")

	t.Run("success", func(t *testing.T) {
		deunaLoginMock, userToken := fixture2.DeunaLoginFixture(t, userID, enterpriseID, nil)

		expectedPaymentFlow, _ := enums.NewPaymentFlowEnum(card.CardType, order.AllowCapture)
		cmd.PaymentFlow = expectedPaymentFlow

		orderAdapterMock, deunaCreateOrderResponse := fixture2.NewCreateOrderFixture(
			t, orderID, cmd.Payment.ID, cmd.CurrencyCode.Code, cmd.Payment.Total.Value, expectedPaymentFlow, nil,
		)

		orderToken := deunaCreateOrderResponse.Token

		resourceMock := mockResource.MakeOrderPaymentFixture(t, userToken, orderToken, nil)
		repositoryMock := fixture.CreatePaymentOrderReadModelRepository(t, orderID, cardID, cmd, nil)
		adapter := NewPaymentOrderAdapter(resourceMock, orderAdapterMock, deunaLoginMock, repositoryMock)
		err := adapter.CreatePaymentOrder(ctx, cmd, order, card)

		assert.NoError(t, err)

		resourceMock.AssertExpectations(t)
		resourceMock.AssertExpectations(t)
		orderAdapterMock.AssertExpectations(t)
		deunaLoginMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("failure_due_to_invalid_card", func(t *testing.T) {
		deunaLoginMock, userToken := fixture2.DeunaLoginFixture(t, userID, enterpriseID, nil)

		expectedPaymentFlow, _ := enums.NewPaymentFlowEnum(card.CardType, order.AllowCapture)
		cmd.PaymentFlow = expectedPaymentFlow

		orderAdapterMock, orderResponse := fixture2.NewCreateOrderFixture(
			t, orderID, cmd.Payment.ID, cmd.CurrencyCode.Code, cmd.Payment.Total.Value, expectedPaymentFlow, nil,
		)

		resourceMock := mockResource.MakeOrderPaymentFixture(t, userToken, orderResponse.Token, fmt.Errorf("invalid card information"))
		repositoryMock := mockRepository.NewPaymentOrderRepositoryIF(t)
		adapter := NewPaymentOrderAdapter(resourceMock, orderAdapterMock, deunaLoginMock, repositoryMock)
		err := adapter.CreatePaymentOrder(ctx, cmd, order, card)

		assert.Error(t, err)
		assert.EqualError(t, err, "invalid card information")

		resourceMock.AssertExpectations(t)
		orderAdapterMock.AssertExpectations(t)
		deunaLoginMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("failure_due_to_missing_country_code", func(t *testing.T) {
		cmdWithoutCountry := cmd
		cmdWithoutCountry.CountryCode = ""

		deunaLoginMock, _ := fixture2.DeunaLoginFixture(t, userID, enterpriseID, nil)
		orderAdapterMock, _ := fixture2.NewCreateOrderFixture(t, orderID, cmd.Payment.ID, cmd.CurrencyCode.Code, cmd.Payment.Total.Value, cmd.PaymentFlow, nil)
		resourceMock := resources.NewDeunaPaymentResourceIF(t)
		resourceMock.On("MakeOrderPayment", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("country code is required"))
		repositoryMock := mockRepository.NewPaymentOrderRepositoryIF(t)

		adapter := NewPaymentOrderAdapter(resourceMock, orderAdapterMock, deunaLoginMock, repositoryMock)
		err := adapter.CreatePaymentOrder(ctx, cmdWithoutCountry, order, card)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "country code is required")

		resourceMock.AssertExpectations(t)
		orderAdapterMock.AssertExpectations(t)
		deunaLoginMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})
}

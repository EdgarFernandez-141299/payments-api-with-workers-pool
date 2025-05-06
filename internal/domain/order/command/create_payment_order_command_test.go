package command

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestWithWithCollectionAccount(t *testing.T) {
	t.Run("should set collection account", func(t *testing.T) {
		oldCmd := NewCreatePaymentOrderCommandBuilder().
			WithReferenceOrderID("32").
			Build()

		got := oldCmd.WithCollectionAccount(entities.CollectionAccount{ID: "1"})

		assert.Equal(t, "1", got.CollectionAccount.ID)
		assert.Equal(t, "32", got.ReferenceOrderID)
		assert.Equal(t, "", oldCmd.CollectionAccount.ID)
	})
}

func TestCreatePaymentOrderCommandBuilder(t *testing.T) {
	associatedOrigin, _ := value_objects.NewFromAssociatedOriginString(enums.Downpayment.String())
	currency := value_objects.CurrencyCode{Code: "USD"}
	user := entities.NewUser(value_objects.Member, "123")
	method := value_objects.PaymentMethod{
		Type: enums.CCMethod,
		CCData: value_objects.PaymentMethodData[value_objects.CardInfo]{
			Data: value_objects.CardInfo{
				CardID: "123",
				CVV:    "123",
			},
		},
	}
	totalAmount, _ := value_objects.NewCurrencyAmount(currency, decimal.NewFromFloat(100.0))

	t.Run("should create a new payment order command", func(t *testing.T) {
		builder := NewCreatePaymentOrderCommandBuilder()

		builder.WithReferenceOrderID("123")
		builder.WithUser(user)
		builder.WithPayment(entities.PaymentOrder{
			ID:         "123",
			OriginType: associatedOrigin,
			Total:      totalAmount,
			Status:     enums.PaymentProcessing,
			Method:     method,
		})
		builder.WithAssociatedOrigin(associatedOrigin)
		builder.WithCurrencyCode(currency)

		builder.WithPaymentOrderID("123")

		command := builder.Build()

		assert.Equal(t, "123", command.ReferenceOrderID)
		assert.Equal(t, "123", command.ID)
		assert.Nil(t, command.Validate())
	})

	t.Run("should return error ReferenceOrderID is required", func(t *testing.T) {
		builder := NewCreatePaymentOrderCommandBuilder()

		builder.WithUser(user)
		builder.WithReferenceOrderID("")
		builder.WithPayment(entities.PaymentOrder{
			ID:         "123",
			OriginType: associatedOrigin,
			Total:      totalAmount,
			Status:     enums.PaymentProcessing,
			Method:     method,
		})

		builder.WithPaymentOrderID("123")

		command := builder.Build()

		assert.Equal(t, "", command.ReferenceOrderID)
		assert.Equal(t, "123", command.ID)
		assert.EqualError(t, command.Validate(), "ReferenceOrderID is required")
	})
	t.Run("should return error payment invalid", func(t *testing.T) {
		builder := NewCreatePaymentOrderCommandBuilder()

		builder.WithUser(user)
		builder.WithReferenceOrderID("123")
		builder.WithPayment(entities.PaymentOrder{})

		command := builder.Build()

		assert.EqualError(t, command.Validate(), "type is not valid")
	})

	t.Run("should return error user invalid", func(t *testing.T) {
		builder := NewCreatePaymentOrderCommandBuilder()

		builder.WithUser(entities.User{
			ID:   "",
			Type: value_objects.Member,
		})
		builder.WithReferenceOrderID("123")
		builder.WithPayment(entities.PaymentOrder{
			ID:         "123",
			OriginType: associatedOrigin,
			Total:      totalAmount,
			Status:     enums.PaymentProcessing,
			Method:     method,
		})

		command := builder.Build()

		assert.EqualError(t, command.Validate(), "Business Error code: INVALID_USER_ID, message: invalid user id: ")
	})
}

func TestSetAuthorizationCode(t *testing.T) {
	t.Run("should set authorization code", func(t *testing.T) {
		paymentOrderCommand := new(CreatePaymentOrderCommand)

		paymentOrder := paymentOrderCommand.WithAuthorizationCode("123")

		assert.Equal(t, "123", paymentOrder.AuthorizationCode)
	})
}

func TestWithPaymentFlow(t *testing.T) {
	t.Run("should set payment flow", func(t *testing.T) {
		oldCmd := NewCreatePaymentOrderCommandBuilder().
			WithReferenceOrderID("32").
			Build()

		got := oldCmd.WithPaymentFlow(enums.PaymentFlowEnum("DOWNPAYMENT"))

		assert.Equal(t, enums.PaymentFlowEnum("DOWNPAYMENT"), got.PaymentFlow)
		assert.Equal(t, "DOWNPAYMENT", got.Payment.PaymentFlow)
		assert.Equal(t, "32", got.ReferenceOrderID)
		assert.Equal(t, "", oldCmd.PaymentFlow.String())
	})
}

func TestWithCardData(t *testing.T) {
	t.Run("should set card data", func(t *testing.T) {
		oldCmd := NewCreatePaymentOrderCommandBuilder().
			WithReferenceOrderID("32").
			Build()

		cardData := entities.Card{
			Brand:    "VISA",
			LastFour: "1234",
			CardType: "CREDIT",
		}

		got := oldCmd.WithCardData(cardData)

		assert.Equal(t, "VISA", got.Payment.PaymentCard.CardBrand)
		assert.Equal(t, "1234", got.Payment.PaymentCard.CardLast4)
		assert.Equal(t, "CREDIT", got.Payment.PaymentCard.CardType)
		assert.Equal(t, "32", got.ReferenceOrderID)
		assert.Equal(t, "", oldCmd.Payment.PaymentCard.CardBrand)
	})
}

func TestCreatePaymentOrderCommandBuilder_WithCountryCode(t *testing.T) {
	t.Run("should set country code", func(t *testing.T) {
		builder := NewCreatePaymentOrderCommandBuilder()
		builder.WithCountryCode("MX")
		command := builder.Build()

		assert.Equal(t, "MX", command.CountryCode)
	})
}

func TestCreatePaymentOrderCommand_Validate(t *testing.T) {
	associatedOrigin, _ := value_objects.NewFromAssociatedOriginString(enums.Downpayment.String())
	currency := value_objects.CurrencyCode{Code: "USD"}
	user := entities.NewUser(value_objects.Member, "123")
	method := value_objects.PaymentMethod{
		Type: enums.CCMethod,
		CCData: value_objects.PaymentMethodData[value_objects.CardInfo]{
			Data: value_objects.CardInfo{
				CardID: "123",
				CVV:    "123",
			},
		},
	}
	totalAmount, _ := value_objects.NewCurrencyAmount(currency, decimal.NewFromFloat(100.0))

	tests := []struct {
		name    string
		command CreatePaymentOrderCommand
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command",
			command: CreatePaymentOrderCommand{
				ReferenceOrderID: "123",
				Payment: entities.PaymentOrder{
					ID:         "123",
					OriginType: associatedOrigin,
					Total:      totalAmount,
					Status:     enums.PaymentProcessing,
					Method:     method,
				},
				User: user,
			},
			wantErr: false,
		},
		{
			name: "missing reference order id",
			command: CreatePaymentOrderCommand{
				ReferenceOrderID: "",
				Payment: entities.PaymentOrder{
					ID:         "123",
					OriginType: associatedOrigin,
					Total:      totalAmount,
					Status:     enums.PaymentProcessing,
					Method:     method,
				},
				User: user,
			},
			wantErr: true,
			errMsg:  "ReferenceOrderID is required",
		},
		{
			name: "invalid payment",
			command: CreatePaymentOrderCommand{
				ReferenceOrderID: "123",
				Payment:          entities.PaymentOrder{},
				User:             user,
			},
			wantErr: true,
			errMsg:  "type is not valid",
		},
		{
			name: "invalid user",
			command: CreatePaymentOrderCommand{
				ReferenceOrderID: "123",
				Payment: entities.PaymentOrder{
					ID:         "123",
					OriginType: associatedOrigin,
					Total:      totalAmount,
					Status:     enums.PaymentProcessing,
					Method:     method,
				},
				User: entities.User{
					ID:   "",
					Type: value_objects.Member,
				},
			},
			wantErr: true,
			errMsg:  "Business Error code: INVALID_USER_ID, message: invalid user id: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.command.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

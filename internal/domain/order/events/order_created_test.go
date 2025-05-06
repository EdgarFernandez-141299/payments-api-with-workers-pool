package events

import (
	"testing"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

func TestFromCreateOrderCommand(t *testing.T) {
	userID := "user123"
	user := entities.NewUser(vo.NewUserType(vo.Member), userID)
	usdCurrencyCode, _ := vo.NewCurrencyCode("USD")
	totalAmountAsFloat := decimal.NewFromFloat(100.5)
	countryCode, _ := vo.NewCountryWithCode("MX")

	totalAmount, _ := vo.NewCurrencyAmount(usdCurrencyCode, totalAmountAsFloat)
	invalidAmount, _ := vo.NewCurrencyAmount(usdCurrencyCode, decimal.NewFromFloat(-10))
	invalidUser := entities.NewUser(vo.NewUserType(vo.Member), "")

	tests := []struct {
		name      string
		command   command.CreateOrderCommand
		wantErr   bool
		errorType error
	}{
		{
			name: "valid command",
			command: command.CreateOrderCommand{
				ReferenceID:  "ref123",
				TotalAmount:  totalAmount,
				PhoneNumber:  "123456789",
				User:         user,
				CurrencyCode: usdCurrencyCode,
				CountryCode:  countryCode,
			},
			wantErr: false,
		},
		{
			name: "missing reference ID",
			command: command.CreateOrderCommand{
				ReferenceID:  "",
				TotalAmount:  totalAmount,
				PhoneNumber:  "123456789",
				User:         user,
				CurrencyCode: usdCurrencyCode,
				CountryCode:  countryCode,
			},
			wantErr:   true,
			errorType: errors.ErrInvalidOrderID,
		},
		{
			name: "invalid total amount",
			command: command.CreateOrderCommand{
				ReferenceID:  "ref123",
				TotalAmount:  invalidAmount,
				PhoneNumber:  "123456789",
				User:         user,
				CurrencyCode: usdCurrencyCode,
				CountryCode:  countryCode,
			},
			wantErr:   true,
			errorType: errors.ErrInvalidOrderTotalAmount,
		},
		{
			name: "missing phone number",
			command: command.CreateOrderCommand{
				ReferenceID:  "ref123",
				TotalAmount:  totalAmount,
				PhoneNumber:  "",
				User:         user,
				CurrencyCode: usdCurrencyCode,
				CountryCode:  countryCode,
			},
			wantErr:   true,
			errorType: errors.ErrInvalidOrderPhoneNumber,
		},
		{
			name: "missing user type",
			command: command.CreateOrderCommand{
				ReferenceID:  "ref123",
				TotalAmount:  totalAmount,
				PhoneNumber:  "123456789",
				User:         invalidUser,
				CurrencyCode: usdCurrencyCode,
				CountryCode:  countryCode,
			},
			wantErr:   true,
			errorType: errors.ErrInvalidOrderUserType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderEv := NewOrderCreatedEventBuilder().
				SetUser(tt.command.User).
				SetID(tt.command.ReferenceID).
				SetTotalAmount(tt.command.TotalAmount).
				SetPhoneNumber(tt.command.PhoneNumber).
				SetCountryCode(tt.command.CountryCode).
				SetBillingAddress(tt.command.BillingAddress).
				SetEnterpriseID(tt.command.EnterpriseID).
				SetShippingAddress(tt.command.BillingAddress).
				SetEmail(tt.command.Email).
				SetCreatedAt(time.Now().UTC()).
				Build()

			assert.Equal(t, tt.command.User, orderEv.User)
			assert.Equal(t, tt.command.ReferenceID, orderEv.ID)
			assert.Equal(t, tt.command.TotalAmount, orderEv.TotalAmount)
			assert.Equal(t, tt.command.PhoneNumber, orderEv.PhoneNumber)
			assert.Equal(t, tt.command.CountryCode, orderEv.CountryCode)
			assert.Equal(t, tt.command.BillingAddress, orderEv.BillingAddress)
			assert.Equal(t, tt.command.EnterpriseID, orderEv.EnterpriseID)
			assert.Equal(t, tt.command.BillingAddress, orderEv.BillingAddress)
			assert.Equal(t, tt.command.Email, orderEv.Email)
		})
	}
}

func TestOrderCreatedEventBuilder(t *testing.T) {
	userID := "user123"
	user := entities.NewUser(vo.NewUserType(vo.Member), userID)
	usdCurrencyCode, _ := vo.NewCurrencyCode("USD")
	totalAmount, _ := vo.NewCurrencyAmount(usdCurrencyCode, decimal.NewFromFloat(100.5))
	countryCode, _ := vo.NewCountryWithCode("MX")
	address := vo.NewAddress("12345", "Test Street", "Test City", countryCode)
	webhookUrl := vo.NewWebhookUrl("https://test.com/webhook")
	createdAt := time.Now().UTC()
	metadata := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	event := NewOrderCreatedEventBuilder().
		SetID("order123").
		SetTotalAmount(totalAmount).
		SetPhoneNumber("123456789").
		SetUser(user).
		SetCreatedAt(createdAt).
		SetCountryCode(countryCode).
		SetBillingAddress(address).
		SetShippingAddress(address).
		SetEnterpriseID("enterprise123").
		SetEmail("test@example.com").
		SetWebhookUrl(webhookUrl).
		SetMetadata(metadata).
		SetAllowCapture(true).
		Build()

	assert.Equal(t, "order123", event.ID)
	assert.Equal(t, totalAmount, event.TotalAmount)
	assert.Equal(t, "123456789", event.PhoneNumber)
	assert.Equal(t, user, event.User)
	assert.Equal(t, createdAt, event.CreatedAt)
	assert.Equal(t, countryCode, event.CountryCode)
	assert.Equal(t, address, event.BillingAddress)
	assert.Equal(t, address, event.ShippingAddress)
	assert.Equal(t, "enterprise123", event.EnterpriseID)
	assert.Equal(t, "test@example.com", event.Email)
	assert.Equal(t, webhookUrl, event.WebhookUrl)
	assert.Equal(t, metadata, event.Metadata)
	assert.True(t, event.AllowCapture)
}

func TestOrderCreatedEventBuilder_EmptyValues(t *testing.T) {
	event := NewOrderCreatedEventBuilder().Build()

	assert.Empty(t, event.ID)
	assert.Empty(t, event.TotalAmount)
	assert.Empty(t, event.PhoneNumber)
	assert.Empty(t, event.User)
	assert.Empty(t, event.CreatedAt)
	assert.Empty(t, event.CountryCode)
	assert.Empty(t, event.BillingAddress)
	assert.Empty(t, event.ShippingAddress)
	assert.Empty(t, event.EnterpriseID)
	assert.Empty(t, event.Email)
	assert.Empty(t, event.WebhookUrl)
	assert.Nil(t, event.Metadata)
	assert.False(t, event.AllowCapture)
}

func TestOrderCreatedEventBuilder_PartialValues(t *testing.T) {
	userID := "user123"
	user := entities.NewUser(vo.NewUserType(vo.Member), userID)
	metadata := map[string]interface{}{
		"test": "value",
	}
	webhookUrl := vo.WebhookUrl{Url: "https://webhook.example.com"}

	event := NewOrderCreatedEventBuilder().
		SetID("order123").
		SetUser(user).
		SetMetadata(metadata).
		SetAllowCapture(true).
		SetWebhookUrl(webhookUrl).
		Build()

	assert.Equal(t, "order123", event.ID)
	assert.Equal(t, user, event.User)
	assert.Equal(t, metadata, event.Metadata)
	assert.True(t, event.AllowCapture)
	assert.Equal(t, webhookUrl, event.WebhookUrl)
	assert.Empty(t, event.TotalAmount)
	assert.Empty(t, event.PhoneNumber)
	assert.Empty(t, event.CreatedAt)
	assert.Empty(t, event.CountryCode)
	assert.Empty(t, event.BillingAddress)
	assert.Empty(t, event.ShippingAddress)
	assert.Empty(t, event.EnterpriseID)
	assert.Empty(t, event.Email)
}

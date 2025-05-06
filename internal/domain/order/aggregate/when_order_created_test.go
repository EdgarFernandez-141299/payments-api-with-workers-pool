package aggregate

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestWhenOrderCreated(t *testing.T) {
	t.Parallel()

	userID := "user123"
	user := entities.NewUser(vo.NewUserType(vo.Member), userID)
	usdCurrencyCode, _ := vo.NewCurrencyCode("USD")
	totalAmountAsFloat := decimal.NewFromFloat(100.5)
	enterpriseID := "enterprise123"
	countryCode, _ := vo.NewCountryWithCode("MX")
	metadata := map[string]interface{}{
		"bill": "123",
	}
	webhookUrl := vo.NewWebhookUrl("https://webhook.example.com")
	allowCapture := true

	totalAmount, _ := vo.NewCurrencyAmount(usdCurrencyCode, totalAmountAsFloat)

	ev := events.NewOrderCreatedEventBuilder().
		SetID("1234").
		SetTotalAmount(totalAmount).
		SetPhoneNumber("3222332").
		SetUser(user).
		SetCountryCode(countryCode).
		SetEnterpriseID(enterpriseID).
		SetMetadata(metadata).
		SetWebhookUrl(webhookUrl).
		SetAllowCapture(allowCapture).
		Build()

	tests := []struct {
		name     string
		order    *Order
		event    events.OrderCreated
		expected *Order
	}{
		{
			name:  "valid order creation",
			order: &Order{},
			event: *ev,
			expected: &Order{
				ID:            "1234",
				TotalAmount:   ev.TotalAmount,
				PhoneNumber:   "3222332",
				User:          ev.User,
				CountryCode:   countryCode,
				Status:        vo.OrderStatusProcessing(),
				EnterpriseID:  enterpriseID,
				Currency:      usdCurrencyCode,
				Metadata:      ev.Metadata,
				OrderPayments: []entities.PaymentOrder{},
				WebhookUrl:    webhookUrl,
				AllowCapture:  allowCapture,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			WhenOrderCreated(tt.order, tt.event)
			assert.Equal(t, tt.expected, tt.order, "WhenOrderCreated() = %v, want %v", tt.order, tt.expected)
		})
	}
}

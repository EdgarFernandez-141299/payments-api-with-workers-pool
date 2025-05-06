package command

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewPaymentOrderCapturedCommand(t *testing.T) {
	tests := []struct {
		name      string
		orderID   string
		paymentID string
		currency  string
		amount    decimal.Decimal
		expected  PaymentOrderCapturedCommand
	}{
		{
			name:      "crear comando de captura de pago exitoso",
			orderID:   "order-123",
			paymentID: "payment-456",
			currency:  "USD",
			amount:    decimal.NewFromFloat(100.50),
			expected: PaymentOrderCapturedCommand{
				OrderID:   "order-123",
				PaymentID: "payment-456",
				Currency:  "USD",
				Amount:    decimal.NewFromFloat(100.50),
			},
		},
		{
			name:      "crear comando de captura de pago con monto cero",
			orderID:   "order-123",
			paymentID: "payment-456",
			currency:  "MXN",
			amount:    decimal.Zero,
			expected: PaymentOrderCapturedCommand{
				OrderID:   "order-123",
				PaymentID: "payment-456",
				Currency:  "MXN",
				Amount:    decimal.Zero,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewPaymentOrderCapturedCommand(tt.orderID, tt.paymentID, tt.currency, tt.amount)

			assert.Equal(t, tt.expected.OrderID, cmd.OrderID)
			assert.Equal(t, tt.expected.PaymentID, cmd.PaymentID)
			assert.Equal(t, tt.expected.Currency, cmd.Currency)
			assert.True(t, tt.expected.Amount.Equal(cmd.Amount))
		})
	}
}

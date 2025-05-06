package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCreatePaymentOrderProcessedCommand(t *testing.T) {
	tests := []struct {
		name      string
		orderID   string
		paymentID string
		authCode  string
		status    string
		expected  CreatePaymentOrderProcessedCommand
	}{
		{
			name:      "create processed payment command",
			orderID:   "order-1",
			paymentID: "payment-1",
			authCode:  "auth-123",
			status:    "completed",
			expected: CreatePaymentOrderProcessedCommand{
				OrderID:           "order-1",
				PaymentID:         "payment-1",
				AuthorizationCode: "auth-123",
				OrderStatusString: "completed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCreatePaymentOrderProcessedCommand(tt.orderID, tt.paymentID, tt.authCode, tt.status, CardData{})
			assert.Equal(t, tt.expected, got)
		})
	}
}

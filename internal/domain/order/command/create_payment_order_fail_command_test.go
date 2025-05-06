package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCreatePaymentOrderFailCommand(t *testing.T) {
	tests := []struct {
		name      string
		orderID   string
		paymentID string
		reason    string
		status    string
		expected  CreatePaymentOrderFailCommand
	}{
		{
			name:      "create failed payment command",
			orderID:   "order-1",
			paymentID: "payment-1",
			reason:    "insufficient funds",
			status:    "failed",
			expected: CreatePaymentOrderFailCommand{
				OrderID:           "order-1",
				PaymentID:         "payment-1",
				PaymentReason:     "insufficient funds",
				OrderStatusString: "failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCreatePaymentOrderFailCommand(tt.orderID, tt.paymentID, tt.reason, tt.status)
			assert.Equal(t, tt.expected, got)
		})
	}
}

package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaymentOrderReleasedCommand(t *testing.T) {
	tests := []struct {
		name      string
		orderID   string
		paymentID string
		reason    string
		expected  PaymentOrderReleasedCommand
	}{
		{
			name:      "create payment released command",
			orderID:   "order-1",
			paymentID: "payment-1",
			reason:    "payment released",
			expected: PaymentOrderReleasedCommand{
				OrderID:   "order-1",
				PaymentID: "payment-1",
				Reason:    "payment released",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPaymentOrderReleasedCommand(tt.orderID, tt.paymentID, tt.reason)
			assert.Equal(t, tt.expected, got)
		})
	}
}

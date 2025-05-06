package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

func TestFromFailedOrderCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      command.CreatePaymentOrderFailCommand
		expected *OrderPaymentFailedEvent
	}{
		{
			name: "create failed order event",
			cmd: command.CreatePaymentOrderFailCommand{
				OrderID:           "order-1",
				PaymentID:         "payment-1",
				PaymentReason:     "insufficient funds",
				OrderStatusString: "failed",
			},
			expected: &OrderPaymentFailedEvent{
				OrderID:           "order-1",
				PaymentID:         "payment-1",
				PaymentReason:     "insufficient funds",
				OrderStatusString: "failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromFailedOrderCommand(tt.cmd)
			assert.Equal(t, tt.expected, got)
		})
	}
}

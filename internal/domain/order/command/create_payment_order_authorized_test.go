package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCreatePaymentOrderAuthorizedCommand(t *testing.T) {
	tests := []struct {
		name              string
		orderID           string
		paymentID         string
		authorizationCode string
		orderStatusString string
		paymentCard       CardData
		expected          CreatePaymentOrderAuthorizedCommand
	}{
		{
			name:              "create authorized payment command",
			orderID:           "order-1",
			paymentID:         "payment-1",
			authorizationCode: "auth-123",
			orderStatusString: "authorized",
			paymentCard:       CardData{},
			expected: CreatePaymentOrderAuthorizedCommand{
				OrderID:           "order-1",
				PaymentID:         "payment-1",
				AuthorizationCode: "auth-123",
				OrderStatusString: "authorized",
				PaymentCard:       CardData{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCreatePaymentOrderAuthorizedCommand(
				tt.orderID,
				tt.paymentID,
				tt.authorizationCode,
				tt.orderStatusString,
				tt.paymentCard,
			)
			assert.Equal(t, tt.expected, got)
		})
	}
}

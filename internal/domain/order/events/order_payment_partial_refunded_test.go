package events

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestFromRefundPartialPaymentCommand(t *testing.T) {
	tests := []struct {
		name         string
		paymentID    string
		reason       string
		amount       decimal.Decimal
		expectAmount decimal.Decimal
	}{
		{
			name:         "valid payment data",
			paymentID:    "payment123",
			reason:       "Product defect",
			amount:       decimal.NewFromFloat(50.75),
			expectAmount: decimal.NewFromFloat(50.75),
		},
		{
			name:         "empty payment ID",
			paymentID:    "",
			reason:       "Refund mistake",
			amount:       decimal.NewFromFloat(0),
			expectAmount: decimal.NewFromFloat(0),
		},
		{
			name:         "empty reason",
			paymentID:    "payment456",
			reason:       "",
			amount:       decimal.NewFromInt(100),
			expectAmount: decimal.NewFromInt(100),
		},
		{
			name:         "zero amount",
			paymentID:    "payment789",
			reason:       "No valid refund request",
			amount:       decimal.NewFromInt(0),
			expectAmount: decimal.NewFromInt(0),
		},
		{
			name:         "negative amount",
			paymentID:    "payment000",
			reason:       "Invalid data entry",
			amount:       decimal.NewFromFloat(-5),
			expectAmount: decimal.NewFromFloat(-5),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			evt := FromRefundPartialPaymentCommand(tc.paymentID, tc.reason, tc.amount)

			assert.NotNil(t, evt)
			assert.Equal(t, tc.paymentID, evt.PaymentID)
			assert.Equal(t, tc.reason, evt.Reason)
			assert.True(t, evt.Amount.Equal(tc.expectAmount), "expected amount to match")
			assert.WithinDuration(t, time.Now(), evt.CreatedAt, time.Second, "createdAt should be within 1 second")
		})
	}
}

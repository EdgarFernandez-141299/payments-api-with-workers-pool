package command

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewCreatePartialPaymentRefundCommand(t *testing.T) {
	// Arrange
	referenceOrderID := "order-123"
	paymentOrderID := "payment-456"
	amount := decimal.NewFromFloat(10.50)
	reason := "customer request"

	// Act
	cmd := NewCreatePartialPaymentRefundCommand(referenceOrderID, paymentOrderID, amount, reason)

	// Assert
	assert.Equal(t, referenceOrderID, cmd.ReferenceOrderID)
	assert.Equal(t, paymentOrderID, cmd.PaymentOrderID)
	assert.Equal(t, amount, cmd.Amount)
	assert.Equal(t, reason, cmd.Reason)
}

func TestCreatePartialPaymentRefundCommand_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cmd     CreatePartialPaymentRefundCommand
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid command",
			cmd: CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Amount:           decimal.NewFromFloat(10.50),
				Reason:           "customer request",
			},
			wantErr: false,
		},
		{
			name: "missing reference order id",
			cmd: CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "",
				PaymentOrderID:   "payment-456",
				Amount:           decimal.NewFromFloat(10.50),
				Reason:           "customer request",
			},
			wantErr: true,
			errMsg:  "reference order id is required",
		},
		{
			name: "missing payment order id",
			cmd: CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "",
				Amount:           decimal.NewFromFloat(10.50),
				Reason:           "customer request",
			},
			wantErr: true,
			errMsg:  "payment order id is required",
		},
		{
			name: "zero amount",
			cmd: CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Amount:           decimal.Zero,
				Reason:           "customer request",
			},
			wantErr: true,
			errMsg:  "amount is required",
		},
		{
			name: "missing reason",
			cmd: CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Amount:           decimal.NewFromFloat(10.50),
				Reason:           "",
			},
			wantErr: true,
			errMsg:  "reason is required",
		},
		{
			name: "negative amount",
			cmd: CreatePartialPaymentRefundCommand{
				ReferenceOrderID: "order-123",
				PaymentOrderID:   "payment-456",
				Amount:           decimal.NewFromFloat(-10.50),
				Reason:           "customer request",
			},
			wantErr: true,
			errMsg:  "amount cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
)

func TestNewRefundTotalCommand(t *testing.T) {
	// Arrange
	referenceOrderID := "ref-123"
	paymentOrderID := "payment-123"
	reason := "customer request"

	// Act
	cmd := NewRefundTotalCommand(referenceOrderID, paymentOrderID, reason)

	// Assert
	assert.Equal(t, referenceOrderID, cmd.ReferenceOrderID)
	assert.Equal(t, paymentOrderID, cmd.PaymentOrderID)
	assert.Equal(t, reason, cmd.Reason)
}

func TestRefundTotalCommand_Validate(t *testing.T) {
	tests := []struct {
		name          string
		command       RefundTotalCommand
		expectedError error
	}{
		{
			name: "valid command",
			command: RefundTotalCommand{
				ReferenceOrderID: "ref-123",
				PaymentOrderID:   "payment-123",
				Reason:           "customer request",
			},
			expectedError: nil,
		},
		{
			name: "missing reference order ID",
			command: RefundTotalCommand{
				ReferenceOrderID: "",
				PaymentOrderID:   "payment-123",
				Reason:           "customer request",
			},
			expectedError: errors.NewRefundCreateValidationError(errors.ErrInvalidRefundReferenceOrderID),
		},
		{
			name: "missing payment order ID",
			command: RefundTotalCommand{
				ReferenceOrderID: "ref-123",
				PaymentOrderID:   "",
				Reason:           "customer request",
			},
			expectedError: errors.NewRefundCreateValidationError(errors.ErrInvalidRefundPaymentOrderID),
		},
		{
			name: "missing reason",
			command: RefundTotalCommand{
				ReferenceOrderID: "ref-123",
				PaymentOrderID:   "payment-123",
				Reason:           "",
			},
			expectedError: errors.NewRefundCreateValidationError(errors.ErrInvalidRefundReason),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.command.Validate()

			// Assert
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}
}

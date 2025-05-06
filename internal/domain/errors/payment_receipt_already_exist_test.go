package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

func TestNewPaymentReceiptAlreadyExistError(t *testing.T) {
	tests := []struct {
		name       string
		paymentID  string
		expectCode string
		expectMsg  string
	}{
		{
			name:       "valid payment ID",
			paymentID:  "12345",
			expectCode: paymentReceiptAlreadyExistErrorCode,
			expectMsg:  "payment receipt for payment 12345 already exists",
		},
		{
			name:       "empty payment ID",
			paymentID:  "",
			expectCode: paymentReceiptAlreadyExistErrorCode,
			expectMsg:  "payment receipt for payment  already exists",
		},
		{
			name:       "special characters in payment ID",
			paymentID:  "payment@#$%",
			expectCode: paymentReceiptAlreadyExistErrorCode,
			expectMsg:  "payment receipt for payment payment@#$% already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPaymentReceiptAlreadyExistError(tt.paymentID)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectMsg)
			assert.True(t, domain.IsBusinessErrorCode(err, tt.expectCode))
		})
	}
}

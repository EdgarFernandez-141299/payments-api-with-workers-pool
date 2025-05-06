package enums_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

func TestNewPaymentMethodsFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    enums.PaymentMethodEnum
		expectError bool
	}{
		{
			name:        "Valid CCMethod",
			input:       "CCData",
			expected:    enums.CCMethod,
			expectError: false,
		},
		{
			name:        "Valid PaymentDevice",
			input:       "DEVICE",
			expected:    enums.PaymentDevice,
			expectError: false,
		},
		{
			name:        "Valid TokenCard",
			input:       "TOKEN_CARD",
			expected:    enums.TokenCard,
			expectError: false,
		},
		{
			name:        "Invalid Payment Method",
			input:       "INVALID",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := enums.NewPaymentMethodsFromString(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, enums.ErrInvalidPaymentMethod, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestPaymentMethodEnum_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    enums.PaymentMethodEnum
		expected bool
	}{
		{
			name:     "Valid CCMethod",
			input:    enums.CCMethod,
			expected: true,
		},
		{
			name:     "Valid PaymentDevice",
			input:    enums.PaymentDevice,
			expected: true,
		},
		{
			name:     "Valid TokenCard",
			input:    enums.TokenCard,
			expected: true,
		},
		{
			name:     "Invalid Payment Method",
			input:    "INVALID",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestPaymentMethodEnum_String(t *testing.T) {
	tests := []struct {
		name     string
		input    enums.PaymentMethodEnum
		expected string
	}{
		{
			name:     "CCMethod String",
			input:    enums.CCMethod,
			expected: "CCData",
		},
		{
			name:     "PaymentDevice String",
			input:    enums.PaymentDevice,
			expected: "DEVICE",
		},
		{
			name:     "TokenCard String",
			input:    enums.TokenCard,
			expected: "TOKEN_CARD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.String())
		})
	}
}

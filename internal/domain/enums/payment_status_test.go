package enums

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaymentStatusFromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    PaymentStatus
		expectError bool
	}{
		{
			name:        "Valid PaymentProcessing",
			input:       "PROCESSING",
			expected:    PaymentProcessing,
			expectError: false,
		},
		{
			name:        "Valid PaymentProcessed",
			input:       "PROCESSED",
			expected:    PaymentProcessed,
			expectError: false,
		},
		{
			name:        "Valid PaymentVoided",
			input:       "VOIDED",
			expected:    PaymentCanceled,
			expectError: false,
		},
		{
			name:        "Valid PaymentDenied",
			input:       "DENIED",
			expected:    PaymentDenied,
			expectError: false,
		},
		{
			name:        "Valid PaymentRefunded",
			input:       "REFUNDED",
			expected:    PaymentRefunded,
			expectError: false,
		},
		{
			name:        "Valid PaymentFailed",
			input:       "FAILED",
			expected:    PaymentFailed,
			expectError: false,
		},
		{
			name:        "Valid PaymentNotProcessed",
			input:       "NOT_PROCESSED",
			expected:    PaymentNotProcessed,
			expectError: false,
		},
		{
			name:        "Valid PartiallyRefunded",
			input:       "PARTIALLY_REFUNDED",
			expected:    PartiallyRefunded,
			expectError: false,
		},
		{
			name:        "Valid PaymentAuthorized",
			input:       "AUTHORIZED",
			expected:    PaymentAuthorized,
			expectError: false,
		},
		{
			name:        "Invalid Payment Status",
			input:       "INVALID",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewPaymentStatusFromString(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidPaymentStatus, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestPaymentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected bool
	}{
		{
			name:     "Valid PaymentProcessing",
			input:    PaymentProcessing,
			expected: true,
		},
		{
			name:     "Valid PaymentProcessed",
			input:    PaymentProcessed,
			expected: true,
		},
		{
			name:     "Invalid PaymentVoided",
			input:    PaymentVoided,
			expected: false,
		},
		{
			name:     "Invalid PaymentDenied",
			input:    PaymentDenied,
			expected: false,
		},
		{
			name:     "Valid PaymentRefunded",
			input:    PaymentRefunded,
			expected: true,
		},
		{
			name:     "Invalid PaymentFailed",
			input:    PaymentFailed,
			expected: false,
		},
		{
			name:     "Invalid PaymentNotProcessed",
			input:    PaymentNotProcessed,
			expected: false,
		},
		{
			name:     "Invalid PartiallyRefunded",
			input:    PartiallyRefunded,
			expected: false,
		},
		{
			name:     "Invalid Empty Status",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsValid())
		})
	}
}

func TestPaymentStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected string
	}{
		{
			name:     "PaymentProcessing String",
			input:    PaymentProcessing,
			expected: "PROCESSING",
		},
		{
			name:     "PaymentProcessed String",
			input:    PaymentProcessed,
			expected: "PROCESSED",
		},
		{
			name:     "PaymentVoided String",
			input:    PaymentVoided,
			expected: "VOIDED",
		},
		{
			name:     "PaymentDenied String",
			input:    PaymentDenied,
			expected: "DENIED",
		},
		{
			name:     "PaymentRefunded String",
			input:    PaymentRefunded,
			expected: "REFUNDED",
		},
		{
			name:     "PaymentFailed String",
			input:    PaymentFailed,
			expected: "FAILED",
		},
		{
			name:     "PaymentNotProcessed String",
			input:    PaymentNotProcessed,
			expected: "NOT_PROCESSED",
		},
		{
			name:     "PartiallyRefunded String",
			input:    PartiallyRefunded,
			expected: "PARTIALLY_REFUNDED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.String())
		})
	}
}

func TestPaymentStatus_IsFailure(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected bool
	}{
		{
			name:     "Failure PaymentFailed",
			input:    PaymentFailed,
			expected: true,
		},
		{
			name:     "Failure PaymentDenied",
			input:    PaymentDenied,
			expected: true,
		},
		{
			name:     "Failure PaymentVoided",
			input:    PaymentVoided,
			expected: true,
		},
		{
			name:     "Success PaymentProcessing",
			input:    PaymentProcessing,
			expected: false,
		},
		{
			name:     "Success PaymentProcessed",
			input:    PaymentProcessed,
			expected: false,
		},
		{
			name:     "Success PaymentNotProcessed",
			input:    PaymentNotProcessed,
			expected: false,
		},
		{
			name:     "Success PartiallyRefunded",
			input:    PartiallyRefunded,
			expected: false,
		},
		{
			name:     "Nonexistent InvalidStatus",
			input:    PaymentStatus("INVALID_STATUS"),
			expected: false,
		},
		{
			name:     "Empty Status",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsFailure())
		})
	}
}

func TestPaymentStatus_IsProcessing(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected bool
	}{
		{
			name:     "Is Processing",
			input:    PaymentProcessing,
			expected: true,
		},
		{
			name:     "Is Not Processing",
			input:    PaymentProcessed,
			expected: false,
		},
		{
			name:     "Empty Status",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsProcessing())
		})
	}
}

func TestPaymentStatus_IsProcessed(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected bool
	}{
		{
			name:     "Is Processed",
			input:    PaymentProcessed,
			expected: true,
		},
		{
			name:     "Is Not Processed",
			input:    PaymentProcessing,
			expected: false,
		},
		{
			name:     "Empty Status",
			input:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsProcessed())
		})
	}
}

func TestPaymentStatus_IsPartiallyRefunded(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected bool
	}{
		{
			name:     "Is Partially Refunded",
			input:    PartiallyRefunded,
			expected: true,
		},
		{
			name:     "Is Not Partially Refunded - Processed",
			input:    PaymentProcessed,
			expected: false,
		},
		{
			name:     "Is Not Partially Refunded - Denied",
			input:    PaymentDenied,
			expected: false,
		},
		{
			name:     "Empty Status",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid Payment Status",
			input:    PaymentStatus("INVALID_STATUS"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsPartiallyRefunded())
		})
	}
}

func TestPaymentStatus_IsCanceled(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected bool
	}{
		{
			name:     "Is Canceled - PaymentCanceled",
			input:    PaymentCanceled,
			expected: true,
		},
		{
			name:     "Is Canceled - PaymentVoided",
			input:    PaymentVoided,
			expected: true,
		},
		{
			name:     "Is Not Canceled - Processed",
			input:    PaymentProcessed,
			expected: false,
		},
		{
			name:     "Is Not Canceled - Processing",
			input:    PaymentProcessing,
			expected: false,
		},
		{
			name:     "Empty Status",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid Payment Status",
			input:    PaymentStatus("INVALID_STATUS"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsCanceled())
		})
	}
}

func TestPaymentStatus_IsAuthorized(t *testing.T) {
	tests := []struct {
		name     string
		input    PaymentStatus
		expected bool
	}{
		{
			name:     "Is Authorized",
			input:    PaymentAuthorized,
			expected: true,
		},
		{
			name:     "Is Not Authorized - Processed",
			input:    PaymentProcessed,
			expected: false,
		},
		{
			name:     "Is Not Authorized - Processing",
			input:    PaymentProcessing,
			expected: false,
		},
		{
			name:     "Empty Status",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid Payment Status",
			input:    PaymentStatus("INVALID_STATUS"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.input.IsAuthorized())
		})
	}
}

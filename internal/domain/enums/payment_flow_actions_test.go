package enums

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentFlowActionEnum_String(t *testing.T) {
	tests := []struct {
		name     string
		action   PaymentFlowActionEnum
		expected string
	}{
		{
			name:     "Capture action",
			action:   CapturePayment,
			expected: "CAPTURE",
		},
		{
			name:     "Release action",
			action:   ReleasePayment,
			expected: "RELEASE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.action.String())
		})
	}
}

func TestNewPaymentFlowActionEnum(t *testing.T) {
	tests := []struct {
		name        string
		action      string
		expected    PaymentFlowActionEnum
		expectError bool
	}{
		{
			name:        "Valid capture action",
			action:      "CAPTURE",
			expected:    CapturePayment,
			expectError: false,
		},
		{
			name:        "Valid release action",
			action:      "RELEASE",
			expected:    ReleasePayment,
			expectError: false,
		},
		{
			name:        "Invalid action",
			action:      "INVALID",
			expected:    "",
			expectError: true,
		},
		{
			name:        "Case insensitive capture",
			action:      "capture",
			expected:    CapturePayment,
			expectError: false,
		},
		{
			name:        "Case insensitive release",
			action:      "release",
			expected:    ReleasePayment,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewPaymentFlowActionEnum(tt.action)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidPaymentFlowAction, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestPaymentFlowActionEnum_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		action   PaymentFlowActionEnum
		expected bool
	}{
		{
			name:     "Valid capture action",
			action:   CapturePayment,
			expected: true,
		},
		{
			name:     "Valid release action",
			action:   ReleasePayment,
			expected: true,
		},
		{
			name:     "Invalid action",
			action:   PaymentFlowActionEnum("INVALID"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.action.IsValid())
		})
	}
}

func TestPaymentFlowActionEnum_IsCapture(t *testing.T) {
	tests := []struct {
		name     string
		action   PaymentFlowActionEnum
		expected bool
	}{
		{
			name:     "Is capture",
			action:   CapturePayment,
			expected: true,
		},
		{
			name:     "Is not capture",
			action:   ReleasePayment,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.action.IsCapture())
		})
	}
}

func TestPaymentFlowActionEnum_IsRelease(t *testing.T) {
	tests := []struct {
		name     string
		action   PaymentFlowActionEnum
		expected bool
	}{
		{
			name:     "Is release",
			action:   ReleasePayment,
			expected: true,
		},
		{
			name:     "Is not release",
			action:   CapturePayment,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.action.IsRelease())
		})
	}
}

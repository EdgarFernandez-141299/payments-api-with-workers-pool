package value_objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderStatusFailed(t *testing.T) {
	status := OrderStatusFailed()
	assert.Equal(t, failed, status.Get())
}

func TestOrderStatusProcessing(t *testing.T) {
	status := OrderStatusProcessing()
	assert.Equal(t, processing, status.Get())
}

func TestOrderStatusRefunded(t *testing.T) {
	status := OrderStatusRefunded()
	assert.Equal(t, refunded, status.Get())
}

func TestOrderStatusPartiallyRefunded(t *testing.T) {
	status := OrderStatusPartiallyRefunded()
	assert.Equal(t, partiallyRefunded, status.Get())
}

func TestOrderStatusPartiallyProcessed(t *testing.T) {
	status := OrderPartialProcessed()
	assert.Equal(t, partiallyProcessed, status.Get())
}

func TestOrderStatusProcessed(t *testing.T) {
	status := OrderStatusProcessed()
	assert.Equal(t, processed, status.Get())
}

func TestOrderStatusAuthorized(t *testing.T) {
	status := OrderStatusAuthorized()
	assert.Equal(t, authorized, status.Get())
}

func TestOrderStatusCanceled(t *testing.T) {
	status := OrderStatusCanceled()
	assert.Equal(t, canceled, status.Get())
}

func TestOrderStatus_Get(t *testing.T) {
	tests := []struct {
		name     string
		status   OrderStatus
		expected string
	}{
		{
			name:     "processing status",
			status:   OrderStatus{status: processing},
			expected: processing,
		},
		{
			name:     "failed status",
			status:   OrderStatus{status: failed},
			expected: failed,
		},
		{
			name:     "refunded status",
			status:   OrderStatus{status: refunded},
			expected: refunded,
		},
		{
			name:     "partially refunded status",
			status:   OrderStatus{status: partiallyRefunded},
			expected: partiallyRefunded,
		},
		{
			name:     "partially processed status",
			status:   OrderStatus{status: partiallyProcessed},
			expected: partiallyProcessed,
		},
		{
			name:     "processed status",
			status:   OrderStatus{status: processed},
			expected: processed,
		},
		{
			name:     "authorized status",
			status:   OrderStatus{status: authorized},
			expected: authorized,
		},
		{
			name:     "canceled status",
			status:   OrderStatus{status: canceled},
			expected: canceled,
		},
		{
			name:     "custom status",
			status:   OrderStatus{status: "CUSTOM_STATUS"},
			expected: "CUSTOM_STATUS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.Get()
			assert.Equal(t, tt.expected, result)
		})
	}
}

package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReleaseResponse(t *testing.T) {
	tests := []struct {
		name             string
		referenceOrderID string
		paymentID        string
		paymentStatus    string
		expected         *ReleaseResponse
	}{
		{
			name:             "Caso exitoso",
			referenceOrderID: "12345",
			paymentID:        "PAY-123",
			paymentStatus:    "COMPLETED",
			expected: &ReleaseResponse{
				ReferenceOrderID: "12345",
				PaymentID:        "PAY-123",
				PaymentStatus:    "COMPLETED",
			},
		},
		{
			name:             "Campos vacíos",
			referenceOrderID: "",
			paymentID:        "",
			paymentStatus:    "",
			expected: &ReleaseResponse{
				ReferenceOrderID: "",
				PaymentID:        "",
				PaymentStatus:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewReleaseResponse(tt.referenceOrderID, tt.paymentID, tt.paymentStatus)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

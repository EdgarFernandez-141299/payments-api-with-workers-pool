package events

import (
	"testing"
	"time"
)

func TestFromRefundOrderCommand(t *testing.T) {
	tests := []struct {
		name      string
		paymentID string
		reason    string
	}{
		{"valid payment and reason", "12345", "Duplicate payment"},
		{"empty payment ID", "", "Reason with empty payment ID"},
		{"empty reason", "12345", ""},
		{"both empty", "", ""},
		{"special characters", "123$%^", "Reason with special chars #@$"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromRefundOrderCommand(tt.paymentID, tt.reason)

			if result == nil {
				t.Errorf("expected non-nil result, got nil")
				return
			}

			if result.PaymentID != tt.paymentID {
				t.Errorf("expected PaymentID: %s, got: %s", tt.paymentID, result.PaymentID)
			}

			if result.Reason != tt.reason {
				t.Errorf("expected Reason: %s, got: %s", tt.reason, result.Reason)
			}

			if time.Since(result.CreatedAt) < 0 || time.Since(result.CreatedAt) > time.Second {
				t.Errorf("expected CreatedAt to be current time, got: %v", result.CreatedAt)
			}
		})
	}
}

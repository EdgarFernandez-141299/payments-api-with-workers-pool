package integration_events

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewOrderVoidedPaidIntegrationEvent(t *testing.T) {
	tests := []struct {
		name       string
		params     IntegrationEventsParams
		reason     string
		wantReason string
	}{
		{
			name: "valid input with non-empty reason",
			params: IntegrationEventsParams{
				ReferenceOrderID:   "order123",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assocPay123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "regular",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			reason:     "Payment voided due to incorrect amount",
			wantReason: "Payment voided due to incorrect amount",
		},
		{
			name: "valid input with empty reason",
			params: IntegrationEventsParams{
				ReferenceOrderID:   "order123",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assocPay123",
				TotalOrderAmount:   50.0,
				Currency:           "EUR",
				UserID:             "user456",
				UserType:           "guest",
				EnterpriseID:       "enterprise456",
				TotalOrderPaid:     50.0,
				TotalPaymentAmount: 50.0,
			},
			reason:     "",
			wantReason: "",
		},
		{
			name: "valid input with special characters in reason",
			params: IntegrationEventsParams{
				ReferenceOrderID:   "order789",
				ReferencePaymentID: "payment789",
				AssociatedPayment:  "assocPay789",
				TotalOrderAmount:   200.0,
				Currency:           "GBP",
				UserID:             "user789",
				UserType:           "premium",
				EnterpriseID:       "enterprise789",
				TotalOrderPaid:     200.0,
				TotalPaymentAmount: 200.0,
			},
			reason:     "Reason with special characters: !@#&*()",
			wantReason: "Reason with special characters: !@#&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewOrderFailedIntegrationEvent(tt.params, tt.reason, "failed")
			assert.Equal(t, got.Type, paymentOrderFailedEventType)
			if !reflect.DeepEqual(got.FailureReason, tt.wantReason) {
				t.Errorf("NewOrderFailedIntegrationEvent().FailureReason = %v, want %v", got.FailureReason, tt.wantReason)
			}
		})
	}
}

package integration_events

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewOrderPaymentPaidIntegrationEvent(t *testing.T) {
	tests := []struct {
		name                string
		params              IntegrationEventsParams
		authorizationCode   string
		orderStatus         string
		expectedAuthCode    string
		expectedOrderStatus string
	}{
		{
			name: "Valid creation",
			params: IntegrationEventsParams{
				ReferenceOrderID:   "order123",
				ReferencePaymentID: "payment456",
				AssociatedPayment:  "associatedPayment789",
				TotalOrderAmount:   120.5,
				Currency:           "USD",
				UserID:             "user001",
				UserType:           "regular",
				EnterpriseID:       "enterprise567",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 20.5,
			},
			authorizationCode:   "authCode123",
			orderStatus:         "PAID",
			expectedAuthCode:    "authCode123",
			expectedOrderStatus: "PAID",
		},
		{
			name:                "Empty params",
			params:              IntegrationEventsParams{},
			authorizationCode:   "",
			orderStatus:         "PENDING",
			expectedAuthCode:    "",
			expectedOrderStatus: "PENDING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewOrderPaymentProcessedIntegrationEvent(tt.params, tt.authorizationCode, tt.orderStatus)
			assert.Equal(t, event.Type, paymentOrderProcessedEventType)

			if event.AuthorizationCode != tt.expectedAuthCode {
				t.Errorf("AuthorizationCode = %v; want %v", event.AuthorizationCode, tt.expectedAuthCode)
			}

			if !reflect.DeepEqual(event.baseOrderIntegrationEvent.OrderStatus, tt.expectedOrderStatus) {
				t.Errorf("OrderStatus = %v; want %v", event.baseOrderIntegrationEvent.OrderStatus, tt.expectedOrderStatus)
			}
		})
	}
}

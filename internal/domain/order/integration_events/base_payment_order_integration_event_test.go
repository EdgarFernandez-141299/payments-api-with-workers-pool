package integration_events

import (
	"errors"
	"testing"
)

func TestBaseOrderIntegrationEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   baseOrderIntegrationEvent
		wantErr error
	}{
		{
			name: "valid event",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
				Metadata: map[string]interface{}{
					"key1": "value1",
				},
			},
			wantErr: nil,
		},
		{
			name: "missing ReferenceOrderID",
			event: baseOrderIntegrationEvent{
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing PaymentStatus",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing OrderStatus",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing ReferencePaymentID",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing AssociatedPayment",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing Currency",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				UserID:             "user123",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing UserID",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserType:           "customer",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing UserType",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				EnterpriseID:       "enterprise123",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
		{
			name: "missing EnterpriseID",
			event: baseOrderIntegrationEvent{
				ReferenceOrderID:   "order123",
				PaymentStatus:      "paid",
				OrderStatus:        "completed",
				ReferencePaymentID: "payment123",
				AssociatedPayment:  "assoc123",
				TotalOrderAmount:   100.0,
				Currency:           "USD",
				UserID:             "user123",
				UserType:           "customer",
				TotalOrderPaid:     100.0,
				TotalPaymentAmount: 100.0,
			},
			wantErr: integrationEventParamError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

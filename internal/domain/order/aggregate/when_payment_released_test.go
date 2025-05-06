package aggregate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestWhenPaymentReleased(t *testing.T) {
	tests := []struct {
		name           string
		order          *Order
		event          events.OrderPaymentReleasedEvent
		expectedStatus value_objects.OrderStatus
	}{
		{
			name: "should update order status to canceled and payment status to canceled",
			order: &Order{
				Status: value_objects.OrderStatusProcessing(),
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
				},
			},
			event: events.OrderPaymentReleasedEvent{
				OrderID:    "order-1",
				PaymentID:  "payment-1",
				Reason:     "test reason",
				ReleasedAt: time.Now(),
			},
			expectedStatus: value_objects.OrderStatusCanceled(),
		},
		{
			name: "should not update payment status if payment not found",
			order: &Order{
				Status: value_objects.OrderStatusProcessing(),
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-2",
						Status: enums.PaymentProcessing,
					},
				},
			},
			event: events.OrderPaymentReleasedEvent{
				OrderID:    "order-1",
				PaymentID:  "payment-1",
				Reason:     "test reason",
				ReleasedAt: time.Now(),
			},
			expectedStatus: value_objects.OrderStatusCanceled(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WhenPaymentReleased(tt.order, tt.event)

			assert.Equal(t, tt.expectedStatus.Get(), tt.order.Status.Get())

			for _, payment := range tt.order.OrderPayments {
				if payment.ID == tt.event.PaymentID {
					assert.Equal(t, enums.PaymentCanceled, payment.Status)
				}
			}
		})
	}
}

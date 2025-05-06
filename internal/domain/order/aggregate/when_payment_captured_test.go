package aggregate

import (
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

	"github.com/stretchr/testify/assert"
)

func TestWhenPaymentCaptured(t *testing.T) {
	tests := []struct {
		name           string
		initialOrder   *Order
		event          events.OrderPaymentCapturedEvent
		expectedStatus value_objects.OrderStatus
	}{
		{
			name: "should update order and payment status when payment is captured",
			initialOrder: &Order{
				Status: value_objects.OrderStatusProcessing(),
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentProcessing,
					},
					{
						ID:     "payment-2",
						Status: enums.PaymentProcessing,
					},
				},
			},
			event: events.OrderPaymentCapturedEvent{
				PaymentID: "payment-1",
			},
			expectedStatus: value_objects.OrderStatusProcessed(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WhenPaymentCaptured(tt.initialOrder, tt.event)

			assert.Equal(t, tt.expectedStatus, tt.initialOrder.Status)

			// Verificar que el pago específico fue actualizado
			for _, payment := range tt.initialOrder.OrderPayments {
				if payment.ID == tt.event.PaymentID {
					assert.Equal(t, enums.PaymentProcessed, payment.Status)
				}
			}
		})
	}
}

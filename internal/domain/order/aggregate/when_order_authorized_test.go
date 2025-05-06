package aggregate

import (
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

	"github.com/stretchr/testify/assert"
)

func TestWhenOrderPaymentAuthorized(t *testing.T) {
	tests := []struct {
		name                string
		order               *Order
		event               events.OrderPaymentAuthorizedEvent
		expectedStatus      enums.PaymentStatus
		expectedOrderStatus vo.OrderStatus
		shouldError         bool
	}{
		{
			name: "should update payment status to authorized",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-123",
						Status: enums.PaymentStatus("PENDING"),
					},
				},
			},
			event: events.OrderPaymentAuthorizedEvent{
				PaymentID:         "payment-123",
				AuthorizationCode: "AUTH123",
				PaymentCard: command.CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "CREDIT",
				},
			},
			expectedStatus:      enums.PaymentStatus("AUTHORIZED"),
			expectedOrderStatus: vo.OrderStatusAuthorized(),
			shouldError:         false,
		},
		{
			name: "should handle multiple payments in order",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-1",
						Status: enums.PaymentStatus("PENDING"),
					},
					{
						ID:     "payment-2",
						Status: enums.PaymentStatus("PENDING"),
					},
				},
			},
			event: events.OrderPaymentAuthorizedEvent{
				PaymentID:         "payment-2",
				AuthorizationCode: "AUTH456",
				PaymentCard: command.CardData{
					CardBrand: "MASTERCARD",
					CardLast4: "5678",
					CardType:  "DEBIT",
				},
			},
			expectedStatus:      enums.PaymentStatus("AUTHORIZED"),
			expectedOrderStatus: vo.OrderStatusAuthorized(),
			shouldError:         false,
		},
		{
			name: "should not find non-existent payment",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID:     "payment-123",
						Status: enums.PaymentStatus("PENDING"),
					},
				},
			},
			event: events.OrderPaymentAuthorizedEvent{
				PaymentID:         "non-existent",
				AuthorizationCode: "AUTH789",
				PaymentCard: command.CardData{
					CardBrand: "AMEX",
					CardLast4: "9012",
					CardType:  "CREDIT",
				},
			},
			expectedStatus:      enums.PaymentStatus("PENDING"),
			expectedOrderStatus: vo.OrderStatusAuthorized(),
			shouldError:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WhenOrderPaymentAuthorized(tt.order, tt.event)

			paymentFound := false
			for _, payment := range tt.order.OrderPayments {
				if payment.ID == tt.event.PaymentID {
					paymentFound = true
					assert.Equal(t, tt.expectedStatus, payment.Status)
					assert.Equal(t, tt.event.AuthorizationCode, payment.AuthorizationCode)
					assert.Equal(t, tt.event.PaymentCard.CardBrand, payment.PaymentCard.CardBrand)
					assert.Equal(t, tt.event.PaymentCard.CardLast4, payment.PaymentCard.CardLast4)
					assert.Equal(t, tt.event.PaymentCard.CardType, payment.PaymentCard.CardType)
					break
				}
			}

			if tt.shouldError {
				assert.False(t, paymentFound, "Payment should not be found")
			} else {
				assert.True(t, paymentFound, "Payment should be found")
			}

			assert.Equal(t, tt.expectedOrderStatus, tt.order.Status, "Order status should be updated correctly")
		})
	}
}

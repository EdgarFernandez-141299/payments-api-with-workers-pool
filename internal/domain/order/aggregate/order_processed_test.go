package aggregate

import (
	"testing"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/events"
	vo "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

func TestWhenOrderPaymentProcessed(t *testing.T) {
	usdCode := vo.CurrencyCode{Code: "USD"}
	amount := decimal.NewFromFloat(100.0)

	tests := []struct {
		name                string
		order               *Order
		event               events.OrderPaymentProcessedEvent
		expectedStatus      enums.PaymentStatus
		expectedAuthCode    string
		expectedOrderStatus vo.OrderStatus
		expectedCardData    entities.CardData
	}{
		{
			name: "procesamiento exitoso con autorización",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID: "payment-1",
						Total: vo.CurrencyAmount{
							Value: amount,
							Code:  usdCode,
						},
						Status: enums.PaymentProcessing,
					},
				},
			},
			event: events.OrderPaymentProcessedEvent{
				OrderID:           "order-1",
				PaymentID:         "payment-1",
				AuthorizationCode: "auth-123",
				OrderStatusString: "processed",
				PaymentCard: command.CardData{
					CardBrand: "VISA",
					CardLast4: "1234",
					CardType:  "credit",
				},
			},
			expectedStatus:      enums.PaymentProcessed,
			expectedAuthCode:    "auth-123",
			expectedOrderStatus: vo.OrderStatusProcessed(),
			expectedCardData: entities.CardData{
				CardBrand: "VISA",
				CardLast4: "1234",
				CardType:  "credit",
			},
		},
		{
			name: "procesamiento exitoso con múltiples pagos",
			order: &Order{
				ID: "order-1",
				OrderPayments: []entities.PaymentOrder{
					{
						ID: "payment-1",
						Total: vo.CurrencyAmount{
							Value: amount,
							Code:  usdCode,
						},
						Status: enums.PaymentProcessing,
					},
					{
						ID: "payment-2",
						Total: vo.CurrencyAmount{
							Value: amount,
							Code:  usdCode,
						},
						Status: enums.PaymentProcessing,
					},
				},
			},
			event: events.OrderPaymentProcessedEvent{
				OrderID:           "order-1",
				PaymentID:         "payment-2",
				AuthorizationCode: "auth-456",
				OrderStatusString: "processed",
				PaymentCard: command.CardData{
					CardBrand: "MASTERCARD",
					CardLast4: "5678",
					CardType:  "debit",
				},
			},
			expectedStatus:      enums.PaymentProcessed,
			expectedAuthCode:    "auth-456",
			expectedOrderStatus: vo.OrderStatusProcessed(),
			expectedCardData: entities.CardData{
				CardBrand: "MASTERCARD",
				CardLast4: "5678",
				CardType:  "debit",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WhenOrderPaymentProcessed(tt.order, tt.event)

			payment, _ := tt.order.FindPaymentByID(tt.event.PaymentID)
			if payment.Status != tt.expectedStatus {
				t.Errorf("WhenOrderPaymentProcessed() status = %v, want %v", payment.Status, tt.expectedStatus)
			}
			if payment.AuthorizationCode != tt.expectedAuthCode {
				t.Errorf("WhenOrderPaymentProcessed() authCode = %v, want %v", payment.AuthorizationCode, tt.expectedAuthCode)
			}
			if tt.order.Status != tt.expectedOrderStatus {
				t.Errorf("WhenOrderPaymentProcessed() order status = %v, want %v", tt.order.Status, tt.expectedOrderStatus)
			}
			if payment.PaymentCard.CardBrand != tt.expectedCardData.CardBrand {
				t.Errorf("WhenOrderPaymentProcessed() card brand = %v, want %v", payment.PaymentCard.CardBrand, tt.expectedCardData.CardBrand)
			}
			if payment.PaymentCard.CardLast4 != tt.expectedCardData.CardLast4 {
				t.Errorf("WhenOrderPaymentProcessed() card last4 = %v, want %v", payment.PaymentCard.CardLast4, tt.expectedCardData.CardLast4)
			}
			if payment.PaymentCard.CardType != tt.expectedCardData.CardType {
				t.Errorf("WhenOrderPaymentProcessed() card type = %v, want %v", payment.PaymentCard.CardType, tt.expectedCardData.CardType)
			}
		})
	}
}

func TestWhenOrderPaymentFailed(t *testing.T) {
	usdCode := vo.CurrencyCode{Code: "USD"}
	amount := decimal.NewFromFloat(100.0)

	tests := []struct {
		name                string
		order               *Order
		event               events.OrderPaymentFailedEvent
		expectedStatus      enums.PaymentStatus
		expectedReason      string
		expectedOrderStatus vo.OrderStatus
	}{
		{
			name: "fallo de pago con razón específica",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID: "payment-1",
						Total: vo.CurrencyAmount{
							Value: amount,
							Code:  usdCode,
						},
						Status: enums.PaymentProcessing,
					},
				},
			},
			event: events.OrderPaymentFailedEvent{
				PaymentID:     "payment-1",
				PaymentReason: "insufficient_funds",
			},
			expectedStatus:      enums.PaymentFailed,
			expectedReason:      "insufficient_funds",
			expectedOrderStatus: vo.OrderStatusFailed(),
		},
		{
			name: "fallo de pago con múltiples pagos",
			order: &Order{
				OrderPayments: []entities.PaymentOrder{
					{
						ID: "payment-1",
						Total: vo.CurrencyAmount{
							Value: amount,
							Code:  usdCode,
						},
						Status: enums.PaymentProcessing,
					},
					{
						ID: "payment-2",
						Total: vo.CurrencyAmount{
							Value: amount,
							Code:  usdCode,
						},
						Status: enums.PaymentProcessing,
					},
				},
			},
			event: events.OrderPaymentFailedEvent{
				PaymentID:     "payment-1",
				PaymentReason: "card_declined",
			},
			expectedStatus:      enums.PaymentFailed,
			expectedReason:      "card_declined",
			expectedOrderStatus: vo.OrderStatusFailed(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WhenOrderPaymentFailed(tt.order, tt.event)

			payment, _ := tt.order.FindPaymentByID(tt.event.PaymentID)
			if payment.Status != tt.expectedStatus {
				t.Errorf("WhenOrderPaymentFailed() status = %v, want %v", payment.Status, tt.expectedStatus)
			}
			if payment.FailureReason != tt.expectedReason {
				t.Errorf("WhenOrderPaymentFailed() reason = %v, want %v", payment.FailureReason, tt.expectedReason)
			}
			if tt.order.Status != tt.expectedOrderStatus {
				t.Errorf("WhenOrderPaymentFailed() order status = %v, want %v", tt.order.Status, tt.expectedOrderStatus)
			}
		})
	}
}

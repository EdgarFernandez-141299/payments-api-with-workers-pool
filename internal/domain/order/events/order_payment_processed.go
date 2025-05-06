package events

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

type OrderPaymentProcessedEvent struct {
	OrderID           string
	PaymentID         string
	AuthorizationCode string
	OrderStatusString string
	PaymentCard       command.CardData
}

func FromProcessedOrderCommand(cmd command.CreatePaymentOrderProcessedCommand) *OrderPaymentProcessedEvent {
	return &OrderPaymentProcessedEvent{
		OrderID:           cmd.OrderID,
		PaymentID:         cmd.PaymentID,
		AuthorizationCode: cmd.AuthorizationCode,
		OrderStatusString: cmd.OrderStatusString,
		PaymentCard:       cmd.PaymentCard,
	}
}

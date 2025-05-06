package events

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

type OrderPaymentAuthorizedEvent struct {
	OrderID           string
	PaymentID         string
	AuthorizationCode string
	OrderStatusString string
	PaymentCard       command.CardData
}

func FromAuthorizedOrderCommand(cmd command.CreatePaymentOrderAuthorizedCommand) *OrderPaymentAuthorizedEvent {
	return &OrderPaymentAuthorizedEvent{
		OrderID:           cmd.OrderID,
		PaymentID:         cmd.PaymentID,
		AuthorizationCode: cmd.AuthorizationCode,
		OrderStatusString: cmd.OrderStatusString,
		PaymentCard:       cmd.PaymentCard,
	}
}

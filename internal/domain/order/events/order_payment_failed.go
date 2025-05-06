package events

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"

type OrderPaymentFailedEvent struct {
	OrderID           string
	PaymentID         string
	PaymentReason     string
	OrderStatusString string
}

func FromFailedOrderCommand(cmd command.CreatePaymentOrderFailCommand) *OrderPaymentFailedEvent {
	return &OrderPaymentFailedEvent{
		OrderID:           cmd.OrderID,
		PaymentID:         cmd.PaymentID,
		PaymentReason:     cmd.PaymentReason,
		OrderStatusString: cmd.OrderStatusString,
	}
}

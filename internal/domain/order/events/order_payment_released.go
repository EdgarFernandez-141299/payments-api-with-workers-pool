package events

import (
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type OrderPaymentReleasedEvent struct {
	OrderID    string
	PaymentID  string
	Reason     string
	ReleasedAt time.Time
}

func FromReleasedOrderCommand(cmd command.PaymentOrderReleasedCommand) *OrderPaymentReleasedEvent {
	return &OrderPaymentReleasedEvent{
		OrderID:    cmd.OrderID,
		PaymentID:  cmd.PaymentID,
		Reason:     cmd.Reason,
		ReleasedAt: time.Now(),
	}
}

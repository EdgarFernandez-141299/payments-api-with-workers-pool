package events

import (
	"time"
)

type OrderPaymentTotalRefundedEvent struct {
	PaymentID string
	Reason    string
	CreatedAt time.Time
}

func FromRefundOrderCommand(paymentID string, reason string) *OrderPaymentTotalRefundedEvent {
	return &OrderPaymentTotalRefundedEvent{
		PaymentID: paymentID,
		Reason:    reason,
		CreatedAt: time.Now(),
	}
}

package events

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderPaymentPartialRefundedEvent struct {
	PaymentID string
	Reason    string
	Amount    decimal.Decimal
	CreatedAt time.Time
}

func FromRefundPartialPaymentCommand(paymentID string, reason string, amount decimal.Decimal) *OrderPaymentPartialRefundedEvent {
	return &OrderPaymentPartialRefundedEvent{
		PaymentID: paymentID,
		Reason:    reason,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
}

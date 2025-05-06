package events

import (
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type OrderPaymentCapturedEvent struct {
	OrderID    string
	PaymentID  string
	Amount     decimal.Decimal
	CapturedAt time.Time
}

func FromCapturedOrderCommand(cmd command.PaymentOrderCapturedCommand) *OrderPaymentCapturedEvent {
	return &OrderPaymentCapturedEvent{
		OrderID:    cmd.OrderID,
		PaymentID:  cmd.PaymentID,
		Amount:     cmd.Amount,
		CapturedAt: time.Now(),
	}
}

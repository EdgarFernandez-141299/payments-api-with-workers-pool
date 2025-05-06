package command

import (
	"github.com/shopspring/decimal"
)

type PaymentOrderCapturedCommand struct {
	OrderID   string
	PaymentID string
	Amount    decimal.Decimal
	Currency  string
}

func NewPaymentOrderCapturedCommand(orderID, paymentID, currency string, amount decimal.Decimal) PaymentOrderCapturedCommand {
	return PaymentOrderCapturedCommand{
		OrderID:   orderID,
		PaymentID: paymentID,
		Amount:    amount,
		Currency:  currency,
	}
}

package models

import "github.com/shopspring/decimal"

type RefundAdapterModel struct {
	Amount    decimal.Decimal
	PaymentID string
	OrderID   string
}

func NewRefundModel(amount decimal.Decimal, paymentID, orderID string) RefundAdapterModel {
	return RefundAdapterModel{
		Amount:    amount,
		PaymentID: paymentID,
		OrderID:   orderID,
	}
}

type RefundsAdapterModel struct {
	Refunds []RefundAdapterModel
}

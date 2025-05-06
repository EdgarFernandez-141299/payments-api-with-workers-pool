package command

import (
	"errors"

	"github.com/shopspring/decimal"
)

type CreatePartialPaymentRefundCommand struct {
	ReferenceOrderID string
	PaymentOrderID   string
	Amount           decimal.Decimal
	Reason           string
}

func NewCreatePartialPaymentRefundCommand(referenceOrderID, paymentOrderID string, amount decimal.Decimal, reason string) *CreatePartialPaymentRefundCommand {
	return &CreatePartialPaymentRefundCommand{
		ReferenceOrderID: referenceOrderID,
		PaymentOrderID:   paymentOrderID,
		Amount:           amount,
		Reason:           reason,
	}
}

func (cmd *CreatePartialPaymentRefundCommand) Validate() error {
	if cmd.ReferenceOrderID == "" {
		return errors.New("reference order id is required")
	}
	if cmd.PaymentOrderID == "" {
		return errors.New("payment order id is required")
	}
	if cmd.Amount.IsZero() {
		return errors.New("amount is required")
	}
	if cmd.Reason == "" {
		return errors.New("reason is required")
	}
	if cmd.Amount.IsNegative() {
		return errors.New("amount cannot be negative")
	}
	return nil
}

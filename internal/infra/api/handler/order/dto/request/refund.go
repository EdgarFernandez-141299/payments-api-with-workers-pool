package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/command"
)

type RefundDTO struct {
	ReferenceOrderID string          `json:"reference_order_id" validate:"required"`
	PaymentOrderID   string          `json:"payment_order_id" validate:"required"`
	Reason           string          `json:"reason" validate:"required"`
	Amount           decimal.Decimal `json:"amount"`
	IsTotal          bool            `json:"is_total"`
}

func (r *RefundDTO) Validate() error {
	err := validator.New().Struct(r)

	if err != nil {
		return err
	}

	if err := r.ValidateAmount(); err != nil {
		return err
	}

	return nil
}

func (r *RefundDTO) ValidateAmount() error {
	if !r.IsTotal {
		if r.Amount.IsZero() {
			return errors.New("amount cannot be zero")
		}

		if r.Amount.IsNegative() {
			return errors.New("amount cannot be negative")
		}
	}

	return nil
}

func (r *RefundDTO) Command() (command.RefundTotalCommand, error) {
	cmd := command.NewRefundTotalCommand(r.ReferenceOrderID, r.PaymentOrderID, r.Reason)
	return *cmd, nil
}

func (r *RefundDTO) CommandPartial() (command.CreatePartialPaymentRefundCommand, error) {
	cmd := command.NewCreatePartialPaymentRefundCommand(r.ReferenceOrderID, r.PaymentOrderID, r.Amount, r.Reason)
	if err := cmd.Validate(); err != nil {
		return command.CreatePartialPaymentRefundCommand{}, err
	}
	return *cmd, nil
}

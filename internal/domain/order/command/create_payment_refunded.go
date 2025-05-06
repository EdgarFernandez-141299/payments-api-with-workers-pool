package command

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"

type RefundTotalCommand struct {
	ReferenceOrderID string
	PaymentOrderID   string
	Reason           string
}

func NewRefundTotalCommand(referenceOrderID, paymentOrderID, reason string) *RefundTotalCommand {
	return &RefundTotalCommand{
		ReferenceOrderID: referenceOrderID,
		PaymentOrderID:   paymentOrderID,
		Reason:           reason,
	}
}

func (cmd *RefundTotalCommand) Validate() error {
	if cmd.ReferenceOrderID == "" {
		return errors.NewRefundCreateValidationError(errors.ErrInvalidRefundReferenceOrderID)
	}
	if cmd.PaymentOrderID == "" {
		return errors.NewRefundCreateValidationError(errors.ErrInvalidRefundPaymentOrderID)
	}
	if cmd.Reason == "" {
		return errors.NewRefundCreateValidationError(errors.ErrInvalidRefundReason)
	}

	return nil
}

package errors

import (
	"errors"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const refundCreateValidationError = "REFUND_CREATE_VALIDATION_ERROR"

var (
	ErrInvalidRefundReferenceOrderID = errors.New("invalid refund reference order ID")
	ErrInvalidRefundPaymentOrderID   = errors.New("invalid refund payment order ID")
	ErrInvalidRefundReason           = errors.New("invalid refund reason")
)

func NewRefundCreateValidationError(err error) error {
	return domain.WrapBusinessError(err, refundCreateValidationError, "PartialRefund validation error", map[string]interface{}{})
}

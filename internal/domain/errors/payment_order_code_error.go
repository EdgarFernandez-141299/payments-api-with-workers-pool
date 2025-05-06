package errors

import (
	"errors"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const paymentOrderCreateValidationError = "PAYMENT_ORDER_CREATE_VALIDATION_ERROR"

var (
	ErrInvalidPaymentOrderReferenceOrderID  = errors.New("invalid payment order reference order ID")
	ErrInvalidPaymentOrderUserType          = errors.New("invalid payment order user ID")
	ErrInvalidPaymentOrderUserID            = errors.New("invalid payment order user type")
	ErrInvalidPaymentOrderEnterpriseID      = errors.New("invalid payment order enterprise ID")
	ErrInvalidPaymentOrderPaymentMethodType = errors.New("invalid payment order payment method type")
)

func NewPaymentOrderValidationError(err error) error {
	return domain.WrapBusinessError(err, paymentOrderCreateValidationError, "Payment Order validation error", map[string]interface{}{})
}

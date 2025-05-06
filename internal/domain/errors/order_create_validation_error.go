package errors

import (
	"errors"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const orderCreateValidationError = "ORDER_CREATE_VALIDATION_ERROR"

var (
	ErrInvalidOrderID           = errors.New("invalid order ID")
	ErrInvalidOrderTotalAmount  = errors.New("invalid order total amount")
	ErrInvalidOrderPhoneNumber  = errors.New("invalid order phone number")
	ErrInvalidOrderUserType     = errors.New("invalid order user type")
	ErrInvalidOrderEnterpriseID = errors.New("invalid order enterprise ID")
)

func NewOrderCreateValidationError(err error) error {
	return domain.WrapBusinessError(err, orderCreateValidationError, "Order validation error", map[string]interface{}{})
}

package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const orderPaymentNotFoundCodeError = "ORDER_PAYMENT_NOT_FOUND_ERROR"

func NewOrderPaymentNotFoundError(orderID string, paymentID string) error {
	err := fmt.Errorf("no payment found for order ID %s and payment ID %s", orderID, paymentID)

	return domain.WrapBusinessError(err, orderPaymentNotFoundCodeError, err.Error(), map[string]interface{}{})
}

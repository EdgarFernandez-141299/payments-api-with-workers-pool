package errors

import (
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const paymentReceiptAlreadyExistErrorCode = "PAYMENT_RECEIPT_ALREADY_EXIST"

func NewPaymentReceiptAlreadyExistError(paymentID string) error {
	err := fmt.Errorf("payment receipt for payment %v already exists", paymentID)
	return domain.WrapBusinessError(err, paymentReceiptAlreadyExistErrorCode, err.Error(), map[string]interface{}{})
}
package errors

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const invalidAmountErrorCode = "INVALID_AMOUNT"

func NewInvalidAmountError(amount decimal.Decimal) error {
	err := fmt.Errorf("invalid amount: %v", amount)
	return domain.WrapBusinessError(err, invalidAmountErrorCode, err.Error(), map[string]interface{}{})
}

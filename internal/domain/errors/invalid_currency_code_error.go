package errors

import (
	"errors"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

var ErrInvalidCurrencyCode = errors.New("invalid currency code")

const invalidCurrencyCode = "INVALID_CURRENCY_CODE"

func NewInvalidCurrencyCodeError(code string) error {
	return domain.WrapBusinessError(
		ErrInvalidCurrencyCode, invalidCurrencyCode, fmt.Sprintf("Invalid currency code %s", code), map[string]interface{}{},
	)
}

package value_objects

import (
	"github.com/shopspring/decimal"
	bizErrors "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
)

type CurrencyAmount struct {
	Code  CurrencyCode
	Value decimal.Decimal
}

func newAmount(currencyCode CurrencyCode, value decimal.Decimal) CurrencyAmount {
	return CurrencyAmount{
		Code:  currencyCode,
		Value: value,
	}
}

// NewCurrencyAmount creates a new instance of CurrencyAmount with the given currency code and value.
// The provided currency code should exist in the predefined `CurrencySymbols` map.
// If the code is invalid, it returns an error of type `bizErrors.InvalidCurrencyCodeError`.
//
// Parameters:
//   - code: A CurrencyCode representing the currency code (e.g., "USD", "EUR").
//   - value: A float64 representing the numeric value of the currency amount.
//
// Returns:
//   - CurrencyAmount: The created CurrencyAmount object.
//   - error: An error of type `bizErrors.InvalidCurrencyCodeError` if the currency code is invalid; otherwise, nil.
func NewCurrencyAmount(code CurrencyCode, value decimal.Decimal) (CurrencyAmount, error) {
	return newAmount(code, value), nil
}

func (c CurrencyAmount) Add(amount decimal.Decimal) CurrencyAmount {
	newAmount := c.Value.Add(amount)
	newCurrencyAmount, _ := NewCurrencyAmount(c.Code, newAmount)

	return newCurrencyAmount
}

func (c CurrencyAmount) Equals(other CurrencyAmount) bool {
	return c.Code.Equals(other.Code) && c.Value.Equal(other.Value)
}

func (c CurrencyAmount) Validate() error {
	if c.Value.LessThan(decimal.NewFromInt(0)) {
		return bizErrors.NewInvalidAmountError(c.Value)
	}

	return nil
}

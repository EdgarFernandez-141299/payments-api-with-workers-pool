package value_objects

import bizErrors "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"

type CurrencyCode struct {
	Code   string
	Symbol string
}

func NewCurrencyCode(code string) (CurrencyCode, error) {
	symbol, ok := CurrencySymbols[code]
	if !ok {
		return CurrencyCode{}, bizErrors.NewInvalidCurrencyCodeError(code)
	}

	return CurrencyCode{Code: code, Symbol: symbol}, nil
}

func (c CurrencyCode) Equals(other CurrencyCode) bool {
	return c.Code == other.Code
}

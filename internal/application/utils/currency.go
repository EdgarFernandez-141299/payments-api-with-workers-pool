package utils

import "github.com/biter777/countries"

func IsValidCurrencyCode(code string) bool {
	currencyFound := countries.CurrencyCodeByName(code)
	return currencyFound != countries.CurrencyUnknown
}

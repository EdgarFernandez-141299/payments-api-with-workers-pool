package utils

import (
	"fmt"

	"github.com/biter777/countries"
)

func IsValidCountryCode(code string) bool {
	country := countries.ByName(code)
	return country != countries.Unknown && country != countries.None
}

func GetCountryIso3ByCode(code string) string {
	country := countries.ByName(code)
	return country.Alpha3()
}

func GetCountryIso2ByCode(code string) string {
	country := countries.ByName(code)
	return country.Alpha2()
}

func RemoveDuplicateCurrencies(currencies []string) []string {
	uniqueCurrencies := make(map[string]struct{})
	filteredCurrencies := []string{}

	for _, currency := range currencies {
		if _, exists := uniqueCurrencies[currency]; !exists {
			uniqueCurrencies[currency] = struct{}{}
			filteredCurrencies = append(filteredCurrencies, currency)
		}
	}

	return filteredCurrencies
}

func ValidateCurrencies(currencies []string) error {
	for _, currency := range currencies {
		if !IsValidCurrencyCode(currency) {
			return fmt.Errorf("currency %s is not valid", currency)
		}
	}

	return nil
}

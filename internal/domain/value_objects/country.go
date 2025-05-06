package value_objects

import (
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/utils"
	bizErrors "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
)

type Country struct {
	Code string
}

func newCountry(code string) Country {
	return Country{Code: code}
}

func NewCountryWithCode(code string) (Country, error) {
	if code == "" {
		return Country{}, bizErrors.NewInvalidCountryCodeError(code)
	}

	if !utils.IsValidCountryCode(code) {
		return Country{}, bizErrors.NewInvalidCountryCodeError(code)
	}

	return Country{
		Code: utils.GetCountryIso3ByCode(code),
	}, nil
}

func (c Country) Equals(other Country) bool {
	return c.Code == other.Code
}

func (c Country) Iso3() string {
	return c.Code
}

func (c Country) Iso2() string {
	return utils.GetCountryIso2ByCode(c.Code)
}

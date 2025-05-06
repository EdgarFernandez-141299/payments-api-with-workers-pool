package errors

import (
	"errors"
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

var ErrInvalidCountryCode = errors.New("invalid country code")

const invalidCountryCode = "INVALID_COUNTRY_CODE"

func NewInvalidCountryCodeError(code string) error {
	return domain.WrapBusinessError(
		ErrInvalidCountryCode, invalidCountryCode, fmt.Sprintf("Invalid country code %s", code), map[string]interface{}{},
	)
}

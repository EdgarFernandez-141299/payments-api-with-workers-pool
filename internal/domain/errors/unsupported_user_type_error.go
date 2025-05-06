package errors

import (
	"errors"
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

var UnsupportedUserTypeErr = errors.New("unsuported user type")

const unsupportedUserTypeErrorCode = "UNSUPPORTED_USER_TYPE"

func NewUnsupportedUserTypeError(userType string) error {
	return domain.WrapBusinessError(
		UnsupportedUserTypeErr, unsupportedUserTypeErrorCode, fmt.Sprintf("Unsuported user type %s", userType), map[string]interface{}{},
	)
}

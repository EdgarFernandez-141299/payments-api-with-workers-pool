package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userPhoneNotFoundError = "USER_PHONE_NOT_FOUND"

func NewUserPhoneNotFoundError(memberId string, err error) error {
	msg := fmt.Sprintf("phone main not found in user: %s", memberId)
	return domain.WrapBusinessError(err, userPhoneNotFoundError, msg, nil)
}

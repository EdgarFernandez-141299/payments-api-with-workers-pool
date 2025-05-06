package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userNotFoundError = "USER_NOT_FOUND"

func NewUserNotFoundError(memberId string, err error) error {
	msg := fmt.Sprintf("user not found with memberId: %s", memberId)
	return domain.WrapBusinessError(err, userNotFoundError, msg, nil)
}

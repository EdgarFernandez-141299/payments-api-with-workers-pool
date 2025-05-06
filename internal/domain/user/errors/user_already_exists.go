package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userAlreadyExistError = "USER_ALREADY_EXISTS"

func NewUserAlreadyExistError(memberId string, err error) error {
	msg := fmt.Sprintf("user already exists with memberId: %s", memberId)
	return domain.WrapBusinessError(err, userAlreadyExistError, msg, nil)
}

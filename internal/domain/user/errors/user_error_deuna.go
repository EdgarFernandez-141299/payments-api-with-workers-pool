package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userFailCreateDeUnaError = "USER_FAIL_CREATE_DEUNA_PROVIDER"

func NewUserFailCreateDeUnaError(memberId string, err error) error {
	msg := fmt.Sprintf("error create user DEUNA provider with id: %s", memberId)
	return domain.WrapBusinessError(err, userFailCreateDeUnaError, msg, nil)
}

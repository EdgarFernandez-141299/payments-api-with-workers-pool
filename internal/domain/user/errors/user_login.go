package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const loginUserError = "LOGIN_USER_DEUNA_ERROR"

func NewLoginError(userId string, err error) error {
	msg := fmt.Sprintf("error login DEUNA with userId: %s", userId)
	return domain.WrapBusinessError(err, loginUserError, msg, nil)
}

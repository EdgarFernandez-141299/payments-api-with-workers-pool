package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userCreateError = "USER_EMAIL_NOT_FOUND"

func NewuserCreateError(err error) error {
	msg := "error creating user db record"
	return domain.WrapBusinessError(err, userCreateError, msg, nil)
}

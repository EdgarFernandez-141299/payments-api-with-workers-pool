package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userEmailNotFoundError = "USER_EMAIL_NOT_FOUND"

func NewUserEmailNotFoundError(memberId string, err error) error {
	msg := fmt.Sprintf("email main not found in user: %s", memberId)
	return domain.WrapBusinessError(err, userEmailNotFoundError, msg, nil)
}

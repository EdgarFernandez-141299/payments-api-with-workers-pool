package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userGettingDBError = "USER_GETTING_DB_ERROR"

func NewUserGettingDBError(memberId string, err error) error {
	msg := fmt.Sprintf("error getting user with memberId: %s", memberId)
	return domain.WrapBusinessError(err, userGettingDBError, msg, nil)
}

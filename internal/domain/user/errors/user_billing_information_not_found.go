package errors

import (
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const userBillingInformationNotFoundError = "USER_BILLLING_INFORMATION_NOT_FOUND"

func NewUserBillingInformationNotFoundError(memberId string, err error) error {
	msg := fmt.Sprintf("billing Information not found in user: %s", memberId)
	return domain.WrapBusinessError(err, userBillingInformationNotFoundError, msg, nil)
}

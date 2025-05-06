package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const InvalidUserIDErrorCode = "INVALID_USER_ID"

func NewInvalidUserIDError(id string) error {
	err := fmt.Errorf("invalid user id: %s", id)
	return domain.WrapBusinessError(err, InvalidUserIDErrorCode, err.Error(), map[string]interface{}{})
}

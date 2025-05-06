package errors

import (
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const orderAlreadyExistErrorCode = "ORDER_ALREADY_EXISTS"

func NewOrderAlreadyExistError(orderReferenceID string) error {
	err := fmt.Errorf("order %v already exists", orderReferenceID)
	return domain.WrapBusinessError(err, orderAlreadyExistErrorCode, err.Error(), map[string]interface{}{})
}

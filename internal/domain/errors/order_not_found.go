package errors

import (
	"fmt"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const OrderNotFoundErrorCode = "ORDER_NOT_FOUND"

func NewOrderNotFoundError(orderReferenceID string) error {
	err := fmt.Errorf("order %v not found", orderReferenceID)
	return domain.WrapBusinessError(err, OrderNotFoundErrorCode, err.Error(), nil)
}

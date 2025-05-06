package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const createCardError = "CREATE_CARD_ERROR"

func NewCreateCardError(err error) error {
	msg := "error creating card DEUNA service"
	return domain.WrapBusinessError(err, createCardError, msg, nil)
}

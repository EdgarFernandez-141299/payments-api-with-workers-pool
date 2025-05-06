package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const createCardDBError = "CREATE_CARD_DB_ERROR"

func NewCreateCardDBError(err error) error {
	msg := "error creating card db"
	return domain.WrapBusinessError(err, createCardDBError, msg, nil)
}

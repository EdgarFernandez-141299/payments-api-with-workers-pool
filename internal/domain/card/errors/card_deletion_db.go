package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const deleteCardDBError = "DELETE_CARD_DB_ERROR"

func NewDeleteCardDBError(err error) error {
	msg := "error deleting card db"
	return domain.WrapBusinessError(err, deleteCardDBError, msg, nil)
}

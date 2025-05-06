package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const CardErrorGetDB = "CARD_ERROR_GET_DB"

func NewCardErrorGetDB(err error) error {
	msg := "error getting card from database"
	return domain.WrapBusinessError(err, CardErrorGetDB, msg, nil)
}

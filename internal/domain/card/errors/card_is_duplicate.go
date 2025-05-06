package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const CardAlreadyExists = "CARD_ALREADY_EXISTS"

func NewCardIsAlreadyExists(err error) error {
	return domain.WrapBusinessError(err, CardAlreadyExists, "the card already exists", nil)
}

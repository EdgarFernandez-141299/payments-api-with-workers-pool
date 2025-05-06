package errors

import (
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const CardNotFound = "Card not found [ %s ]"
const CardNotFoundCode = "CARD_NOT_FOUND"

func NewCardNotFoundError(cardID string, err error) error {
	msg := fmt.Sprintf(CardNotFound, cardID)
	return domain.WrapBusinessError(err, CardNotFoundCode, msg, nil)
}

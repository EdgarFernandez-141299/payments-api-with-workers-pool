package errors

import (
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const cardIsNotFromMember = "Card is not from member %s "
const isNotCode = "CARD_IS_NOT_FROM_MEMBER"

func NewCardIsNotFromMember(memberID string, err error) error {
	msg := fmt.Sprintf(cardIsNotFromMember, memberID)
	return domain.WrapBusinessError(err, isNotCode, msg, nil)
}

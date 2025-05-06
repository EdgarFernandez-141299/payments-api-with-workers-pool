package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
	enums "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card"
)

const deleteCardError = "DELETE_CARD_ERROR"

func NewDeleteCardError(provider enums.Providers, err error) error {
	msg := "error deleting card " + string(provider) + " service"
	return domain.WrapBusinessError(err, deleteCardError, msg, nil)
}

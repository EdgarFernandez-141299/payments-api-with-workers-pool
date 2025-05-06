package errors

import (
	"fmt"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const collectionAccountNotFoundError = "COLLECTION_ACCOUNT_NOT_FOUND"

func NewCollectionAccountNotFound() error {
	return domain.WrapBusinessError(fmt.Errorf("collection account not found"), collectionAccountNotFoundError, "collection account not found", map[string]interface{}{})
}

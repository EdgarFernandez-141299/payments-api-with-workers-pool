package errors

import (
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
)

const (
	collectionCenterErrorNotFoundError = "COLLECTION_CENTER_NOT_FOUND"
	invalidCountryCodeError            = "INVALID_COUNTRY_CODE"
	invalidCurrencyCodeError           = "INVALID_CURRENCY_CODE"
	invalidAccountTypeError            = "INVALID_ACCOUNT_TYPE"
	invalidAssociatedOriginsError      = "INVALID_ASSOCIATED_ORIGINS"
	collectionAccountErrorDB           = "COLLECTION_ACCOUNT_ERROR_DB"
	currencyNoMatchError               = "CURRENCY_NO_MATCH"
)

func NewCollectionCenterNotFoundError(err error) error {
	msg := "collection center not found"
	return domain.WrapBusinessError(err, collectionCenterErrorNotFoundError, msg, err)
}

func NewInvalidCountryCodeError(err error) error {
	msg := "invalid country code"
	return domain.WrapBusinessError(err, invalidCountryCodeError, msg, err)
}

func NewInvalidCurrencyCodeError(err error) error {
	msg := "invalid currency code"
	return domain.WrapBusinessError(err, invalidCurrencyCodeError, msg, err)
}

func NewInvalidAccountTypeError(err error) error {
	msg := "invalid account type"
	return domain.WrapBusinessError(err, invalidAccountTypeError, msg, err)
}

func NewInvalidAssociatedOriginsError(err error) error {
	msg := "invalid associated origins"
	return domain.WrapBusinessError(err, invalidAssociatedOriginsError, msg, err)
}

func NewCollectionAccountError(err error) error {
	msg := "error creating collection account"
	return domain.WrapBusinessError(err, collectionAccountErrorDB, msg, err)
}

func NewCollectionAccountAlreadyExistsError(err error) error {
	msg := "collection account number already exists"
	return domain.WrapBusinessError(err, collectionAccountErrorDB, msg, err)
}

func NewCollectionAccountFindByError(err error) error {
	msg := "error finding collection account"
	return domain.WrapBusinessError(err, collectionAccountErrorDB, msg, err)
}

func NewCollectionAccountRouteAlreadyExist(err error) error {
	msg := "collection account route already exists"
	return domain.WrapBusinessError(err, collectionAccountErrorDB, msg, err)
}

func NewCollectionAccountNotFoundError(err error) error {
	msg := "collection account not found"
	return domain.WrapBusinessError(err, collectionAccountErrorDB, msg, err)
}

func NewCollectionAccountRouteCreateError(err error) error {
	msg := "error creating collection account route"
	return domain.WrapBusinessError(err, collectionAccountErrorDB, msg, err)
}

func NewCollectionAccountCurrencyCodeError(err error) error {
	msg := "currency code does not match"
	return domain.WrapBusinessError(err, currencyNoMatchError, msg, err)
}

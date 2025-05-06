package entities

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"

type CollectionAccount struct {
	ID                     string
	AccountType            string
	CollectionCenterID     string
	CurrencyCode           value_objects.CurrencyCode
	AccountNumber          string
	BankName               string
	InterbankAccountNumber string
	EnterpriseID           string
}

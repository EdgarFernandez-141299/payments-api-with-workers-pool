package response

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"

type CollectionAccountResponse struct {
	ID                 string `json:"id"`
	AccountType        string `json:"account_type"`
	CollectionCenterID string `json:"collection_center_id"`
	Currency           string `json:"currency_code"`
	AccountNumber      string `json:"account_number"`
	BankName           string `json:"bank_name"`
}

func NewCollectionAccountResponse(entity entities.CollectionAccountEntity) CollectionAccountResponse {
	return CollectionAccountResponse{
		ID:                 entity.ID.String(),
		AccountType:        entity.AccountType,
		CollectionCenterID: entity.CollectionCenterID,
		Currency:           entity.CurrencyCode,
		AccountNumber:      entity.AccountNumber,
		BankName:           entity.BankName,
	}
}

package request

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type CollectionAccountRequest struct {
	AccountType            enums.AccountType `json:"account_type" validate:"required"`
	CollectionCenterID     string            `json:"collection_center_id" validate:"required"`
	CurrencyCode           string            `json:"currency_code" validate:"required"`
	AccountNumber          string            `json:"account_number" validate:"required,numeric,min=8,max=20"`
	InterbankAccountNumber string            `json:"interbank_account_number"`
	BankName               string            `json:"bank_name"`
	BranchCode             string            `json:"branch_code"`
}

func (p *CollectionAccountRequest) Validate() error {
	err := validator.New().Struct(p)

	if err != nil {
		return err
	}

	return nil
}

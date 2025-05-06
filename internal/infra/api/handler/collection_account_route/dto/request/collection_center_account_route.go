package request

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
)

type CollectionAccountRouteRequest struct {
	CollectionAccountID string                 `json:"collection_account_id" validate:"required"`
	CountryCode         string                 `json:"country_code" validate:"required"`
	CurrencyCode        string                 `json:"currency_code" validate:"required"`
	AssociatedOrigin    enums.AssociatedOrigin `json:"associated_origin" validate:"required"`
}

func (c *CollectionAccountRouteRequest) Validate() error {
	err := validator.New().Struct(c)

	if err != nil {
		return err
	}

	return nil
}

package response

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"

type CollectionAccountRouteResponse struct {
	ID                  string `json:"id"`
	CollectionAccountID string `json:"collection_account_id"`
	CountryCode         string `json:"country_code"`
	CurrencyCode        string `json:"currency_code"`
}

func NewCollectionAccountRouteResponse(
	entity entities.CollectionAccountRouteEntity,
) CollectionAccountRouteResponse {
	return CollectionAccountRouteResponse{
		ID:                  entity.ID.String(),
		CollectionAccountID: entity.CollectionAccountID,
		CountryCode:         entity.CountryCode,
		CurrencyCode:        entity.CurrencyCode,
		
	}
}

type CollectionAccountRouteDisableResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	DisabledAt string `json:"disabled_at"`
}

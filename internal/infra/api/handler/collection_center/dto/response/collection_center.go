package response

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"

type CollectionCenterResponse struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	AvailableCurrencies []string `json:"available_currencies"`
	Description         string   `json:"description"`
}

func NewCollectionCenterResponse(entity entities.CollectionCenterEntity) CollectionCenterResponse {
	return CollectionCenterResponse{
		ID:                  entity.ID.String(),
		Name:                entity.Name,
		Description:         entity.Description,
		AvailableCurrencies: entity.AvailableCurrencies,
	}
}

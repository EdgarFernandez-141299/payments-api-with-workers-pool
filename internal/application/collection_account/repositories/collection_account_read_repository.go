package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"
)

type CollectionAccountReadRepositoryIF interface {
	GetCollectionAccountRoute(
		ctx context.Context, country, associatedOrigin, currency, enterpriseId string,
	) (entities.CollectionAccountEntity, error)
}

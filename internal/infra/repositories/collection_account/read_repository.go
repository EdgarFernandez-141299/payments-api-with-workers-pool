package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account/repositories"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"
	"gorm.io/gorm"
)

type CollectionAccountReadRepository struct {
	db *gorm.DB
}

func NewCollectionAccountReadRepository(db *gorm.DB) repositories.CollectionAccountReadRepositoryIF {
	return &CollectionAccountReadRepository{
		db: db,
	}
}

func (c *CollectionAccountReadRepository) GetCollectionAccountRoute(
	ctx context.Context, country, associatedOrigin, currency, enterpriseId string,
) (entities.CollectionAccountEntity, error) {
	var collectionAccount entities.CollectionAccountEntity

	err := c.db.WithContext(ctx).
		Select("collection_account.*").
		Where(`cr.country_code = ? 
			AND cr.associated_origin = ? 
			AND cr.currency_code = ? 
			AND collection_account.enterprise_id = ?`, country, associatedOrigin, currency, enterpriseId).
		Joins("JOIN collection_account_route cr ON cr.collection_account_id = collection_account.id").
		First(&collectionAccount).Error

	if err != nil {
		return entities.CollectionAccountEntity{}, err
	}

	return collectionAccount, nil
}

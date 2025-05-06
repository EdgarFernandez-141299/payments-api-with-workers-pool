package repositories

import (
	"context"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gorm.io/gorm"
)

type CollectionCenterAccountRouteRepositoryIF interface {
	Create(
		ctx context.Context,
		collectionRoute entities.CollectionAccountRouteEntity,
	) error
	Disable(
		ctx context.Context,
		id, enterpriseId string,
	) error
	FindRouteBy(ctx context.Context,
		countryCode, currencyCode, associatedOrigin, enterpriseID string,
	) (entities.CollectionAccountRouteEntity, error)
}

type CollectionCenterRepository struct {
	db *gorm.DB
}

func NewCollectionCenterRepositoryIF(db *gorm.DB) CollectionCenterAccountRouteRepositoryIF {
	return &CollectionCenterRepository{
		db: db,
	}
}

func (c *CollectionCenterRepository) Create(
	ctx context.Context,
	collectionRoute entities.CollectionAccountRouteEntity,
) error {
	return c.db.WithContext(ctx).Create(&collectionRoute).Error
}

func (c *CollectionCenterRepository) Disable(
	ctx context.Context,
	id, enterpriseId string,
) error {
	return c.db.WithContext(ctx).
		Model(&entities.CollectionAccountRouteEntity{}).
		Where("id = ? AND enterprise_id = ?", id, enterpriseId).
		Update("disabled_at", time.Now()).Error
}

func (c *CollectionCenterRepository) FindRouteBy(
	ctx context.Context, countryCode, currencyCode, associatedOrigin, enterpriseID string) (entities.CollectionAccountRouteEntity, error) {
	var route entities.CollectionAccountRouteEntity
	err := c.db.WithContext(ctx).
		Where("country_code = ? AND currency_code = ? AND associated_origin = ? AND enterprise_id = ?",
			countryCode, currencyCode, associatedOrigin, enterpriseID).
		First(&route).Error
	return route, err
}

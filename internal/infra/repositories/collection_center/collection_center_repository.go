package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gorm.io/gorm"
)

type CollectionCenterRepositoryIF interface {
	Create(ctx context.Context, collectionCenter entities.CollectionCenterEntity) error
	FindByID(
		ctx context.Context,
		id, enterpriseId string,
	) (entities.CollectionCenterEntity, error)
}

type CollectionCenterRepository struct {
	db *gorm.DB
}

func NewCollectionCenterRepository(db *gorm.DB) CollectionCenterRepositoryIF {
	return &CollectionCenterRepository{
		db: db,
	}
}

func (r *CollectionCenterRepository) Create(
	ctx context.Context,
	collectionCenter entities.CollectionCenterEntity,
) error {
	return r.db.WithContext(ctx).Create(&collectionCenter).Error
}

func (r *CollectionCenterRepository) FindByID(
	ctx context.Context,
	id, enterpriseId string,
) (entities.CollectionCenterEntity, error) {
	var entity entities.CollectionCenterEntity

	err := r.db.WithContext(ctx).
		Where("id = ? AND enterprise_id = ?", id, enterpriseId).
		First(&entity).Error

	return entity, err
}

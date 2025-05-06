package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"
	"gorm.io/gorm"
)

type CollectionAccountRepositoryIF interface {
	Create(ctx context.Context, entity entities.CollectionAccountEntity) error
	FindById(ctx context.Context, id, enterpriseId string) (entities.CollectionAccountEntity, error)
	FindByAccountNumber(
		ctx context.Context,
		accountNumber, enterpriseId string,
	) (entities.CollectionAccountEntity, error)
}

type CollectionAccountRepository struct {
	db *gorm.DB
}

func NewCollectionAccountRepository(db *gorm.DB) CollectionAccountRepositoryIF {
	return &CollectionAccountRepository{db: db}
}

func (r *CollectionAccountRepository) Create(
	ctx context.Context,
	entity entities.CollectionAccountEntity,
) error {
	return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *CollectionAccountRepository) FindById(
	ctx context.Context, id, enterpriseId string,
) (entities.CollectionAccountEntity, error) {
	var entity entities.CollectionAccountEntity

	err := r.db.WithContext(ctx).
		Where("id = ? AND enterprise_id = ?", id, enterpriseId).
		First(&entity).
		Error

	return entity, err
}

func (r *CollectionAccountRepository) FindByAccountNumber(
	ctx context.Context,
	accountNumber, enterpriseId string,
) (entities.CollectionAccountEntity, error) {
	var entity entities.CollectionAccountEntity

	err := r.db.WithContext(ctx).
		Where("account_number = ? AND enterprise_id = ?",
			accountNumber, enterpriseId).
		First(&entity).Error

	return entity, err
}

package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	"gorm.io/gorm"
)

var (
	ctx = context.Background()
)

type OrderWriteRepository struct {
	db *gorm.DB
}

func NewOrderWriteRepository(db *gorm.DB) repository.OrderWriteRepositoryIF {
	return &OrderWriteRepository{
		db: db,
	}
}

func (r *OrderWriteRepository) CreateOrder(
	ctx context.Context, entity entities.OrderEntity,
) error {
	return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *OrderWriteRepository) UpdateOrder(
	ctx context.Context, entity entities.OrderEntity,
) error {
	return r.db.WithContext(ctx).Model(&entity).Updates(entity).Error
}

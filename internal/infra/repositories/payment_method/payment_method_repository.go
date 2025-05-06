package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gorm.io/gorm"
)

type PaymentMethodRepositoryIF interface {
	Create(ctx context.Context, entity entities.PaymentMethodEntity) error
}

type PaymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepositoryIF {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) Create(ctx context.Context, entity entities.PaymentMethodEntity) error {
	return r.db.WithContext(ctx).Create(&entity).Error
}

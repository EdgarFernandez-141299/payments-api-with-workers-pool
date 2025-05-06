package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gorm.io/gorm"
)

type PaymentOrderWriteRespository struct {
	db *gorm.DB
}

func NewPaymentOrderWriteRespository(db *gorm.DB) repository.PaymentOrderRepositoryIF {
	return &PaymentOrderWriteRespository{
		db: db,
	}
}

func (r *PaymentOrderWriteRespository) CreatePaymentOrder(
	ctx context.Context, entity entities.PaymentOrderEntity,
) error {
	return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *PaymentOrderWriteRespository) UpdatePaymentOrder(
	ctx context.Context, entity entities.PaymentOrderEntity,
) error {
	return r.db.WithContext(ctx).Model(&entity).Updates(entity).Error
}

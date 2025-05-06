package refund

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
	"gorm.io/gorm"
)

type RefundWriteRepositoryIF struct {
	db *gorm.DB
}

func NewRefundWriteRepository(db *gorm.DB) repository.RefundWriteRepositoryIF {
	return &RefundWriteRepositoryIF{
		db: db,
	}
}

func (r *RefundWriteRepositoryIF) Create(ctx context.Context, refund entities.RefundEntity) error {
	if err := r.db.WithContext(ctx).Create(&refund).Error; err != nil {
		return err
	}

	return nil
}

func (r *RefundWriteRepositoryIF) UpdateRefund(
	ctx context.Context, entity entities.RefundEntity,
) error {
	return r.db.WithContext(ctx).Model(&entity).Updates(entity).Error
}

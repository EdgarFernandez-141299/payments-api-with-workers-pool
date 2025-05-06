package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gorm.io/gorm"
)

type PaymentConceptRepositoryIF interface {
	Create(ctx context.Context, entity entities.PaymentConceptEntity) error
	FindByCode(ctx context.Context, code, enterpriseId string) (entities.PaymentConceptEntity, error)
}

type PaymentConceptRepository struct {
	db *gorm.DB
}

func NewPaymentConceptRepository(db *gorm.DB) PaymentConceptRepositoryIF {
	return &PaymentConceptRepository{db: db}
}

func (r *PaymentConceptRepository) Create(
	ctx context.Context,
	entity entities.PaymentConceptEntity,
) error {
	return r.db.WithContext(ctx).Create(&entity).Error
}

func (r *PaymentConceptRepository) FindByCode(
	ctx context.Context, code, enterpriseId string,
) (entities.PaymentConceptEntity, error) {
	var entity entities.PaymentConceptEntity

	err := r.db.WithContext(ctx).
		Where("code = ? AND enterprise_id = ?", code, enterpriseId).
		First(&entity).
		Error

	return entity, err
}

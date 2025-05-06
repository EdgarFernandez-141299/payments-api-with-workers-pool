package repositories

import (
	"context"
	domainEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_receipt/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_receipt/entities"
	"gorm.io/gorm"
)

type PaymentReceiptRepository struct {
	db *gorm.DB
}

func NewPaymentReceiptRepository(db *gorm.DB) repository.PaymentReceiptRepository {
	return &PaymentReceiptRepository{
		db: db,
	}
}

func (r *PaymentReceiptRepository) GetByPaymentID(ctx context.Context, paymentID string) (domainEntities.PaymentReceipt, error) {
	var paymentReceiptDTO entities.PaymentReceiptDTO
	err := r.db.WithContext(ctx).
		Where("payment_id =?", paymentID).
		First(&paymentReceiptDTO).Error

	if err != nil {
		return domainEntities.PaymentReceipt{}, err
	}

	return paymentReceiptDTO.ToDomain(), nil
}

func (r *PaymentReceiptRepository) CreatePaymentReceipt(
	ctx context.Context, ent domainEntities.PaymentReceipt,
) error {
	entity := entities.FromEntity(ent)
	return r.db.WithContext(ctx).Create(&entity).Error
}

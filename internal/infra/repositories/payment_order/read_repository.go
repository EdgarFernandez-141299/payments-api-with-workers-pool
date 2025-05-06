package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gorm.io/gorm"
)

type PaymentOrderReadRepository struct {
	db *gorm.DB
}

func NewPaymentOrderReadRepository(db *gorm.DB) repository.GetPaymentOrderByReferenceIF {
	return &PaymentOrderReadRepository{
		db: db,
	}
}

func (r *PaymentOrderReadRepository) GetPaymentOrderByReference(
	ctx context.Context, referenceID, orderID, enterpriseID string,
) (entities.PaymentOrderEntity, error) {
	var paymentOrder entities.PaymentOrderEntity
	err := r.db.WithContext(ctx).
		Joins(`JOIN "order" ON "order".id = payment.order_id`).
		Where(`"order".reference_order_id = ? AND payment.payment_order_id = ? AND "order".enterprise_id = ?`, referenceID, orderID, enterpriseID).
		First(&paymentOrder).Error
	if err != nil {
		return entities.PaymentOrderEntity{}, err
	}
	return paymentOrder, nil
}

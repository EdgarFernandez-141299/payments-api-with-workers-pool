package refund

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/refund/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/refund/entities"
	"gorm.io/gorm"
)

type RefundReadRepository struct {
	db *gorm.DB
}

func NewRefundReadRepository(db *gorm.DB) repository.RefundReadRepositoryIF {
	return &RefundReadRepository{db: db}
}

func (r *RefundReadRepository) GetRefundsByReferenceOrderID(ctx context.Context, referenceOrderID, paymentID string, enterpriseID string) ([]entities.RefundEntity, error) {
	var refunds []entities.RefundEntity

	if err := r.db.
		WithContext(ctx).
		Joins(`JOIN payment ON payment.id = refund.payment_id`).
		Joins(`JOIN "order" ON "order".id = payment.order_id`).
		Where(`"order".reference_order_id = ? AND refund.payment_id = ? AND "order".enterprise_id = ?`, referenceOrderID, paymentID, enterpriseID).
		Find(&refunds).Error; err != nil {
		return []entities.RefundEntity{}, err
	}
	return refunds, nil
}

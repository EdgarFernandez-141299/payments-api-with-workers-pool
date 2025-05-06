package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	"gorm.io/gorm"
)

type OrderReadRepository struct {
	db *gorm.DB
}

func NewOrderReadRepository(db *gorm.DB) repository.OrderReadRepositoryIF {
	return &OrderReadRepository{
		db: db,
	}
}

func (r *OrderReadRepository) GetOrderByReferenceID(ctx context.Context, referenceOrderID, enterpriseID string) (entities.OrderEntity, error) {
	var order entities.OrderEntity

	err := r.db.WithContext(ctx).
		Where("reference_order_id = ?", referenceOrderID).
		First(&order).Error
	return order, err
}

func (r *OrderReadRepository) GetOrderPayments(ctx context.Context, referenceOrderID, enterpriseID string) ([]projections.OrderPaymentsProjection, error) {
	var orderPaymentProjection []projections.OrderPaymentsProjection

	err := r.db.WithContext(ctx).
		Model(&entities.OrderEntity{}).
		Select(`
			"order".reference_order_id, 
			"order".user_id, 
			"order".total_amount, 
			"order".currency_code, 
			"order".country_code, 
			"order".status as order_status,
			"order".metadata as metadata, 
			p.id as payment_id,
			p.card_id, 
			p.payment_method, 
			p.status as payment_status,
			p.authorization_code, 
			p.payment_order_id, 
			p.created_at as payment_date
		`).
		Joins(`LEFT JOIN payment p on "order".id = p.order_id`).
		Where(`"order".reference_order_id = ? AND "order".enterprise_id = ? AND "order".deleted_at IS NULL`, referenceOrderID, enterpriseID).
		Scan(&orderPaymentProjection).Error

	return orderPaymentProjection, err
}

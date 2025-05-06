package repositories

import (
	"context"

	orderRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	paymentOrderRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gorm.io/gorm"
)

type UpdateOrderStatusRepository struct {
	db                         *gorm.DB
	paymentOrderReadRepository paymentOrderRepository.GetPaymentOrderByReferenceIF
	orderReadRepository        orderRepository.OrderReadRepositoryIF
}

func NewUpdateOrderStatusRepository(db *gorm.DB,
	paymentOrderReadRepository paymentOrderRepository.GetPaymentOrderByReferenceIF,
	orderReadRepository orderRepository.OrderReadRepositoryIF) repository.TransactionsRepositoryIF {
	return &UpdateOrderStatusRepository{
		db:                         db,
		paymentOrderReadRepository: paymentOrderReadRepository,
		orderReadRepository:        orderReadRepository,
	}
}

func (r *UpdateOrderStatusRepository) UpdatePaymentOrderStatus(ctx context.Context, orderReferenceId, paymentId, enterpriseID string, status enums.PaymentStatus) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	paymentOrder, err := r.paymentOrderReadRepository.GetPaymentOrderByReference(ctx, orderReferenceId, paymentId, enterpriseID)
	if err != nil {
		tx.Rollback()
		return err
	}

	paymentOrder.SetStatus(status)
	if err := tx.Save(&paymentOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	order, err := r.orderReadRepository.GetOrderByReferenceID(ctx, orderReferenceId, enterpriseID)
	if err != nil {
		tx.Rollback()
		return err
	}

	order.SetStatus(status.String())
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

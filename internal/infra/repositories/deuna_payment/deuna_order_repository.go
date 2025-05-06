package repository

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/deuna_payment/entities"
	"gorm.io/gorm"
)

var OrderTokeNotFoundError = errors.New("order token not found")

type DeunaPaymentWriteRepository struct {
	db *gorm.DB
}

func NewDeunaPaymentWriteRepository(db *gorm.DB) repository.DeunaOrderRepository {
	return &DeunaPaymentWriteRepository{
		db: db,
	}
}

func (r *DeunaPaymentWriteRepository) CreatePaymentOrderDeuna(ctx context.Context, paymentID, orderID, deunaOrderToken string) error {
	deunaOrderID := utils.NewDeunaOrderID(orderID, paymentID)
	paymentOrderID := entities.NewDeunaPaymentEntity(deunaOrderID.GetID(), deunaOrderToken)

	return r.db.WithContext(ctx).Create(&paymentOrderID).Error
}

func (r *DeunaPaymentWriteRepository) GetTokenByOrderAndPaymentID(ctx context.Context, orderID string, paymentID string) (token string, err error) {
	deunaOrderID := utils.NewDeunaOrderID(orderID, paymentID)
	var paymentEntity entities.DeunaPaymentEntity
	err = r.db.WithContext(ctx).
		Where("payment_id = ?", deunaOrderID.GetID()).
		First(&paymentEntity).Error

	if err != nil {
		return "", OrderTokeNotFoundError
	}

	if paymentEntity.IsEmptyToken() {
		return "", nil
	}

	return paymentEntity.OrderToken, nil
}

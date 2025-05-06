package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	orderEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	paymentOrderEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTransactionsTestDB(t *testing.T) *gorm.DB {
	cxn := ":memory:?cache=shared"
	gormDB, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{})

	if err != nil {
		t.Errorf("failed to open stub database: %v", err)
	}

	if migrateErr := gormDB.AutoMigrate(
		&paymentOrderEntities.PaymentOrderEntity{},
		&orderEntities.OrderEntity{},
	); migrateErr != nil {
		t.Errorf("failed to migrate database: %v", migrateErr)
	}

	return gormDB
}

func TestUpdatePaymentOrderStatus_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockDB := setupTransactionsTestDB(t)
	mockPaymentOrderRepo := repository.NewGetPaymentOrderByReferenceIF(t)
	mockOrderRepo := repository.NewOrderReadRepositoryIF(t)

	// Mock payment order
	paymentOrder := paymentOrderEntities.PaymentOrderEntity{}
	mockPaymentOrderRepo.On("GetPaymentOrderByReference", ctx, "order123", "payment123", "enterprise123").Return(paymentOrder, nil)

	// Mock order
	order := orderEntities.OrderEntity{}
	mockOrderRepo.On("GetOrderByReferenceID", ctx, "order123", "enterprise123").Return(order, nil)

	// Create repository
	repo := NewUpdateOrderStatusRepository(mockDB, mockPaymentOrderRepo, mockOrderRepo)

	// Execute
	err := repo.UpdatePaymentOrderStatus(ctx, "order123", "payment123", "enterprise123", enums.PaymentProcessed)

	// Assert
	assert.NoError(t, err)
	mockPaymentOrderRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

func TestUpdatePaymentOrderStatus_ErrorGettingPaymentOrder(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockDB := setupTransactionsTestDB(t)
	mockPaymentOrderRepo := repository.NewGetPaymentOrderByReferenceIF(t)
	mockOrderRepo := repository.NewOrderReadRepositoryIF(t)

	// Mock payment order error
	mockPaymentOrderRepo.On("GetPaymentOrderByReference", ctx, "order123", "payment123", "enterprise123").Return(paymentOrderEntities.PaymentOrderEntity{}, errors.New("payment order not found"))

	// Create repository
	repo := NewUpdateOrderStatusRepository(mockDB, mockPaymentOrderRepo, mockOrderRepo)

	// Execute
	err := repo.UpdatePaymentOrderStatus(ctx, "order123", "payment123", "enterprise123", enums.PaymentProcessed)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "payment order not found", err.Error())
	mockPaymentOrderRepo.AssertExpectations(t)
}

func TestUpdatePaymentOrderStatus_ErrorGettingOrder(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockDB := setupTransactionsTestDB(t)
	mockPaymentOrderRepo := repository.NewGetPaymentOrderByReferenceIF(t)
	mockOrderRepo := repository.NewOrderReadRepositoryIF(t)

	// Mock payment order
	paymentOrder := paymentOrderEntities.PaymentOrderEntity{}
	mockPaymentOrderRepo.On("GetPaymentOrderByReference", ctx, "order123", "payment123", "enterprise123").Return(paymentOrder, nil)

	// Mock order error
	mockOrderRepo.On("GetOrderByReferenceID", ctx, "order123", "enterprise123").Return(orderEntities.OrderEntity{}, errors.New("order not found"))

	// Create repository
	repo := NewUpdateOrderStatusRepository(mockDB, mockPaymentOrderRepo, mockOrderRepo)

	// Execute
	err := repo.UpdatePaymentOrderStatus(ctx, "order123", "payment123", "enterprise123", enums.PaymentProcessed)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "order not found", err.Error())
	mockPaymentOrderRepo.AssertExpectations(t)
	mockOrderRepo.AssertExpectations(t)
}

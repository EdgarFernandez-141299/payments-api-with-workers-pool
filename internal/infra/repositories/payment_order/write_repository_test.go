package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&entities.PaymentOrderEntity{})
	return db, nil
}

var (
	id, _ = uid.NewUniqueID(uid.WithID("id"))
)

func TestPaymentOrderRepository_CreatePaymentOrder(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	repo := NewPaymentOrderWriteRespository(db)

	ctx := context.Background()
	paymentOrder := entities.PaymentOrderEntity{
		ID:     uid.GenerateID(),
		Status: enums.Active.String(),
	}

	err = repo.CreatePaymentOrder(ctx, paymentOrder)
	assert.NoError(t, err)

	var checkPayment entities.PaymentOrderEntity
	result := db.First(&checkPayment, "id = ?", paymentOrder.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, paymentOrder.ID, checkPayment.ID)
	assert.Equal(t, paymentOrder.Status, checkPayment.Status)
}

func TestPaymentOrderRepository_UpdatePaymentOrder(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("failed to setup test db: %v", err)
	}
	repo := NewPaymentOrderWriteRespository(db)

	ctx := context.Background()

	paymentOrder := entities.PaymentOrderEntity{
		ID:     uid.GenerateID(),
		Status: enums.Active.String(),
	}
	err = repo.CreatePaymentOrder(ctx, paymentOrder)
	assert.NoError(t, err)

	paymentOrder.Status = enums.PaymentProcessed.String()
	err = repo.UpdatePaymentOrder(ctx, paymentOrder)
	assert.NoError(t, err)

	var checkPayment entities.PaymentOrderEntity
	result := db.First(&checkPayment, "id = ?", paymentOrder.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, paymentOrder.ID, checkPayment.ID)
	assert.Equal(t, enums.PaymentProcessed.String(), checkPayment.Status)
}

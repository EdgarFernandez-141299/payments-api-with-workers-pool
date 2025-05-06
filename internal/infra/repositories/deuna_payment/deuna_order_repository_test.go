package repository

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/deuna_payment/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeunaPaymentWriteRepository_CreatePaymentOrderDeuna(t *testing.T) {
	db := setupTestDB()
	repo := &DeunaPaymentWriteRepository{db: db}

	type args struct {
		paymentID       string
		orderID         string
		deunaOrderToken string
	}
	tests := []struct {
		name          string
		setup         func()
		args          args
		expectedError error
	}{
		{
			name: "successful creation",
			setup: func() {
				// Prepare database without conflicting data
			},
			args: args{
				paymentID:       "payment123",
				orderID:         "order123",
				deunaOrderToken: "validtoken",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.Migrator().DropTable(&entities.DeunaPaymentEntity{})
			_ = db.AutoMigrate(&entities.DeunaPaymentEntity{})
			tt.setup()

			err := repo.CreatePaymentOrderDeuna(context.Background(), tt.args.paymentID, tt.args.orderID, tt.args.deunaOrderToken)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&entities.DeunaPaymentEntity{})
	return db
}

func TestDeunaPaymentWriteRepository_GetTokenByOrderAndPaymentID(t *testing.T) {
	db := setupTestDB()
	repo := &DeunaPaymentWriteRepository{db: db}

	type args struct {
		orderID   string
		paymentID string
	}
	tests := []struct {
		name          string
		setup         func()
		args          args
		expectedToken string
		expectedError error
	}{
		{
			name: "valid token found",
			setup: func() {
				orderID := "order123"
				paymentID := "payment123"
				token := "validtoken"
				entity := entities.NewDeunaPaymentEntity(utils.NewDeunaOrderID(orderID, paymentID).GetID(), token)
				_ = db.Create(&entity).Error
			},
			args: args{
				orderID:   "order123",
				paymentID: "payment123",
			},
			expectedToken: "validtoken",
			expectedError: nil,
		},
		{
			name: "token not found",
			setup: func() {
				// No entry created
			},
			args: args{
				orderID:   "notfound",
				paymentID: "notfound",
			},
			expectedToken: "",
			expectedError: OrderTokeNotFoundError,
		},
		{
			name: "empty token",
			setup: func() {
				orderID := "order456"
				paymentID := "payment456"
				entity := entities.NewDeunaPaymentEntity(utils.NewDeunaOrderID(orderID, paymentID).GetID(), "")
				_ = db.Create(&entity).Error
			},
			args: args{
				orderID:   "order456",
				paymentID: "payment456",
			},
			expectedToken: "",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db.Migrator().DropTable(&entities.DeunaPaymentEntity{})
			_ = db.AutoMigrate(&entities.DeunaPaymentEntity{})
			tt.setup()

			token, err := repo.GetTokenByOrderAndPaymentID(context.Background(), tt.args.orderID, tt.args.paymentID)

			assert.Equal(t, tt.expectedToken, token)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

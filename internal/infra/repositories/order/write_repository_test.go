package repositories

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/gommon/uid"
	orderentities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	paymententities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/payment_order/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ()

func setupTestDB(t *testing.T, ent *orderentities.OrderEntity) *gorm.DB {
	cxn := ":memory:?cache=shared"
	gormDB, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{})

	if err != nil {
		t.Errorf("failed to open stub database: %v", err)
	}

	if migrateErr := gormDB.AutoMigrate(&orderentities.OrderEntity{}, &paymententities.PaymentOrderEntity{}); migrateErr != nil {
		t.Errorf("failed to migrate database: %v", migrateErr)
	}

	if ent != nil {
		if err := gormDB.Create(&ent).Error; err != nil {
			t.Errorf("failed to create entity: %v", err)
		}
	}

	return gormDB
}

func TestWriteRepository_CreateOrder(t *testing.T) {
	id, _ := uid.NewUniqueID()

	validOrder := orderentities.OrderEntity{
		ID:               id.String(),
		ReferenceOrderID: "test-reference-order-id",
		UserID:           "test-user-id",
		TotalAmount:      decimal.NewFromFloat(100),
		CountryCode:      "US",
		CurrencyCode:     "USD",
	}

	tests := []struct {
		name    string
		Order   orderentities.OrderEntity
		wantErr error
		setup   func(t *testing.T, ent *orderentities.OrderEntity) *gorm.DB
	}{
		{
			name:    "Create order success",
			Order:   validOrder,
			wantErr: nil,
			setup:   setupTestDB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup(t, nil)

			repo := NewOrderWriteRepository(db)

			err := repo.CreateOrder(ctx, tt.Order)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error: got %v, want nil", err)
				}
			}
		})
	}
}

func TestWriteRepository_UpdateOrder(t *testing.T) {
	id, _ := uid.NewUniqueID()

	initialOrder := orderentities.OrderEntity{
		ID:               id.String(),
		ReferenceOrderID: "test-reference-order-id",
		UserID:           "test-user-id",
		TotalAmount:      decimal.NewFromFloat(100),
		CountryCode:      "US",
		CurrencyCode:     "USD",
		Status:           "pending",
	}

	updatedOrder := initialOrder
	updatedOrder.Status = "completed"

	tests := []struct {
		name    string
		initial orderentities.OrderEntity
		update  orderentities.OrderEntity
		wantErr error
		setup   func(t *testing.T, ent *orderentities.OrderEntity) *gorm.DB
	}{
		{
			name:    "Update order success",
			initial: initialOrder,
			update:  updatedOrder,
			wantErr: nil,
			setup:   setupTestDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup(t, &tt.initial)

			repo := NewOrderWriteRepository(db)

			err := repo.UpdateOrder(ctx, tt.update)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error: got %v, want nil", err)
				}

				// Verificar que la actualización se realizó correctamente
				var result orderentities.OrderEntity
				if err := db.First(&result, "id = ?", tt.update.ID).Error; err != nil {
					t.Errorf("failed to fetch updated order: %v", err)
				}

				if result.Status != tt.update.Status {
					t.Errorf("status not updated: got %v, want %v", result.Status, tt.update.Status)
				}
			}
		})
	}
}

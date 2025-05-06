package repositories

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/order/entities"
	"gorm.io/gorm"
)

func TestReadRepository_GetOrderByReferenceID(t *testing.T) {
	id := "34"
	validOrder := entities.OrderEntity{
		ID:               id,
		ReferenceOrderID: "referenceOrderID",
		UserID:           "userID",
		TotalAmount:      decimal.NewFromFloat(100),
		CountryCode:      "US",
		CurrencyCode:     "USD",
	}

	tests := []struct {
		name             string
		order            entities.OrderEntity
		referenceOrderID string
		enterpriseID     string
		wantErr          error
		setup            func(t *testing.T, ent *entities.OrderEntity) *gorm.DB
	}{
		{
			name:             "get order by reference id success",
			order:            validOrder,
			wantErr:          nil,
			setup:            setupTestDB,
			referenceOrderID: "referenceOrderID",
			enterpriseID:     "enterpriseID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup(t, &validOrder)

			repo := NewOrderReadRepository(db)

			_, err := repo.GetOrderByReferenceID(ctx, tt.referenceOrderID, tt.enterpriseID)

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

func TestReadRepository_GetOrderPayments(t *testing.T) {
	id := "35"
	validOrder := entities.OrderEntity{
		ID:               id,
		ReferenceOrderID: "referenceOrderID",
		UserID:           "userID",
		TotalAmount:      decimal.NewFromFloat(100),
		CountryCode:      "US",
		CurrencyCode:     "USD",
	}

	tests := []struct {
		name             string
		order            entities.OrderEntity
		referenceOrderID string
		enterpriseID     string
		wantErr          error
		setup            func(t *testing.T, ent *entities.OrderEntity) *gorm.DB
	}{
		{
			name:             "get order status success",
			order:            validOrder,
			wantErr:          nil,
			setup:            setupTestDB,
			referenceOrderID: "referenceOrderID",
			enterpriseID:     "enterpriseID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup(t, &validOrder)

			repo := NewOrderReadRepository(db)

			_, err := repo.GetOrderPayments(ctx, tt.referenceOrderID, tt.enterpriseID)

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

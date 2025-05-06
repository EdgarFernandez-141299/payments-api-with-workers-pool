package repositories

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/gommon/uid"
	entitiesDomain "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T, ent *entities.CollectionAccountEntity, entR *entitiesDomain.CollectionAccountRouteEntity) *gorm.DB {
	cxn := ":memory:?cache=shared"
	gormDB, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{})

	if err != nil {
		t.Errorf("failed to open stub database: %v", err)
	}

	if migrateErr := gormDB.AutoMigrate(
		&entities.CollectionAccountEntity{},
		&entitiesDomain.CollectionAccountRouteEntity{},
	); migrateErr != nil {
		t.Errorf("failed to migrate database: %v", migrateErr)
	}

	if ent != nil {
		if err := gormDB.Create(&ent).Error; err != nil {
			t.Errorf("failed to create entity: %v", err)
		}
	}

	if entR != nil {
		if err := gormDB.Create(&entR).Error; err != nil {
			t.Errorf("failed to create entity: %v", err)
		}
	}

	return gormDB
}

func TestCollectionAccountFindById(t *testing.T) {
	enterpriseID := "test-enterprise"
	validID := uid.GenerateID()

	validEntity := entities.CollectionAccountEntity{
		ID:                 validID,
		AccountNumber:      "123456",
		EnterpriseID:       enterpriseID,
		AccountType:        "mix",
		BankName:           "Test Bank",
		CollectionCenterID: "test-collection-center",
		CurrencyCode:       "USD",
	}

	validEntityRoute := entitiesDomain.CollectionAccountRouteEntity{
		ID:                  validID,
		CollectionAccountID: validID.String(),
		CountryCode:         "MX",
		CurrencyCode:        "MXN",
		EnterpriseID:        enterpriseID,
		AssociatedOrigin:    "DOWNPAYMENT",
	}

	tests := []struct {
		name            string
		id              string
		enterpriseID    string
		wantErr         error
		wantEntity      *entities.CollectionAccountEntity
		wantEntityRoute *entitiesDomain.CollectionAccountRouteEntity
		setup           func(t *testing.T, ent *entities.CollectionAccountEntity,
			entR *entitiesDomain.CollectionAccountRouteEntity) *gorm.DB
	}{
		{
			name:            "Found",
			id:              validID.String(),
			enterpriseID:    enterpriseID,
			wantErr:         nil,
			wantEntity:      &validEntity,
			wantEntityRoute: &validEntityRoute,
			setup:           setupTestDB,
		},
		{
			name:         "Not Found",
			id:           "wrong-id",
			enterpriseID: enterpriseID,
			wantErr:      gorm.ErrRecordNotFound,
			wantEntity:   nil,
			setup:        setupTestDB,
		},
		{
			name:         "Wrong Enterprise ID",
			id:           validID.String(),
			enterpriseID: "wrong-enterprise",
			wantErr:      gorm.ErrRecordNotFound,
			wantEntity:   nil,
			setup:        setupTestDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db := tc.setup(t, tc.wantEntity, tc.wantEntityRoute)

			repo := NewCollectionAccountRepository(db)

			gotEntity, gotErr := repo.FindById(ctx, tc.id, tc.enterpriseID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tc.wantErr)
			}

			if tc.wantEntity != nil && gotEntity.ID.String() != tc.wantEntity.ID.String() {
				t.Errorf("unexpected entity: got %v, want %v", gotEntity, tc.wantEntity)
			}
		})
	}
}

func TestCollectionAccountFindByAccountNumber(t *testing.T) {
	enterpriseID := "test-enterprise"
	validID := uid.GenerateID()

	validEntity := entities.CollectionAccountEntity{
		ID:                 validID,
		AccountNumber:      "123456",
		EnterpriseID:       enterpriseID,
		AccountType:        "mix",
		BankName:           "Test Bank",
		CollectionCenterID: "test-collection-center",
		CurrencyCode:       "USD",
	}

	tests := []struct {
		name          string
		accountNumber string
		enterpriseID  string
		wantErr       error
		wantEntity    *entities.CollectionAccountEntity
		setup         func(t *testing.T, ent *entities.CollectionAccountEntity,
			entR *entitiesDomain.CollectionAccountRouteEntity) *gorm.DB
	}{
		{
			name:          "Found",
			accountNumber: validEntity.AccountNumber,
			enterpriseID:  enterpriseID,
			wantErr:       nil,
			wantEntity:    &validEntity,
			setup:         setupTestDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db := tc.setup(t, tc.wantEntity, nil)

			repo := NewCollectionAccountRepository(db)

			gotEntity, gotErr := repo.FindByAccountNumber(ctx, tc.accountNumber, tc.enterpriseID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tc.wantErr)
			}

			if tc.wantEntity != nil && gotEntity.ID.String() != tc.wantEntity.ID.String() {
				t.Errorf("unexpected entity: got %v, want %v", gotEntity, tc.wantEntity)
			}
		})
	}
}

func TestAccountCreate(t *testing.T) {
	tests := []struct {
		name    string
		entity  *entities.CollectionAccountEntity
		wantErr error
		setup   func(t *testing.T, ent *entities.CollectionAccountEntity,
			entR *entitiesDomain.CollectionAccountRouteEntity) *gorm.DB
	}{
		{
			name: "Create Success",
			entity: &entities.CollectionAccountEntity{
				ID:                 uid.GenerateID(),
				AccountNumber:      "654321",
				EnterpriseID:       "test-enterprise",
				AccountType:        "mix",
				BankName:           "Test Bank",
				CollectionCenterID: "test-collection-center",
				CurrencyCode:       "USD",
			},
			wantErr: nil,
			setup:   setupTestDB,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db := tc.setup(t, nil, nil)

			repo := NewCollectionAccountRepository(db)

			gotErr := repo.Create(ctx, *tc.entity)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tc.wantErr)
			}

			if tc.wantErr == nil {
				var gotEntity entities.CollectionAccountEntity
				if err := db.First(&gotEntity, "id = ?", tc.entity.ID).Error; err != nil {
					t.Errorf("failed to find created entity: %v", err)
				}

				if gotEntity.ID.String() != tc.entity.ID.String() {
					t.Errorf("unexpected entity: got %v, want %v", gotEntity, tc.entity)
				}
			}
		})
	}
}

package repositories

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/gommon/uid"
	entitiesDomain "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"

	"gorm.io/gorm"
)

func TestCollectionAccountReadRepository(t *testing.T) {
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
		name              string
		id                string
		enterpriseID      string
		country           string
		associatedOrigin  string
		currency          string
		wantErr           error
		wantEntityAccount *entities.CollectionAccountEntity
		wantEntityRoute   *entitiesDomain.CollectionAccountRouteEntity

		setup func(t *testing.T, ent *entities.CollectionAccountEntity, entR *entitiesDomain.CollectionAccountRouteEntity) *gorm.DB
	}{
		{
			name:              "Found",
			id:                validID.String(),
			enterpriseID:      enterpriseID,
			wantErr:           nil,
			wantEntityAccount: &validEntity,
			wantEntityRoute:   &validEntityRoute,
			setup:             setupTestDB,
			country:           "MX",
			associatedOrigin:  "DOWNPAYMENT",
			currency:          "MXN",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db := tc.setup(t, tc.wantEntityAccount, tc.wantEntityRoute)

			repo := NewCollectionAccountReadRepository(db)

			gotEntity, gotErr := repo.GetCollectionAccountRoute(ctx, tc.country, tc.associatedOrigin, tc.currency, tc.enterpriseID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tc.wantErr)
			}

			if tc.wantEntityAccount != nil && gotEntity.ID.String() != tc.wantEntityAccount.ID.String() {
				t.Errorf("unexpected entity: got %v, want %v", gotEntity, tc.wantEntityAccount)

			}
		})
	}
}

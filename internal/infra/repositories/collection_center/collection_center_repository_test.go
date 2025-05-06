package repositories

import (
	"context"
	"errors"
	"testing"

	"gorm.io/driver/sqlite"

	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T, ent *entities.CollectionCenterEntity) *gorm.DB {
	cxn := ":memory:?cache=shared"
	gormDB, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{})

	if err != nil {
		t.Errorf("failed to open stub database: %v", err)
	}

	if migrateErr := gormDB.AutoMigrate(&entities.CollectionCenterEntity{}); migrateErr != nil {
		t.Errorf("failed to migrate database: %v", migrateErr)
	}

	if ent != nil {
		if err := gormDB.Create(&ent).Error; err != nil {
			t.Errorf("failed to create entity: %v", err)
		}
	}

	return gormDB
}

func TestCollectionCenterRepository_FindByID(t *testing.T) {
	enterpriseID := "test-enterprise"
	validID := uid.GenerateID()

	validEntity := entities.CollectionCenterEntity{
		ID:           validID,
		Name:         "Test Collection Center",
		Description:  "Description",
		EnterpriseID: enterpriseID,
	}

	tests := []struct {
		name         string
		id           string
		enterpriseID string
		wantErr      error
		wantEntity   *entities.CollectionCenterEntity
		setup        func(t *testing.T, ent *entities.CollectionCenterEntity) *gorm.DB
	}{
		{
			name:         "Found",
			id:           validID.String(),
			enterpriseID: enterpriseID,
			wantErr:      nil,
			wantEntity:   &validEntity,
			setup:        setupTestDB,
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

			db := tc.setup(t, tc.wantEntity)

			repo := NewCollectionCenterRepository(db)
			gotEntity, gotErr := repo.FindByID(ctx, tc.id, tc.enterpriseID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tc.wantErr)
			}

			if tc.wantEntity != nil && gotEntity.ID.String() != tc.wantEntity.ID.String() {
				t.Errorf("unexpected entity: got %v, want %v", gotEntity, tc.wantEntity)
			}
		})
	}
}

package repositories

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ctx = context.TODO()
)

func setupTestDB(t *testing.T, ent *entities.UserEntity) *gorm.DB {
	cxn := ":memory:?cache=shared"
	gormDB, err := gorm.Open(sqlite.Open(cxn), &gorm.Config{})

	if err != nil {
		t.Errorf("failed to open stub database: %v", err)
	}

	if migrateErr := gormDB.AutoMigrate(&entities.UserEntity{}); migrateErr != nil {
		t.Errorf("failed to migrate database: %v", migrateErr)
	}

	if ent != nil {
		if err := gormDB.Create(&ent).Error; err != nil {
			t.Errorf("failed to create entity: %v", err)
		}
	}

	return gormDB
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	validUser := entities.UserEntity{
		ID:             "user-id",
		ExternalUserID: "test-member-id",
		Email:          "test-email",
		EnterpriseID:   "test-enterprise-id",
	}

	tests := []struct {
		name       string
		email      string
		enterprise string
		wantErr    error
		setup      func(t *testing.T, ent *entities.UserEntity) *gorm.DB
	}{
		{
			name:       "GetTokenByOrderAndPaymentID Success",
			email:      validUser.Email,
			enterprise: validUser.EnterpriseID,
			wantErr:    nil,
			setup:      setupTestDB,
		},
		{
			name:       "GetTokenByOrderAndPaymentID User Empty",
			email:      "non-existent-email",
			enterprise: validUser.EnterpriseID,
			wantErr:    nil,
			setup:      setupTestDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup(t, &validUser)

			repo := NewUserReadRepository(db)

			user, gotErr := repo.GetUserByEmail(ctx, tt.email, tt.enterprise)

			if tt.name == "GetTokenByOrderAndPaymentID Success" {
				if gotErr != nil {
					t.Errorf("unexpected error: got %v, want nil", gotErr)
				}

				if user.ID != validUser.ID {
					t.Errorf("unexpected ID: got %v, want %v", user.ID, validUser.ID)
				}
				if user.Email != validUser.Email {
					t.Errorf("unexpected Email: got %v, want %v", user.Email, validUser.Email)
				}
				if user.ExternalUserID != validUser.ExternalUserID {
					t.Errorf("unexpected ExternalUserID: got %v, want %v", user.ExternalUserID, validUser.ExternalUserID)
				}
				if user.EnterpriseID != validUser.EnterpriseID {
					t.Errorf("unexpected EnterpriseID: got %v, want %v", user.EnterpriseID, validUser.EnterpriseID)
				}
			} else if tt.name == "GetTokenByOrderAndPaymentID User Empty" {
				if gotErr != nil {
					t.Errorf("unexpected error: got %v, want nil", gotErr)
				}

				if user.ID != "" || user.Email != "" || user.ExternalUserID != "" || user.EnterpriseID != "" {
					t.Errorf("expected empty entity, got: %+v", user)
				}
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	validUser := entities.UserEntity{
		ID:             "user-id",
		ExternalUserID: "test-member-id",
		Email:          "test-email",
		EnterpriseID:   "test-enterprise-id",
	}

	tests := []struct {
		name       string
		memberID   string
		enterprise string
		wantErr    error
		setup      func(t *testing.T, ent *entities.UserEntity) *gorm.DB
	}{
		{
			name:       "GetTokenByOrderAndPaymentID Success",
			memberID:   validUser.ID,
			enterprise: validUser.EnterpriseID,
			wantErr:    nil,
			setup:      setupTestDB,
		},
		{
			name:       "GetTokenByOrderAndPaymentID User Empty",
			memberID:   "non-existent-member-id",
			enterprise: validUser.EnterpriseID,
			wantErr:    nil,
			setup:      setupTestDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup(t, &validUser)

			repo := NewUserReadRepository(db)

			user, gotErr := repo.GetUserByID(ctx, tt.memberID, tt.enterprise)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tt.wantErr)
			}

			if tt.name == "GetTokenByOrderAndPaymentID Success" {
				if gotErr != nil {
					t.Errorf("unexpected error: got %v, want nil", gotErr)
				}

				if user.ID != validUser.ID {
					t.Errorf("unexpected ID: got %v, want %v", user.ID, validUser.ID)
				}
				if user.Email != validUser.Email {
					t.Errorf("unexpected Email: got %v, want %v", user.Email, validUser.Email)
				}
				if user.ExternalUserID != validUser.ExternalUserID {
					t.Errorf("unexpected ExternalUserID: got %v, want %v", user.ExternalUserID, validUser.ExternalUserID)
				}
				if user.EnterpriseID != validUser.EnterpriseID {
					t.Errorf("unexpected EnterpriseID: got %v, want %v", user.EnterpriseID, validUser.EnterpriseID)
				}
			} else if tt.name == "GetTokenByOrderAndPaymentID User Empty" {
				if gotErr != nil {
					t.Errorf("unexpected error: got %v, want nil", gotErr)
				}

				if user.ID != "" || user.Email != "" || user.ExternalUserID != "" || user.EnterpriseID != "" {
					t.Errorf("expected empty entity, got: %+v", user)
				}
			}
		})
	}
}

func TestUserRepository_GetEmailByUserID(t *testing.T) {
	validUser := entities.UserEntity{
		ID:             "user-id",
		ExternalUserID: "test-member-id",
		Email:          "test-email",
		EnterpriseID:   "test-enterprise-id",
	}

	tests := []struct {
		name       string
		userID     string
		enterprise string
		wantErr    error
		setup      func(t *testing.T, ent *entities.UserEntity) *gorm.DB
	}{
		{
			name:       "GetEmailByUserID Success",
			userID:     validUser.ID,
			enterprise: validUser.EnterpriseID,
			wantErr:    nil,
			setup:      setupTestDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := tt.setup(t, &validUser)

			repo := NewUserReadRepository(db)

			email, gotErr := repo.GetEmailByUserID(ctx, tt.userID, tt.enterprise)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tt.wantErr)
			}

			if tt.name == "GetEmailByUserID Success" {
				if gotErr != nil {
					t.Errorf("unexpected error: got %v, want nil", gotErr)
				}

				if email != validUser.Email {
					t.Errorf("unexpected Email: got %v, want %v", email, validUser.Email)
				}
			}
		})
	}
}

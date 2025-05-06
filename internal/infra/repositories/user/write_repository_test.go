package repositories

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
	"gorm.io/gorm"
)

func TestUserRepository_CreateUser(t *testing.T) {
	validUser := entities.UserEntity{
		ID:             "user-id",
		ExternalUserID: "test-member-id",
		Email:          "test-email",
		EnterpriseID:   "test-enterprise-id",
	}

	tests := []struct {
		name    string
		user    entities.UserEntity
		wantErr error
		setup   func(t *testing.T, ent *entities.UserEntity) *gorm.DB
	}{
		{
			name:  "Create Success",
			user:  validUser,
			setup: setupTestDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			db := tt.setup(t, nil)

			repo := NewUserWriteRepository(db)

			gotErr := repo.CreateUser(ctx, tt.user)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("unexpected error: got %v, want %v", gotErr, tt.wantErr)
			}

			if tt.wantErr == nil {
				var gotEntity entities.UserEntity
				if err := db.First(&gotEntity, "id = ?", tt.user.ID).Error; err != nil {
					t.Errorf("failed to find created entity: %v", err)
				}

				if gotEntity.ID != tt.user.ID {
					t.Errorf("unexpected entity: got %v, want %v", gotEntity, tt.user)
				}
			}
		})
	}
}

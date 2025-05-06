package repositories

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
	"gorm.io/gorm"
)

type UserWriteRepository struct {
	db *gorm.DB
}

func NewUserWriteRepository(db *gorm.DB) repository.UserWriteRepositoryIF {
	return &UserWriteRepository{
		db: db,
	}
}

func (r UserWriteRepository) CreateUser(ctx context.Context, user entities.UserEntity) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

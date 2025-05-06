package repositories

import (
	"context"
	"errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
	"gorm.io/gorm"
)

type UserReadRepository struct {
	db *gorm.DB
}

func NewUserReadRepository(db *gorm.DB) repository.UserReadRepositoryIF {
	return &UserReadRepository{
		db: db,
	}
}

func (r *UserReadRepository) GetUserByID(ctx context.Context, userID, enterpriseID string) (entities.UserEntity, error) {
	var user entities.UserEntity
	err := r.db.WithContext(ctx).
		Where("id = ? AND enterprise_id = ?", userID, enterpriseID).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.UserEntity{}, nil
	}

	if err != nil {
		return entities.UserEntity{}, err
	}

	return user, nil
}

func (r UserReadRepository) GetUserByEmail(ctx context.Context, email, enterpriseId string) (entities.UserEntity, error) {
	var user entities.UserEntity
	err := r.db.WithContext(ctx).Where("email = ? AND enterprise_id = ? ", email, enterpriseId).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entities.UserEntity{}, nil
	}

	if err != nil {
		return entities.UserEntity{}, err
	}

	return user, err
}

func (r UserReadRepository) GetEmailByUserID(ctx context.Context, userID, enterpriseID string) (string, error) {
	var email string
	err := r.db.WithContext(ctx).
		Model(&entities.UserEntity{}).
		Where("id = ? AND enterprise_id = ?", userID, enterpriseID).
		Select("email").
		First(&email).Error

	return email, err
}

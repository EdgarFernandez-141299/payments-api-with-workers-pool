package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
)

type UserReadRepositoryIF interface {
	GetUserByID(ctx context.Context, userID, enterpriseID string) (entities.UserEntity, error)
	GetUserByEmail(ctx context.Context, email, enterpriseId string) (entities.UserEntity, error)
	GetEmailByUserID(ctx context.Context, userID, enterpriseID string) (string, error)
}

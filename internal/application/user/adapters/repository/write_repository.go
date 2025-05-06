package repository

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
)

type UserWriteRepositoryIF interface {
	CreateUser(ctx context.Context, user entities.UserEntity) error
}

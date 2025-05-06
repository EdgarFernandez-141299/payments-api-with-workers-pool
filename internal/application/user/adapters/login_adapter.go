package adapters

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	errorsBusiness "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/user/errors"
)

type DeunaLoginAdapter interface {
	Login(ctx context.Context, userId string) (string, error)
	LoginWitUserID(ctx context.Context, userId, enterpriseID string) (string, error)
}

type DeunaLoginAdapterImpl struct {
	client     resources.DeunaLoginResourceIF
	clientAuth resources.DeunaAuthResourceIF
	repository repository.UserReadRepositoryIF
}

func NewDeunaLoginAdapter(
	client resources.DeunaLoginResourceIF,
	repository repository.UserReadRepositoryIF,
	clientAuth resources.DeunaAuthResourceIF,
) DeunaLoginAdapter {
	return &DeunaLoginAdapterImpl{
		client:     client,
		repository: repository,
		clientAuth: clientAuth,
	}
}

func (u *DeunaLoginAdapterImpl) Login(ctx context.Context, userId string) (string, error) {
	tokenResponse, err := u.client.Login(ctx, request.LoginUserDeUnaRequestDTO{
		UserID: userId,
	})
	if err != nil {
		return "", errorsBusiness.NewLoginError(userId, err)
	}

	return tokenResponse.Token, nil
}

func (u *DeunaLoginAdapterImpl) LoginWitUserID(ctx context.Context, userId, enterpriseID string) (string, error) {
	user, err := u.repository.GetUserByID(ctx, userId, enterpriseID)

	if err != nil {
		return "", errorsBusiness.NewUserNotFoundError(userId, err)
	}

	tokenResponse, err := u.clientAuth.AuthUser(ctx, request.DeunaAuthUserRequestDTO{
		Email: user.EmailAlias,
	})

	if err != nil {
		return "", errorsBusiness.NewLoginError(userId, err)
	}

	return tokenResponse.AuthToken, nil
}

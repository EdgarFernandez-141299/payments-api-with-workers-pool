package queries

import (
	"context"
	"errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/repository"
	repositoryUser "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/repository"
	errorsBusiness "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/errors"
	errorsBusinessUser "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/user/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/response"
	"gorm.io/gorm"
)

type GetCardUsecaseIF interface {
	GetCardsByUserID(ctx context.Context, memberId, enterpriseId string) ([]response.CardResponse, error)
}

type GetCardUsecase struct {
	repository     repository.CardReadRepositoryIF
	repositoryUser repositoryUser.UserReadRepositoryIF
}

func NewGetCardUsecase(
	repository repository.CardReadRepositoryIF,
	repositoryUser repositoryUser.UserReadRepositoryIF) GetCardUsecaseIF {
	return &GetCardUsecase{
		repository:     repository,
		repositoryUser: repositoryUser,
	}
}

func (p *GetCardUsecase) GetCardsByUserID(
	ctx context.Context, userID, enterpriseID string,
) ([]response.CardResponse, error) {
	_, err := p.repositoryUser.GetUserByID(ctx, userID, enterpriseID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorsBusinessUser.NewUserNotFoundError(userID, err)
		}

		return []response.CardResponse{}, err
	}

	cards, err := p.repository.GetCardsByUserID(ctx, userID, enterpriseID)

	if err != nil {
		return []response.CardResponse{}, errorsBusiness.NewCardErrorGetDB(err)
	}

	return cards.ToDTO(), nil
}

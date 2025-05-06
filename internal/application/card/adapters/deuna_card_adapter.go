package adapters

import (
	"context"

	"gitlab.com/clubhub.ai1/go-libraries/observability/apm/decorators"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/deuna"
)

type CardAdapter interface {
	CreateCard(
		ctx context.Context,
		userID string,
		userType string,
		enterpriseID string,
		body request.CreateCardRequestDTO,
	) (response.CardResponseDataDTO, error)

	DeleteCard(
		ctx context.Context,
		cardID string,
		externalCardID string,
		userID string,
		enterpriseID string,
	) (response.DeleteCardResponseDTO, error)
}

type DeunaCardAdapterImpl struct {
	client       resources.DeunaCardResourceIF
	loginAdapter adapters.DeunaLoginAdapter
	userAdapter  deuna.CreateUserUseDeunaAdapterIF
}

func NewDeunaCardAdapter(
	client resources.DeunaCardResourceIF,
	loginAdapter adapters.DeunaLoginAdapter,
	userAdapter deuna.CreateUserUseDeunaAdapterIF,
) CardAdapter {
	return &DeunaCardAdapterImpl{
		userAdapter:  userAdapter,
		client:       client,
		loginAdapter: loginAdapter,
	}
}

func (d *DeunaCardAdapterImpl) CreateCard(
	oldContext context.Context,
	userID string,
	userType string,
	enterpriseID string,
	body request.CreateCardRequestDTO,
) (response.CardResponseDataDTO, error) {
	return decorators.TraceDecorator[response.CardResponseDataDTO](
		oldContext,
		"DeunaCardAdapterImpl.CreateCard",
		func(ctx context.Context, span decorators.Span) (response.CardResponseDataDTO, error) {
			user, err := d.userAdapter.GetOrCreateUser(ctx, userID, userType, enterpriseID)

			if err != nil {
				return response.CardResponseDataDTO{}, err
			}

			token, err := d.loginAdapter.Login(ctx, user.ExternalUserID)
			if err != nil {
				return response.CardResponseDataDTO{}, err
			}

			cardResponse, err := d.client.CreateCard(ctx, body, user.ExternalUserID, token)
			if err != nil {
				return response.CardResponseDataDTO{}, err
			}

			cardResponse.InternalUserID = user.ID

			return cardResponse, nil
		},
	)
}

func (d *DeunaCardAdapterImpl) DeleteCard(
	ctx context.Context,
	cardID string,
	externalCardID string,
	userID string,
	enterpriseID string,
) (response.DeleteCardResponseDTO, error) {
	user, err := d.userAdapter.GetUser(ctx, userID, enterpriseID)

	if err != nil {
		return response.DeleteCardResponseDTO{}, err
	}

	cardBody := request.DeleteCardRequestDTO{
		CardId: externalCardID,
		UserId: user.ExternalUserID,
	}

	token, err := d.loginAdapter.Login(ctx, cardBody.UserId)

	if err != nil {
		return response.DeleteCardResponseDTO{}, err
	}

	return d.client.DeleteCard(ctx, cardBody, token)
}

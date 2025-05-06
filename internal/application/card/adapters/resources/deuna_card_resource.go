package resources

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/response"
)

type DeunaCardResourceIF interface {
	CreateCard(
		ctx context.Context,
		body request.CreateCardRequestDTO,
		userID string, token string,
	) (response.CardResponseDataDTO, error)

	DeleteCard(
		ctx context.Context,
		body request.DeleteCardRequestDTO,
		token string,
	) (response.DeleteCardResponseDTO, error)
}

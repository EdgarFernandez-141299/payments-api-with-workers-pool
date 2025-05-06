package usecases

import (
	"context"
	"time"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/response"
)

func (c *CollectionAccountRouteUsecase) Disable(
	ctx context.Context, id, enterpriseId string,
) (response.CollectionAccountRouteDisableResponse, error) {
	err := c.repository.Disable(ctx, id, enterpriseId)
	if err != nil {
		return response.CollectionAccountRouteDisableResponse{}, err
	}

	return response.CollectionAccountRouteDisableResponse{
		ID:         id,
		Status:     "disabled",
		DisabledAt: time.Now().String(),
	}, nil
}

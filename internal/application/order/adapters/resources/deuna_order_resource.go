package resources

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources/dto/response"
)

type DeunaOrderResourceIF interface {
	CreateOrder(
		ctx context.Context,
		body request.CreateDeunaOrderRequestDTO,
	) (response.DeunaOrderResponseDTO, error)

	GetOrder(
		ctx context.Context,
		orderToken string,
	) (response.DeunaOrderResponseDTO, error)

	ExpireOrder(
		ctx context.Context,
		orderToken string,
	) error
}

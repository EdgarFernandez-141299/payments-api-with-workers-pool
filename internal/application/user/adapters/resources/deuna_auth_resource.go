package resources

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

type DeunaAuthResourceIF interface {
	AuthUser(
		ctx context.Context,
		request request.DeunaAuthUserRequestDTO,
	) (response.DeunaAuthResponseDTO, error)
}

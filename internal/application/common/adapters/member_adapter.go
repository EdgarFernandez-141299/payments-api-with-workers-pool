package adapters

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

type MemberAdapterIF interface {
	GetUserProfileInfo(
		oldCtx context.Context,
		userId string,
		enterpriseId string,
	) (*response.UserProfileInfoDTO, error)
}

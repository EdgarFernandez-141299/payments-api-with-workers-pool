package resources

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

type MemberAPIResourceIF interface {
	GetMemberByID(ctx context.Context, id, enterpriseId string) (response.MemberDTO, error)
	GetUserProfileInfo(oldCtx context.Context, userId string, enterpriseId string) (response.UserProfileInfoDTO, error)
}

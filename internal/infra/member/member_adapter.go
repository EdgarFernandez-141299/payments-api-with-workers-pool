package member

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters"
	memberResources "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
)

type memberAdapterImpl struct {
	memberResource memberResources.MemberAPIResourceIF
}

func NewMemberAdapterIF(
	memberResource memberResources.MemberAPIResourceIF,
) adapters.MemberAdapterIF {
	return &memberAdapterImpl{
		memberResource: memberResource,
	}
}

func (mai *memberAdapterImpl) GetUserProfileInfo(
	ctx context.Context,
	userId string,
	enterpriseId string,
) (*response.UserProfileInfoDTO, error) {

	userProfileInfo, err := mai.memberResource.GetUserProfileInfo(
		ctx,
		userId,
		enterpriseId,
	)

	if err != nil {
		return nil, err
	}

	return &userProfileInfo, nil
}

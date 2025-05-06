package member

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

func TestGetUserProfileInfo(t *testing.T) {
	memberResourceMock := resources.NewMemberAPIResourceIF(t)

	t.Run("should return user profile info", func(t *testing.T) {
		memberResourceMock := resources.NewMemberAPIResourceIF(t)

		ctx := context.Background()
		memberAdapter := NewMemberAdapterIF(memberResourceMock)

		memberResourceMock.On("GetUserProfileInfo", mock.Anything, mock.Anything, mock.Anything).
			Return(response.UserProfileInfoDTO{
				Email:     "test@test.com",
				FirstName: "Test",
				LastName:  "Test",
			}, nil)

		userProfileInfo, err := memberAdapter.GetUserProfileInfo(ctx, "123", "456")

		assert.NoError(t, err)
		assert.NotNil(t, userProfileInfo)
		memberResourceMock.ExpectedCalls = nil

	})

	t.Run("should return error", func(t *testing.T) {

		ctx := context.Background()
		memberAdapter := NewMemberAdapterIF(memberResourceMock)

		memberResourceMock.On("GetUserProfileInfo", mock.Anything, mock.Anything, mock.Anything).
			Return(response.UserProfileInfoDTO{}, errors.New("error"))

		userProfileInfo, err := memberAdapter.GetUserProfileInfo(ctx, "123", "456")

		assert.Error(t, err)
		assert.Nil(t, userProfileInfo)
		memberResourceMock.ExpectedCalls = nil
	})
}

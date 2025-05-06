package adapters

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	deunaErrors "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/user/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
)

func TestDeunaLoginAdapterImpl_Login(t *testing.T) {
	type testCase struct {
		name        string
		userId      string
		mockSetup   func(client *resources.DeunaLoginResourceIF)
		expectedRes string
		expectedErr error
	}

	tests := []testCase{
		{
			name:   "successful login",
			userId: "validUserId",
			mockSetup: func(client *resources.DeunaLoginResourceIF) {
				client.On("Login", mock.Anything, request.LoginUserDeUnaRequestDTO{
					UserID: "validUserId",
				}).Return(response.LoginResponseDTO{
					Token: "validToken",
				}, nil).Once()
			},
			expectedRes: "validToken",
			expectedErr: nil,
		},
		{
			name:   "login request fails with error",
			userId: "failedUserId",
			mockSetup: func(client *resources.DeunaLoginResourceIF) {
				client.On("Login", mock.Anything, request.LoginUserDeUnaRequestDTO{
					UserID: "failedUserId",
				}).Return(response.LoginResponseDTO{}, assert.AnError).Once()
			},
			expectedRes: "",
			expectedErr: deunaErrors.NewLoginError("failedUserId", assert.AnError),
		},
		{
			name:   "empty user ID",
			userId: "",
			mockSetup: func(client *resources.DeunaLoginResourceIF) {
				client.On("Login", mock.Anything, request.LoginUserDeUnaRequestDTO{
					UserID: "",
				}).Return(response.LoginResponseDTO{}, errors.New("user ID is empty")).Once()
			},
			expectedRes: "",
			expectedErr: deunaErrors.NewLoginError("", errors.New("user ID is empty")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientMock := resources.NewDeunaLoginResourceIF(t)
			tc.mockSetup(clientMock)

			adapter := DeunaLoginAdapterImpl{client: clientMock}

			res, err := adapter.Login(context.Background(), tc.userId)

			assert.Equal(t, tc.expectedRes, res)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			clientMock.AssertExpectations(t)
		})
	}
}

package adapters

import (
	"context"
	"errors"
	"testing"

	response2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
	deuna2 "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/deuna"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/response"
)

type mocks struct {
	client       *resources.DeunaCardResourceIF
	loginAdapter *adapters.DeunaLoginAdapter
	userAdapter  *deuna2.CreateUserUseDeunaAdapterIF
}

func setupMocks(t *testing.T) *mocks {
	return &mocks{
		client:       resources.NewDeunaCardResourceIF(t),
		loginAdapter: adapters.NewDeunaLoginAdapter(t),
		userAdapter:  deuna2.NewCreateUserUseDeunaAdapterIF(t),
	}
}

func TestCreateCard(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*mocks)
		userID         string
		userType       string
		enterpriseID   string
		request        request.CreateCardRequestDTO
		expectedResult response.CardResponseDataDTO
		expectedError  error
	}{
		{
			name: "successful card creation",
			setupMocks: func(m *mocks) {
				m.userAdapter.On("GetOrCreateUser", mock.Anything, "user1", "usertype1", "enterprise1").
					Return(response2.CreatedUserResponse{ID: "user1", Email: "email@email.com", ExternalUserID: "externalUser1"}, nil)
				m.loginAdapter.On("Login", mock.Anything, "externalUser1").Return("valid_token", nil)
				m.client.On("CreateCard", mock.Anything, mock.Anything, "externalUser1", "valid_token").
					Return(response.CardResponseDataDTO{
						ID:         "card1",
						UserID:     "user1",
						CardHolder: "John Doe",
					}, nil)
			},
			userID:       "user1",
			userType:     "usertype1",
			enterpriseID: "enterprise1",
			request: request.CreateCardRequestDTO{
				CardHolder: "John Doe",
				CardNumber: "4111111111111111",
			},
			expectedResult: response.CardResponseDataDTO{
				ID:             "card1",
				UserID:         "user1",
				CardHolder:     "John Doe",
				InternalUserID: "user1",
			},
			expectedError: nil,
		},
		{
			name: "error getting or creating user",
			setupMocks: func(m *mocks) {
				m.userAdapter.On("GetOrCreateUser", mock.Anything, "user1", "usertype1", "enterprise1").
					Return(response2.CreatedUserResponse{}, errors.New("user error"))
			},
			userID:         "user1",
			userType:       "usertype1",
			enterpriseID:   "enterprise1",
			request:        request.CreateCardRequestDTO{},
			expectedResult: response.CardResponseDataDTO{},
			expectedError:  errors.New("user error"),
		},
		{
			name: "error logging in user",
			setupMocks: func(m *mocks) {
				m.userAdapter.On("GetOrCreateUser", mock.Anything, "user1", "usertype1", "enterprise1").
					Return(response2.CreatedUserResponse{ID: "user1", ExternalUserID: "user1"}, nil)
				m.loginAdapter.On("Login", mock.Anything, "user1").Return("", errors.New("login error"))
			},
			userID:         "user1",
			userType:       "usertype1",
			enterpriseID:   "enterprise1",
			request:        request.CreateCardRequestDTO{},
			expectedResult: response.CardResponseDataDTO{},
			expectedError:  errors.New("login error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupMocks(t)
			adapter := NewDeunaCardAdapter(mocks.client, mocks.loginAdapter, mocks.userAdapter)
			tt.setupMocks(mocks)

			got, err := adapter.CreateCard(context.Background(), tt.userID, tt.userType, tt.enterpriseID, tt.request)

			assert.Equal(t, tt.expectedResult, got)
			assert.Equal(t, tt.expectedError, err)
			mocks.client.AssertExpectations(t)
			mocks.loginAdapter.AssertExpectations(t)
			mocks.userAdapter.AssertExpectations(t)
		})
	}
}

func TestDeleteCard(t *testing.T) {
	tests := []struct {
		name           string
		setupMocks     func(*mocks)
		cardID         string
		externalCard   string
		memberID       string
		enterpriseID   string
		expectedResult response.DeleteCardResponseDTO
		expectedError  error
	}{
		{
			name: "successful card deletion",
			setupMocks: func(m *mocks) {
				m.userAdapter.On("GetUser", mock.Anything, "member1", "enterprise1").
					Return(response2.CreatedUserResponse{ID: "user1", ExternalUserID: "externalUser1"}, nil)
				m.loginAdapter.On("Login", mock.Anything, "externalUser1").Return("valid_token", nil)
				m.client.On("DeleteCard", mock.Anything,
					request.DeleteCardRequestDTO{CardId: "cardExt", UserId: "externalUser1"}, "valid_token").
					Return(response.DeleteCardResponseDTO{
						Code:    "200",
						Message: "Card deleted successfully",
					}, nil)
			},
			cardID:       "card1",
			memberID:     "member1",
			externalCard: "cardExt",
			enterpriseID: "enterprise1",
			expectedResult: response.DeleteCardResponseDTO{
				Code:    "200",
				Message: "Card deleted successfully",
			},
			expectedError: nil,
		},
		{
			name: "error getting user",
			setupMocks: func(m *mocks) {
				m.userAdapter.On("GetUser", mock.Anything, "member1", "enterprise1").
					Return(response2.CreatedUserResponse{}, errors.New("user error"))
			},
			cardID:         "card1",
			memberID:       "member1",
			enterpriseID:   "enterprise1",
			expectedResult: response.DeleteCardResponseDTO{},
			expectedError:  errors.New("user error"),
		},
		{
			name: "error logging in user",
			setupMocks: func(m *mocks) {
				m.userAdapter.On("GetUser", mock.Anything, "member1", "enterprise1").
					Return(response2.CreatedUserResponse{ID: "user1", ExternalUserID: "externalUser1"}, nil)
				m.loginAdapter.On("Login", mock.Anything, "externalUser1").
					Return("", errors.New("login error"))
			},
			cardID:         "card1",
			memberID:       "member1",
			enterpriseID:   "enterprise1",
			expectedResult: response.DeleteCardResponseDTO{},
			expectedError:  errors.New("login error"),
		},
		{
			name: "error deleting card",
			setupMocks: func(m *mocks) {
				m.userAdapter.On("GetUser", mock.Anything, "member1", "enterprise1").
					Return(response2.CreatedUserResponse{ID: "user1", ExternalUserID: "externalUser1"}, nil)
				m.loginAdapter.On("Login", mock.Anything, "externalUser1").Return("valid_token", nil)
				m.client.On("DeleteCard", mock.Anything,
					request.DeleteCardRequestDTO{CardId: "ext1", UserId: "externalUser1"}, "valid_token").
					Return(response.DeleteCardResponseDTO{}, errors.New("delete card error"))
			},
			cardID:         "card1",
			memberID:       "member1",
			enterpriseID:   "enterprise1",
			externalCard:   "ext1",
			expectedResult: response.DeleteCardResponseDTO{},
			expectedError:  errors.New("delete card error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocks := setupMocks(t)
			adapter := NewDeunaCardAdapter(mocks.client, mocks.loginAdapter, mocks.userAdapter)
			tt.setupMocks(mocks)

			got, err := adapter.DeleteCard(context.Background(), tt.cardID, tt.externalCard, tt.memberID, tt.enterpriseID)

			assert.Equal(t, tt.expectedResult, got)
			assert.Equal(t, tt.expectedError, err)

			mocks.client.AssertExpectations(t)
			mocks.loginAdapter.AssertExpectations(t)
			mocks.userAdapter.AssertExpectations(t)
		})
	}
}

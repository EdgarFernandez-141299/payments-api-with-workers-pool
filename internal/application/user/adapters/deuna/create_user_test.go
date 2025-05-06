package deuna

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	requestDeUna "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	responseDeUna "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user/dto/request"
	userResponsedto "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	mockResource "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

var (
	ctx          = context.Background()
	enterpriseId = "enterpriseId"
	userType     = "userType"
	userID       = "userID"
	emails       = []response.EmailDTO{
		{Email: "email@gmail.com", IsDefault: true},
	}
)

func TestCreateUser(t *testing.T) {
	t.Run("error getting member by id member-api", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)

		requestMock := request.CreateUserRequest{
			UserID: userID,
		}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(response.UserProfileInfoDTO{}, assert.AnError)

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		_, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "USER_NOT_FOUND"))

		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})
	t.Run("error getting main email from member", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)

		requestMock := request.CreateUserRequest{
			UserID: userID,
		}

		memberFound := response.UserProfileInfoDTO{
			UserID: userID,
		}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(memberFound, nil)

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		_, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "USER_EMAIL_NOT_FOUND"))

		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})
	t.Run("should retunr an error user not found error", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)

		requestMock := request.CreateUserRequest{
			UserID: userID,
		}

		memberFound := response.UserProfileInfoDTO{}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(memberFound, errors.New("USER_NOT_FOUND"))

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		_, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "USER_NOT_FOUND"))

		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should return an error user already exists", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)

		requestMock := request.CreateUserRequest{
			UserID: userID,
		}

		memberFound := response.UserProfileInfoDTO{
			UserID: userID,
			Email:  "email@gmail.com",
		}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(memberFound, nil)
		repositoryMock.On("GetUserByEmail", ctx, "email@gmail.com", enterpriseId).Return(entities.UserEntity{
			ID: "userID",
		}, nil)

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		_, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "USER_ALREADY_EXISTS"))

		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should return an error  phone main not found in", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)

		requestMock := request.CreateUserRequest{
			UserID: userID,
		}

		memberFound := response.UserProfileInfoDTO{
			UserID: userID,
			Email:  "email@gmail.com",
		}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(memberFound, nil)
		repositoryMock.On("GetUserByEmail", ctx, "email@gmail.com", enterpriseId).Return(entities.UserEntity{}, nil)

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		_, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "USER_PHONE_NOT_FOUND"))

		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should return error create user DEUNA", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)

		requestMock := request.CreateUserRequest{
			UserID:   userID,
			UserType: userType,
		}

		memberFound := response.UserProfileInfoDTO{
			UserID:       userID,
			Email:        "email@gmail.com",
			PrimaryPhone: response.PhoneInfo{Number: "123456789", CountryCode: "+52"},
			FirstName:    "John",
			LastName:     "Doe",
			Address: responseDeUna.AddressInfo{
				AddressLine: "address",
				ZipCode:     "123456",
				City:        "city",
				State:       "state",
			},
		}

		emailAlias := requestDeUna.NewUserEmailAlias(userID)

		newMember := requestDeUna.CreateUserRequestDTO{
			Email:     emailAlias,
			FirstName: memberFound.FirstName,
			LastName:  memberFound.LastName,
			Phone:     "123456789",
		}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(memberFound, nil)
		repositoryMock.On("GetUserByEmail", ctx, "email@gmail.com", enterpriseId).Return(entities.UserEntity{}, nil)
		deunaAdapterMock.On("CreateUser", ctx, newMember).Return(response.CreatedUserResponse{}, assert.AnError)

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		_, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.NotNil(t, err)
		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should retorun create user db", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)
		requestMock := request.CreateUserRequest{
			UserID:   userID,
			UserType: userType,
		}

		memberFound := response.UserProfileInfoDTO{
			UserID:       userID,
			Email:        "email@gmail.com",
			PrimaryPhone: response.PhoneInfo{Number: "123456789", CountryCode: "+52"},
			FirstName:    "John",
			LastName:     "Doe",
			Address: responseDeUna.AddressInfo{
				AddressLine: "address",
				ZipCode:     "123456",
				City:        "city",
				State:       "state",
			},
		}

		emailAlias := requestDeUna.NewUserEmailAlias(userID)

		newMember := requestDeUna.CreateUserRequestDTO{
			Email:     emailAlias,
			FirstName: memberFound.FirstName,
			LastName:  memberFound.LastName,
			Phone:     "123456789",
		}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(memberFound, nil)
		repositoryMock.On("GetUserByEmail", ctx, "email@gmail.com", enterpriseId).Return(entities.UserEntity{}, nil)
		deunaAdapterMock.On("CreateUser", ctx, newMember).Return(response.CreatedUserResponse{
			Token:  "token",
			UserID: "userId",
		}, nil)

		writeRepository.On("CreateUser", ctx, mock.MatchedBy(func(entities.UserEntity) bool {
			return true
		})).Return(assert.AnError)

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		_, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.NotNil(t, err)
		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})

	t.Run("should create user successfully", func(t *testing.T) {
		memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
		deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
		deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
		repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
		writeRepository := mockRepository.NewUserWriteRepositoryIF(t)
		requestMock := request.CreateUserRequest{
			UserID: userID,
		}

		memberFound := response.UserProfileInfoDTO{
			UserID:       userID,
			Email:        "email@gmail.com",
			PrimaryPhone: response.PhoneInfo{Number: "123456789", CountryCode: "+52"},
			FirstName:    "John",
			LastName:     "Doe",
			Address: responseDeUna.AddressInfo{
				AddressLine: "address",
				ZipCode:     "123456",
				City:        "city",
				State:       "state",
			},
		}

		emailAlias := requestDeUna.NewUserEmailAlias(userID)

		newMember := requestDeUna.CreateUserRequestDTO{
			Email:     emailAlias,
			FirstName: memberFound.FirstName,
			LastName:  memberFound.LastName,
			Phone:     "123456789",
		}

		memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(memberFound, nil)
		repositoryMock.On("GetUserByEmail", ctx, "email@gmail.com", enterpriseId).Return(entities.UserEntity{}, nil)
		deunaAdapterMock.On("CreateUser", ctx, newMember).Return(response.CreatedUserResponse{
			Token:  "token",
			UserID: "userID",
		}, nil)

		writeRepository.On("CreateUser", ctx, mock.MatchedBy(func(user entities.UserEntity) bool {
			return user.Email == "email@gmail.com" && user.EnterpriseID == enterpriseId
		})).Return(nil)

		userApp := NewCreateUserUseCases(
			memberAdapterMock,
			deunaAdapterMock,
			deunaAuthAdapterMock,
			repositoryMock,
			writeRepository,
		)

		createdUser, err := userApp.Create(ctx, requestMock, enterpriseId)

		assert.Nil(t, err)
		assert.Equal(t, "userID", createdUser.ID)
		assert.Equal(t, "email@gmail.com", createdUser.Email)
		assert.Equal(t, memberFound.FirstName, createdUser.FirstName)
		assert.Equal(t, memberFound.LastName, createdUser.LastName)

		memberAdapterMock.AssertExpectations(t)
		deunaAdapterMock.AssertExpectations(t)
		repositoryMock.AssertExpectations(t)
	})
}

func TestGetOrCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(memberAdapterMock *mockResource.MemberAPIResourceIF,
			deunaAdapterMock *mockResource.DeUnaUserResourceIF,
			repositoryMock *mockRepository.UserReadRepositoryIF,
			userWriteRepository *mockRepository.UserWriteRepositoryIF,
		)
		expectedResp  userResponsedto.CreatedUserResponse
		expectedError error
	}{
		{
			name: "user exists in repository",
			setupMocks: func(memberAdapterMock *mockResource.MemberAPIResourceIF,
				deunaAdapterMock *mockResource.DeUnaUserResourceIF,
				repositoryMock *mockRepository.UserReadRepositoryIF,
				userWriteRepository *mockRepository.UserWriteRepositoryIF,
			) {
				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).Return(entities.UserEntity{
					ID:             "userID",
					Email:          "email@gmail.com",
					ExternalUserID: "externaluserID",
				}, nil)
			},
			expectedResp: userResponsedto.CreatedUserResponse{
				ID:             "userID",
				Email:          "email@gmail.com",
				ExternalUserID: "externaluserID",
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			setupMocks: func(memberAdapterMock *mockResource.MemberAPIResourceIF,
				deunaAdapterMock *mockResource.DeUnaUserResourceIF,
				repositoryMock *mockRepository.UserReadRepositoryIF,
				userWriteRepository *mockRepository.UserWriteRepositoryIF) {
				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).
					Return(entities.UserEntity{}, assert.AnError)
			},
			expectedResp:  userResponsedto.CreatedUserResponse{},
			expectedError: assert.AnError,
		},
		{
			name: "user needs to be created",
			setupMocks: func(memberAdapterMock *mockResource.MemberAPIResourceIF,
				deunaAdapterMock *mockResource.DeUnaUserResourceIF,
				repositoryMock *mockRepository.UserReadRepositoryIF,
				userWriteRepository *mockRepository.UserWriteRepositoryIF) {
				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).
					Return(entities.UserEntity{}, nil)
				repositoryMock.On("GetUserByEmail", ctx, "email@gmail.com", enterpriseId).
					Return(entities.UserEntity{}, nil)
				memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).
					Return(response.UserProfileInfoDTO{
						UserID:       userID,
						Email:        "email@gmail.com",
						PrimaryPhone: response.PhoneInfo{Number: "123456789", CountryCode: "+52"},
						Address: responseDeUna.AddressInfo{
							AddressLine: "address",
							ZipCode:     "123456",
							City:        "city",
							State:       "state",
						},
					}, nil)
				deunaAdapterMock.On("CreateUser", ctx, mock.MatchedBy(
					func(req requestDeUna.CreateUserRequestDTO) bool {
						return req.Email == requestDeUna.NewUserEmailAlias(userID)
					})).Return(response.CreatedUserResponse{
					Token:  "token",
					UserID: userID,
				}, nil)
				userWriteRepository.On("CreateUser", ctx, mock.Anything).Return(nil)
			},
			expectedResp: userResponsedto.CreatedUserResponse{
				ID:                    userID,
				Email:                 "email@gmail.com",
				ExternalUserID:        userID,
				PaymentsExternalEmail: requestDeUna.NewUserEmailAlias(userID),
			},
			expectedError: nil,
		},
		{
			name: "create user error",
			setupMocks: func(memberAdapterMock *mockResource.MemberAPIResourceIF,
				deunaAdapterMock *mockResource.DeUnaUserResourceIF,
				repositoryMock *mockRepository.UserReadRepositoryIF,
				userWriteRepository *mockRepository.UserWriteRepositoryIF) {
				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).Return(entities.UserEntity{}, nil)
				repositoryMock.On("GetUserByEmail", ctx, "email@gmail.com", enterpriseId).Return(entities.UserEntity{}, nil)
				memberAdapterMock.On("GetUserProfileInfo", ctx, userID, enterpriseId).Return(response.UserProfileInfoDTO{
					UserID:       userID,
					Email:        "email@gmail.com",
					PrimaryPhone: response.PhoneInfo{Number: "123456789", CountryCode: "+52"},
					Address: responseDeUna.AddressInfo{
						AddressLine: "address",
						ZipCode:     "123456",
						City:        "city",
						State:       "state",
					},
				}, nil)
				deunaAdapterMock.On("CreateUser", ctx, mock.Anything).Return(response.CreatedUserResponse{}, assert.AnError)
			},
			expectedResp:  userResponsedto.CreatedUserResponse{},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
			deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
			deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
			repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
			userWriteRepository := mockRepository.NewUserWriteRepositoryIF(t)

			tt.setupMocks(memberAdapterMock, deunaAdapterMock, repositoryMock, userWriteRepository)

			userApp := NewCreateUserUseCases(
				memberAdapterMock,
				deunaAdapterMock,
				deunaAuthAdapterMock,
				repositoryMock,
				userWriteRepository,
			)

			resp, err := userApp.GetOrCreateUser(ctx, userID, userType, enterpriseId)

			assert.Equal(t, tt.expectedResp, resp)

			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			memberAdapterMock.AssertExpectations(t)
			deunaAdapterMock.AssertExpectations(t)
			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(repositoryMock *mockRepository.UserReadRepositoryIF)
		expectedResp  userResponsedto.CreatedUserResponse
		expectedError error
	}{
		{
			name: "user exists in repository",
			setupMocks: func(repositoryMock *mockRepository.UserReadRepositoryIF) {
				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).Return(entities.UserEntity{
					ID:             "user-id",
					Email:          "email@gmail.com",
					ExternalUserID: "externaluserID",
				}, nil)
			},
			expectedResp: userResponsedto.CreatedUserResponse{
				ID:             "user-id",
				Email:          "email@gmail.com",
				ExternalUserID: "externaluserID",
			},
			expectedError: nil,
		},
		{
			name: "repository error",
			setupMocks: func(repositoryMock *mockRepository.UserReadRepositoryIF) {
				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).Return(entities.UserEntity{}, assert.AnError)
			},
			expectedResp:  userResponsedto.CreatedUserResponse{},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := mockRepository.NewUserReadRepositoryIF(t)

			if tt.setupMocks != nil {
				tt.setupMocks(repositoryMock)
			}

			userApp := &CreateUserUseDeunaAdapterImpl{
				repository: repositoryMock,
			}

			resp, err := userApp.GetUser(ctx, userID, enterpriseId)

			assert.Equal(t, tt.expectedResp, resp)

			if tt.expectedError != nil {
				assert.Error(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			repositoryMock.AssertExpectations(t)
		})
	}
}

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(memberAdapterMock *mockResource.MemberAPIResourceIF,
			deunaAdapterMock *mockResource.DeUnaUserResourceIF,
			deunaAuthAdapterMock *mockResource.DeunaAuthResourceIF,
			repositoryMock *mockRepository.UserReadRepositoryIF,
			userWriteRepository *mockRepository.UserWriteRepositoryIF)
		request       request.CreateUserRequest
		expectedResp  userResponsedto.UserValidatedResponse
		expectedError error
	}{
		{
			name: "successful validation",
			setupMocks: func(memberAdapterMock *mockResource.MemberAPIResourceIF,
				deunaAdapterMock *mockResource.DeUnaUserResourceIF,
				deunaAuthAdapterMock *mockResource.DeunaAuthResourceIF,
				repositoryMock *mockRepository.UserReadRepositoryIF,
				userWriteRepository *mockRepository.UserWriteRepositoryIF) {

				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).Return(entities.UserEntity{
					ID:             userID,
					Email:          "email@gmail.com",
					ExternalUserID: "external-id",
					EmailAlias:     requestDeUna.NewUserEmailAlias(userID),
				}, nil)

				deunaAuthAdapterMock.On("AuthUser", ctx, requestDeUna.DeunaAuthUserRequestDTO{
					Email: requestDeUna.NewUserEmailAlias(userID),
				}).Return(responseDeUna.DeunaAuthResponseDTO{
					AuthToken: "auth-token-123",
				}, nil)
			},
			request: request.CreateUserRequest{
				UserID:   userID,
				UserType: userType,
			},
			expectedResp: userResponsedto.UserValidatedResponse{
				ID:                userID,
				Email:             "email@gmail.com",
				ExternalUserID:    "external-id",
				ExternalAuthToken: "auth-token-123",
			},
			expectedError: nil,
		},
		{
			name: "error getting user",
			setupMocks: func(memberAdapterMock *mockResource.MemberAPIResourceIF,
				deunaAdapterMock *mockResource.DeUnaUserResourceIF,
				deunaAuthAdapterMock *mockResource.DeunaAuthResourceIF,
				repositoryMock *mockRepository.UserReadRepositoryIF,
				userWriteRepository *mockRepository.UserWriteRepositoryIF) {

				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).Return(entities.UserEntity{}, assert.AnError)
			},
			request: request.CreateUserRequest{
				UserID:   userID,
				UserType: userType,
			},
			expectedResp:  userResponsedto.UserValidatedResponse{},
			expectedError: assert.AnError,
		},
		{
			name: "error in deuna auth",
			setupMocks: func(memberAdapterMock *mockResource.MemberAPIResourceIF,
				deunaAdapterMock *mockResource.DeUnaUserResourceIF,
				deunaAuthAdapterMock *mockResource.DeunaAuthResourceIF,
				repositoryMock *mockRepository.UserReadRepositoryIF,
				userWriteRepository *mockRepository.UserWriteRepositoryIF) {

				repositoryMock.On("GetUserByID", ctx, userID, enterpriseId).Return(entities.UserEntity{
					ID:             userID,
					Email:          "email@gmail.com",
					ExternalUserID: "external-id",
					EmailAlias:     requestDeUna.NewUserEmailAlias(userID),
				}, nil)

				deunaAuthAdapterMock.On("AuthUser", ctx, requestDeUna.DeunaAuthUserRequestDTO{
					Email: requestDeUna.NewUserEmailAlias(userID),
				}).Return(responseDeUna.DeunaAuthResponseDTO{}, assert.AnError)
			},
			request: request.CreateUserRequest{
				UserID:   userID,
				UserType: userType,
			},
			expectedResp:  userResponsedto.UserValidatedResponse{},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memberAdapterMock := mockResource.NewMemberAPIResourceIF(t)
			deunaAdapterMock := mockResource.NewDeUnaUserResourceIF(t)
			deunaAuthAdapterMock := mockResource.NewDeunaAuthResourceIF(t)
			repositoryMock := mockRepository.NewUserReadRepositoryIF(t)
			userWriteRepository := mockRepository.NewUserWriteRepositoryIF(t)

			tt.setupMocks(memberAdapterMock, deunaAdapterMock, deunaAuthAdapterMock, repositoryMock, userWriteRepository)

			userApp := NewCreateUserUseCases(
				memberAdapterMock,
				deunaAdapterMock,
				deunaAuthAdapterMock,
				repositoryMock,
				userWriteRepository,
			)

			resp, err := userApp.ValidateUser(ctx, tt.request, enterpriseId)

			assert.Equal(t, tt.expectedResp, resp)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			memberAdapterMock.AssertExpectations(t)
			deunaAdapterMock.AssertExpectations(t)
			deunaAuthAdapterMock.AssertExpectations(t)
			repositoryMock.AssertExpectations(t)
			userWriteRepository.AssertExpectations(t)
		})
	}
}

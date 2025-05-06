package deuna

import (
	"context"
	"errors"

	commonResources "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources"
	requestDeUna "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/utils"
	errorsBusiness "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/user/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/user/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
)

type CreateUserUseDeunaAdapterIF interface {
	Create(ctx context.Context,
		request request.CreateUserRequest,
		enterpriseId string,
	) (response.CreatedUserResponse, error)
	GetOrCreateUser(
		ctx context.Context,
		userID string,
		userType string,
		enterpriseID string,
	) (response.CreatedUserResponse, error)
	GetUser(
		ctx context.Context,
		memberId string,
		enterpriseId string,
	) (response.CreatedUserResponse, error)
	ValidateUser(
		ctx context.Context,
		request request.CreateUserRequest,
		enterpriseId string,
	) (response.UserValidatedResponse, error)
}

type CreateUserUseDeunaAdapterImpl struct {
	member              commonResources.MemberAPIResourceIF
	deuna               resources.DeUnaUserResourceIF
	deunaAuth           resources.DeunaAuthResourceIF
	repository          repository.UserReadRepositoryIF
	userWriteRepository repository.UserWriteRepositoryIF
}

func NewCreateUserUseCases(
	member commonResources.MemberAPIResourceIF,
	deuna resources.DeUnaUserResourceIF,
	deunaAuth resources.DeunaAuthResourceIF,
	repository repository.UserReadRepositoryIF,
	userWriteRepository repository.UserWriteRepositoryIF,
) CreateUserUseDeunaAdapterIF {
	return &CreateUserUseDeunaAdapterImpl{
		member:              member,
		deuna:               deuna,
		deunaAuth:           deunaAuth,
		repository:          repository,
		userWriteRepository: userWriteRepository,
	}
}

func (u *CreateUserUseDeunaAdapterImpl) Create(
	ctx context.Context,
	request request.CreateUserRequest,
	enterpriseID string,
) (response.CreatedUserResponse, error) {
	userFound, err := u.member.GetUserProfileInfo(ctx, request.UserID, enterpriseID)
	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewUserNotFoundError(request.UserID, err)
	}

	email, err := utils.GetUserEmail(userFound)
	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewUserEmailNotFoundError(request.UserID, err)
	}

	userExists, err := u.repository.GetUserByEmail(ctx, email, enterpriseID)
	if err != nil {
		return response.CreatedUserResponse{}, err
	}

	if !userExists.IsEmpty() {
		return response.CreatedUserResponse{},
			errorsBusiness.NewUserAlreadyExistError(request.UserID, errors.New("user already exists"))
	}

	phoneInfo, err := utils.GetUserPhone(userFound)
	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewUserPhoneNotFoundError(request.UserID, err)
	}

	billingInformation, err := utils.GetUserBillingInformation(userFound)
	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewUserBillingInformationNotFoundError(request.UserID, err)
	}

	emailAlias := requestDeUna.NewUserEmailAlias(request.UserID)

	newMember := requestDeUna.CreateUserRequestDTO{
		Email:     emailAlias,
		FirstName: userFound.FirstName,
		LastName:  userFound.LastName,
		Phone:     phoneInfo.Number,
	}

	memberCreatedDeUna, err := u.deuna.CreateUser(ctx, newMember)
	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewUserFailCreateDeUnaError(request.UserID, err)
	}

	newUser := entities.UserEntity{
		ID:             request.UserID,
		UserType:       request.UserType,
		ExternalUserID: memberCreatedDeUna.UserID,
		Email:          email,
		Address:        billingInformation.AddressLine,
		Zip:            billingInformation.ZipCode,
		City:           billingInformation.City,
		State:          billingInformation.State,
		CountryCode:    phoneInfo.CountryCode,
		Phone:          phoneInfo.Number,
		EnterpriseID:   enterpriseID,
		EmailAlias:     emailAlias,
	}

	user := entities.NewUserEntity(newUser)

	err = u.userWriteRepository.CreateUser(ctx, user)
	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewuserCreateError(err)
	}

	return response.CreatedUserResponse{
		ID:                    user.ID,
		FirstName:             userFound.FirstName,
		LastName:              userFound.LastName,
		Email:                 newUser.Email,
		ExternalUserID:        memberCreatedDeUna.UserID,
		PaymentsExternalEmail: newMember.Email,
	}, nil
}

func (u *CreateUserUseDeunaAdapterImpl) GetOrCreateUser(
	ctx context.Context,
	userID string,
	userType string,
	enterpriseID string,
) (response.CreatedUserResponse, error) {
	user, err := u.repository.GetUserByID(ctx, userID, enterpriseID)
	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewUserGettingDBError(userID, err)
	}

	if !user.IsEmpty() {
		return response.CreatedUserResponse{
			ID:                    user.ID,
			Email:                 user.Email,
			ExternalUserID:        user.ExternalUserID,
			PaymentsExternalEmail: user.EmailAlias,
		}, nil
	}

	newUser := request.CreateUserRequest{
		UserID: userID,
	}

	createdUser, err := u.Create(ctx, newUser, enterpriseID)
	if err != nil {
		return response.CreatedUserResponse{}, err
	}

	return createdUser, nil
}

func (u *CreateUserUseDeunaAdapterImpl) GetUser(
	ctx context.Context,
	memberId string,
	enterpriseId string,
) (response.CreatedUserResponse, error) {
	user, err := u.repository.GetUserByID(ctx, memberId, enterpriseId)

	if err != nil {
		return response.CreatedUserResponse{}, errorsBusiness.NewUserGettingDBError(memberId, err)
	}

	return response.CreatedUserResponse{
		ID:             user.ID,
		Email:          user.Email,
		ExternalUserID: user.ExternalUserID,
	}, nil
}

func (u *CreateUserUseDeunaAdapterImpl) ValidateUser(
	ctx context.Context,
	request request.CreateUserRequest,
	enterpriseId string,
) (response.UserValidatedResponse, error) {
	user, err := u.GetOrCreateUser(ctx, request.UserID, request.UserType, enterpriseId)

	if err != nil {
		return response.UserValidatedResponse{}, err
	}

	externalAuthReq := requestDeUna.DeunaAuthUserRequestDTO{
		Email: user.PaymentsExternalEmail,
	}

	externalAuthResponse, err := u.deunaAuth.AuthUser(
		ctx,
		externalAuthReq,
	)

	if err != nil {
		return response.UserValidatedResponse{}, errorsBusiness.NewUserFailCreateDeUnaError(request.UserID, err)
	}

	return response.UserValidatedResponse{
		ID:                user.ID,
		Email:             user.Email,
		ExternalUserID:    user.ExternalUserID,
		ExternalAuthToken: externalAuthResponse.AuthToken,
	}, nil
}

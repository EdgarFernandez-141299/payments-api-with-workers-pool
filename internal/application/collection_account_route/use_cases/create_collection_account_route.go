package usecases

import (
	"context"
	"errors"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/utils"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/response"
	"gorm.io/gorm"

	repositoriesCollectionAccount "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account"
	repositories "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account_route"

	errorBusiness "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/collection_center/errors"
)

type CollectionAccountRouteUsecaseIF interface {
	Create(
		ctx context.Context,
		collection request.CollectionAccountRouteRequest,
		enterpriseId string,
	) (response.CollectionAccountRouteResponse, error)
	Disable(ctx context.Context, id, enterpriseId string) (response.CollectionAccountRouteDisableResponse, error)
}

type CollectionAccountRouteUsecase struct {
	repository                  repositories.CollectionCenterAccountRouteRepositoryIF
	repositoryCollectionAccount repositoriesCollectionAccount.CollectionAccountRepositoryIF
}

func NewCollectionAccountRepositoryIF(
	repository repositories.CollectionCenterAccountRouteRepositoryIF,
	repositoryCollectionAccount repositoriesCollectionAccount.CollectionAccountRepositoryIF,
) CollectionAccountRouteUsecaseIF {
	return &CollectionAccountRouteUsecase{
		repository:                  repository,
		repositoryCollectionAccount: repositoryCollectionAccount,
	}
}

func (c *CollectionAccountRouteUsecase) Create(
	ctx context.Context,
	collection request.CollectionAccountRouteRequest,
	enterpriseId string,
) (response.CollectionAccountRouteResponse, error) {
	if !utils.IsValidCountryCode(collection.CountryCode) {
		return response.CollectionAccountRouteResponse{}, errorBusiness.NewInvalidCountryCodeError(
			errors.New("invalid country code"),
		)
	}

	if !utils.IsValidCurrencyCode(collection.CurrencyCode) {
		return response.CollectionAccountRouteResponse{}, errorBusiness.NewInvalidCountryCodeError(
			errors.New("invalid currency code"),
		)
	}

	if !collection.AssociatedOrigin.IsValid() {
		return response.CollectionAccountRouteResponse{}, errorBusiness.NewInvalidAssociatedOriginsError(
			errors.New("invalid associated origin"),
		)
	}

	routeFound, err := c.repository.FindRouteBy(
		ctx,
		collection.CountryCode,
		collection.CurrencyCode,
		collection.AssociatedOrigin.String(),
		enterpriseId,
	)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response.CollectionAccountRouteResponse{}, errorBusiness.NewCollectionAccountFindByError(err)
		}
	}

	if !routeFound.IsEmpty() {
		return response.CollectionAccountRouteResponse{},
			errorBusiness.NewCollectionAccountRouteAlreadyExist(
				errors.New("collection account route already exists"),
			)
	}

	collectionAccount, err := c.repositoryCollectionAccount.FindById(ctx, collection.CollectionAccountID, enterpriseId)
	if err != nil || collectionAccount.IsEmpty() {
		return response.CollectionAccountRouteResponse{}, errorBusiness.NewCollectionAccountNotFoundError(
			errors.New("collection account not found"),
		)
	}

	if collectionAccount.CurrencyCode != collection.CurrencyCode {
		return response.CollectionAccountRouteResponse{}, errorBusiness.NewCollectionAccountCurrencyCodeError(
			errors.New("currency code does not match"),
		)
	}

	entity := entities.NewCollectionAccountRouteEntity(
		collection,
		enterpriseId,
	)

	err = c.repository.Create(ctx, entity)
	if err != nil {
		return response.CollectionAccountRouteResponse{}, errorBusiness.NewCollectionAccountRouteCreateError(err)
	}

	return response.NewCollectionAccountRouteResponse(
		entity,
	), nil
}

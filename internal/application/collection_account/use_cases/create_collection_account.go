package usecases

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/utils"
	errorsBussines "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/collection_center/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account/dto/response"
	repository "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"
	repositoryCollection "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_center"
	"gorm.io/gorm"
)

type CollectionAccountUsecaseIF interface {
	Create(
		ctx context.Context,
		collectionaccount request.CollectionAccountRequest,
		enterpriseId string,
	) (response.CollectionAccountResponse, error)
}

type CollectionAccountUsecaseImpl struct {
	repository           repository.CollectionAccountRepositoryIF
	repoCollectionCenter repositoryCollection.CollectionCenterRepositoryIF
}

func NewCollectionAccountUsecase(
	repository repository.CollectionAccountRepositoryIF,
	repoCollectionCenter repositoryCollection.CollectionCenterRepositoryIF,
) CollectionAccountUsecaseIF {
	return &CollectionAccountUsecaseImpl{
		repository:           repository,
		repoCollectionCenter: repoCollectionCenter,
	}
}

func (p *CollectionAccountUsecaseImpl) Create(
	ctx context.Context,
	collectionAccount request.CollectionAccountRequest,
	enterpriseId string,
) (response.CollectionAccountResponse, error) {
	collectionCenter, err := p.repoCollectionCenter.FindByID(ctx, collectionAccount.CollectionCenterID, enterpriseId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.CollectionAccountResponse{}, errorsBussines.NewCollectionCenterNotFoundError(err)
		}

		return response.CollectionAccountResponse{}, err
	}

	if !collectionCenter.IsEmpty() {
		if !slices.Contains(collectionCenter.AvailableCurrencies, collectionAccount.CurrencyCode) {
			return response.CollectionAccountResponse{}, errorsBussines.NewInvalidCurrencyCodeError(
				errors.New("currency code not available"),
			)
		}
	}

	accountFound, err := p.repository.FindByAccountNumber(ctx, collectionAccount.AccountNumber, enterpriseId)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response.CollectionAccountResponse{}, errorsBussines.NewCollectionAccountAlreadyExistsError(
				err,
			)
		}
	}

	if !accountFound.IsEmpty() {
		return response.CollectionAccountResponse{}, errorsBussines.NewCollectionAccountAlreadyExistsError(
			errors.New("account number already exists"),
		)
	}

	if !utils.IsValidCurrencyCode(collectionAccount.CurrencyCode) {
		return response.CollectionAccountResponse{}, errorsBussines.NewInvalidCurrencyCodeError(
			fmt.Errorf("%s is not a valid currency code", collectionAccount.CurrencyCode),
		)
	}

	if !collectionAccount.AccountType.IsValid() {
		return response.CollectionAccountResponse{}, errorsBussines.NewInvalidAccountTypeError(
			errors.New("invalid account type"),
		)
	}

	entity := entities.NewCollectionAccountEntity(collectionAccount, enterpriseId)
	err = p.repository.Create(ctx, entity)

	if err != nil {
		return response.CollectionAccountResponse{}, errorsBussines.NewCollectionAccountError(err)
	}

	return response.NewCollectionAccountResponse(entity), nil
}

package queries

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/collection_account/repositories"
	entitiesDomain "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type GetCollectionAccountByRouteUsecaseIF interface {
	GetCollectionAccountByRoute(
		ctx context.Context,
		country, associatedOrigin,
		currency, enterpriseId string,
	) (entitiesDomain.CollectionAccount, error)
}

type GetCollectionAccountByRouteUsecase struct {
	repository repositories.CollectionAccountReadRepositoryIF
}

func NewGetCollectionAccountByRouteUsecase(
	repository repositories.CollectionAccountReadRepositoryIF,
) GetCollectionAccountByRouteUsecaseIF {
	return &GetCollectionAccountByRouteUsecase{
		repository: repository,
	}
}

func (p *GetCollectionAccountByRouteUsecase) GetCollectionAccountByRoute(
	ctx context.Context, country, associatedOrigin, currency, enterpriseId string,
) (entitiesDomain.CollectionAccount, error) {
	collectionAccount, err := p.repository.GetCollectionAccountRoute(ctx, country, associatedOrigin, currency, enterpriseId)

	if err != nil {
		return entitiesDomain.CollectionAccount{}, errors.NewCollectionAccountNotFound()
	}

	if collectionAccount.IsEmpty() {
		return entitiesDomain.CollectionAccount{}, errors.NewCollectionAccountNotFound()
	}

	currencyCode, err := value_objects.NewCurrencyCode(collectionAccount.CurrencyCode)
	if err != nil {
		return entitiesDomain.CollectionAccount{}, err
	}

	return entitiesDomain.CollectionAccount{
		ID:                     collectionAccount.ID.String(),
		AccountType:            collectionAccount.AccountType,
		CollectionCenterID:     collectionAccount.CollectionCenterID,
		CurrencyCode:           currencyCode,
		AccountNumber:          collectionAccount.AccountNumber,
		BankName:               collectionAccount.BankName,
		InterbankAccountNumber: collectionAccount.InterbankAccountNumber,
		EnterpriseID:           collectionAccount.EnterpriseID,
	}, nil
}

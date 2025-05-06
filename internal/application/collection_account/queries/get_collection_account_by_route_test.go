package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/gommon/uid"
	entitiesDomain "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repositories"
)

var (
	ctx                     = context.TODO()
	country          string = "MX"
	associatedOrigin string = enums.Downpayment.String()
	currency, _             = value_objects.NewCurrencyCode("MXN")
	enterpriseId     string = "enterpriseId"
	id, _                   = uid.NewUniqueID(uid.WithID("id"))
)

func TestQueryCollectionAccountByRoute(t *testing.T) {
	t.Run("should return an error getting collection account", func(t *testing.T) {
		repositoryMock := repositories.NewCollectionAccountReadRepositoryIF(t)

		repositoryMock.On("GetCollectionAccountRoute", ctx, country, associatedOrigin, currency.Code, enterpriseId).
			Return(entities.CollectionAccountEntity{}, errors.New("some error"))

		collectionAccount := NewGetCollectionAccountByRouteUsecase(repositoryMock)

		response, err := collectionAccount.GetCollectionAccountByRoute(ctx, country, associatedOrigin, currency.Code, enterpriseId)

		assert.EqualError(t, err, "Business Error code: COLLECTION_ACCOUNT_NOT_FOUND, message: collection account not found")
		assert.Zero(t, response)
	})

	t.Run("should return an error getting collection account, collection empty", func(t *testing.T) {
		repositoryMock := repositories.NewCollectionAccountReadRepositoryIF(t)

		repositoryMock.On("GetCollectionAccountRoute", ctx, country, associatedOrigin, currency.Code, enterpriseId).
			Return(entities.CollectionAccountEntity{}, nil)

		collectionAccount := NewGetCollectionAccountByRouteUsecase(repositoryMock)

		response, err := collectionAccount.GetCollectionAccountByRoute(ctx, country, associatedOrigin, currency.Code, enterpriseId)

		assert.EqualError(t, err, "Business Error code: COLLECTION_ACCOUNT_NOT_FOUND, message: collection account not found")
		assert.Zero(t, response)
	})

	t.Run("should return an error getting collection account, collection empty", func(t *testing.T) {
		repositoryMock := repositories.NewCollectionAccountReadRepositoryIF(t)

		repositoryMock.On("GetCollectionAccountRoute", ctx, country, associatedOrigin, currency.Code, enterpriseId).
			Return(entities.CollectionAccountEntity{
				ID:                     id,
				AccountType:            enums.Payers.String(),
				CollectionCenterID:     "collectionCenterID",
				CurrencyCode:           currency.Code,
				AccountNumber:          "accountNumber",
				BankName:               "bankName",
				InterbankAccountNumber: "interbankAccountNumber",
				EnterpriseID:           enterpriseId,
			}, nil)

		collectionAccount := NewGetCollectionAccountByRouteUsecase(repositoryMock)

		response, err := collectionAccount.GetCollectionAccountByRoute(ctx, country, associatedOrigin, currency.Code, enterpriseId)

		assert.Nil(t, err)

		assert.Equal(t, entitiesDomain.CollectionAccount{
			ID:                     id.String(),
			AccountType:            enums.Payers.String(),
			CollectionCenterID:     "collectionCenterID",
			CurrencyCode:           currency,
			AccountNumber:          "accountNumber",
			BankName:               "bankName",
			InterbankAccountNumber: "interbankAccountNumber",
			EnterpriseID:           enterpriseId,
		},
			response,
		)
	})
}

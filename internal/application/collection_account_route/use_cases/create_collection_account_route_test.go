package usecases

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	entitiesRoute "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"
	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repositories"
)

var (
	ctx   = context.TODO()
	id, _ = uid.NewUniqueID(uid.WithID("id"))
)

func TestCreate(t *testing.T) {
	t.Run("should return invalid country", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)

		mockColeectionAccount := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "INVALID",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    enums.Downpayment,
		}
		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockColeectionAccount,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "INVALID_COUNTRY_CODE"))
	})

	t.Run("should return invalid currency", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "INVALID",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    "WRONG",
		}

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(context.Background(), req, "enterprise1")

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "INVALID_COUNTRY_CODE"))
	})

	t.Run("should return associated origin wrong", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    "INVALID",
		}

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")
		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "INVALID_ASSOCIATED_ORIGINS"))
	})

	t.Run("should return error finding route", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    enums.Downpayment,
		}

		mockCollectionCenterRepo.On("FindRouteBy", ctx, "US", "USD", "DOWNPAYMENT", "enterprise1").
			Return(entitiesRoute.CollectionAccountRouteEntity{}, errors.New("error find collection account route"))

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")
		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "COLLECTION_ACCOUNT_ERROR_DB"))
	})

	t.Run("should return error route already exist", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    enums.Downpayment,
		}

		mockCollectionCenterRepo.On("FindRouteBy", ctx, "US", "USD", "DOWNPAYMENT", "enterprise1").
			Return(entitiesRoute.CollectionAccountRouteEntity{
				ID: id,
			}, nil)

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")
		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "COLLECTION_ACCOUNT_ERROR_DB"))
	})

	t.Run("should return error account doesn't exist", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    enums.Downpayment,
		}

		mockCollectionCenterRepo.On("FindRouteBy", ctx, "US", "USD", "DOWNPAYMENT", "enterprise1").
			Return(entitiesRoute.CollectionAccountRouteEntity{}, nil)

		mockCollectionAccountRepo.On("FindById", ctx, "CA001", "enterprise1").
			Return(entities.CollectionAccountEntity{}, errors.New("not found"))

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")
		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "COLLECTION_ACCOUNT_ERROR_DB"))
	})

	t.Run("should return currency not match error", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    enums.Downpayment,
		}

		mockCollectionCenterRepo.On("FindRouteBy", ctx, "US", "USD", "DOWNPAYMENT", "enterprise1").
			Return(entitiesRoute.CollectionAccountRouteEntity{}, nil)

		mockCollectionAccountRepo.On("FindById", ctx, "CA001", "enterprise1").
			Return(entities.CollectionAccountEntity{
				ID:           id,
				CurrencyCode: "EUR",
			}, nil)

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")
		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "CURRENCY_NO_MATCH"))
	})

	t.Run("should return an error creating the collection account route", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    enums.Downpayment,
		}

		mockCollectionCenterRepo.On("FindRouteBy", ctx, "US", "USD", "DOWNPAYMENT", "enterprise1").Return(entitiesRoute.CollectionAccountRouteEntity{}, nil)
		mockCollectionAccountRepo.On("FindById", ctx, "CA001", "enterprise1").Return(entities.CollectionAccountEntity{
			ID:           id,
			CurrencyCode: "USD",
		}, nil)

		mockCollectionCenterRepo.On("Create", ctx, mock.Anything).Return(errors.New("error creating collection account route"))

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")
		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "COLLECTION_ACCOUNT_ERROR_DB"))
	})

	t.Run("should create a route successful", func(t *testing.T) {
		mockCollectionCenterRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
		mockCollectionAccountRepo := mockRepository.NewCollectionAccountRepositoryIF(t)

		req := request.CollectionAccountRouteRequest{
			CountryCode:         "US",
			CurrencyCode:        "USD",
			CollectionAccountID: "CA001",
			AssociatedOrigin:    enums.Downpayment,
		}

		mockCollectionCenterRepo.On("FindRouteBy", ctx, "US", "USD", "DOWNPAYMENT", "enterprise1").Return(entitiesRoute.CollectionAccountRouteEntity{}, nil)
		mockCollectionAccountRepo.On("FindById", ctx, "CA001", "enterprise1").Return(entities.CollectionAccountEntity{
			ID:           id,
			CurrencyCode: "USD",
		}, nil)

		mockCollectionCenterRepo.On("Create", ctx, mock.Anything).Return(
			nil,
		)

		usecase := NewCollectionAccountRepositoryIF(
			mockCollectionCenterRepo,
			mockCollectionAccountRepo,
		)

		_, err := usecase.Create(ctx, req, "enterprise1")
		assert.Nil(t, err)
	})
}

package usecases

import (
	"context"
	"errors"
	"testing"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"

	"gitlab.com/clubhub.ai1/gommon/uid"
	errorsBussines "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/collection_center/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/collection_account/entities"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	entitiesCollectionCenter "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account/dto/request"

	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repositories"
)

var (
	accountType            enums.AccountType = enums.Mixed
	ctx                                      = context.TODO()
	collectionCenterId                       = "SA55xX7aH8PsmN"
	enterpriseId                             = "enterpriseId"
	currencyCode                             = "USD"
	currencyCodeWrong                        = "USDX"
	accountNumber                            = "12345678"
	banckName                                = "bancolombia"
	branchCode                               = "1234"
	InterbankAccountNumber                   = "19248394"
	CollectionCenterName                     = "Collection mix"
	UID, _                                   = uid.NewUniqueID(uid.WithID("id"))
)

func TestCreate(t *testing.T) {
	t.Run("should return an error creating when account not exist collection account", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:        accountType,
			CollectionCenterID: collectionCenterId,
			CurrencyCode:       currencyCode,
			AccountNumber:      accountNumber,
			BankName:           banckName,
			BranchCode:         branchCode,
		}

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{}, gorm.ErrRecordNotFound)

		paymentUseCaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := paymentUseCaseMock.Create(ctx, requestMock, enterpriseId)

		assert.NotNil(t, err)
		assert.EqualError(t, err, errorsBussines.NewCollectionCenterNotFoundError(gorm.ErrRecordNotFound).Error())

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should return an error creating collection account", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            accountType,
			CollectionCenterID:     collectionCenterId,
			CurrencyCode:           currencyCode,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		repoErr := errors.New("error getting collection account")

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{}, repoErr)

		paymentUseCaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := paymentUseCaseMock.Create(ctx, requestMock, enterpriseId)

		assert.NotNil(t, err)
		assert.EqualError(t, err, repoErr.Error())

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should return an error, currency code not availabe", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            accountType,
			CollectionCenterID:     collectionCenterId,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{
				Name:                CollectionCenterName,
				ID:                  UID,
				Description:         CollectionCenterName,
				EnterpriseID:        enterpriseId,
				AvailableCurrencies: []string{currencyCode},
			}, nil)

		paymentUseCaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := paymentUseCaseMock.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "INVALID_CURRENCY_CODE"))

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should return an error creating collection account, error number", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            accountType,
			CollectionCenterID:     collectionCenterId,
			CurrencyCode:           currencyCode,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		repoErr := errors.New("error getting collection account")

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{
				Name:                CollectionCenterName,
				ID:                  UID,
				Description:         CollectionCenterName,
				EnterpriseID:        enterpriseId,
				AvailableCurrencies: []string{currencyCode},
			}, nil)

		repositoryMock.On("FindByAccountNumber", ctx, accountNumber, enterpriseId).
			Return(entities.CollectionAccountEntity{}, repoErr)

		paymentUseCaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := paymentUseCaseMock.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "COLLECTION_ACCOUNT_ERROR_DB"))

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should return an error creating collection account, number account already exist", func(t *testing.T) {
		// Arrange
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            accountType,
			CollectionCenterID:     collectionCenterId,
			CurrencyCode:           currencyCode,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		// repoErr := errors.New("error getting collection account")

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{
				Name:                CollectionCenterName,
				ID:                  UID,
				Description:         CollectionCenterName,
				EnterpriseID:        enterpriseId,
				AvailableCurrencies: []string{currencyCode},
			}, nil)

		repositoryMock.On("FindByAccountNumber", ctx, accountNumber, enterpriseId).
			Return(entities.CollectionAccountEntity{
				ID:            UID,
				AccountType:   accountType.String(),
				CurrencyCode:  currencyCode,
				AccountNumber: accountNumber,
			}, nil)

		paymentUseCaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := paymentUseCaseMock.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "COLLECTION_ACCOUNT_ERROR_DB"))

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should return an error, invalid currency", func(t *testing.T) {
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            accountType,
			CollectionCenterID:     collectionCenterId,
			CurrencyCode:           currencyCodeWrong,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{
				Name:                CollectionCenterName,
				ID:                  UID,
				Description:         CollectionCenterName,
				EnterpriseID:        enterpriseId,
				AvailableCurrencies: []string{currencyCodeWrong},
			}, nil)

		repositoryMock.On("FindByAccountNumber", ctx, accountNumber, enterpriseId).
			Return(entities.CollectionAccountEntity{}, nil)

		collectionAccountUseCaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := collectionAccountUseCaseMock.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "INVALID_CURRENCY_CODE"))

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should return an error creating collection account", func(t *testing.T) {
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            accountType,
			CollectionCenterID:     collectionCenterId,
			CurrencyCode:           currencyCode,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{
				Name:                CollectionCenterName,
				ID:                  UID,
				Description:         CollectionCenterName,
				EnterpriseID:        enterpriseId,
				AvailableCurrencies: []string{currencyCode},
			}, nil)

		repositoryMock.On("FindByAccountNumber", ctx, accountNumber, enterpriseId).
			Return(entities.CollectionAccountEntity{}, nil)

		repositoryMock.On("Create", ctx, mock.Anything).
			Return(errors.New("error creating collection account"))

		collectionaccountUsecaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := collectionaccountUsecaseMock.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "COLLECTION_ACCOUNT_ERROR_DB"))

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should return error validate type account", func(t *testing.T) {
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            "accountType",
			CollectionCenterID:     collectionCenterId,
			CurrencyCode:           currencyCode,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{
				Name:                CollectionCenterName,
				ID:                  UID,
				Description:         CollectionCenterName,
				EnterpriseID:        enterpriseId,
				AvailableCurrencies: []string{currencyCode},
			}, nil)

		repositoryMock.On("FindByAccountNumber", ctx, accountNumber, enterpriseId).
			Return(entities.CollectionAccountEntity{}, nil)

		collectionaccountUsecaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := collectionaccountUsecaseMock.Create(ctx, requestMock, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "INVALID_ACCOUNT_TYPE"))

		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})

	t.Run("should create a collection account succesfull", func(t *testing.T) {
		repositoryMock := mockRepository.NewCollectionAccountRepositoryIF(t)
		repositoryMockCollectionCenter := mockRepository.NewCollectionCenterRepositoryIF(t)

		requestMock := request.CollectionAccountRequest{
			AccountType:            accountType,
			CollectionCenterID:     collectionCenterId,
			CurrencyCode:           currencyCode,
			AccountNumber:          accountNumber,
			BankName:               banckName,
			BranchCode:             branchCode,
			InterbankAccountNumber: InterbankAccountNumber,
		}

		repositoryMockCollectionCenter.On("FindByID", ctx, collectionCenterId, enterpriseId).
			Return(entitiesCollectionCenter.CollectionCenterEntity{
				Name:                CollectionCenterName,
				ID:                  UID,
				Description:         CollectionCenterName,
				EnterpriseID:        enterpriseId,
				AvailableCurrencies: []string{currencyCode},
			}, nil)

		repositoryMock.On("FindByAccountNumber", ctx, accountNumber, enterpriseId).
			Return(entities.CollectionAccountEntity{}, nil)

		repositoryMock.On("Create", ctx, mock.Anything).
			Return(nil)

		collectionaccountUsecaseMock := NewCollectionAccountUsecase(repositoryMock, repositoryMockCollectionCenter)
		_, err := collectionaccountUsecaseMock.Create(ctx, requestMock, enterpriseId)

		assert.Nil(t, err)
		repositoryMock.AssertExpectations(t)
		repositoryMockCollectionCenter.AssertExpectations(t)
	})
}

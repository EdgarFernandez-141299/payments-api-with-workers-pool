package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_center/dto/request"
	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repositories"
)

var (
	ctx = context.TODO()
)

func TestCreateCollectionCenter(t *testing.T) {
	t.Run("error create collection center", func(t *testing.T) {
		mockRepository := new(mockRepository.CollectionCenterRepositoryIF)
		useCase := NewCollectionCenterUsecase(mockRepository)

		mockRepository.On("Create", ctx, mock.Anything).Return(errors.New("error creating collection center"))

		mockCollectionCenter := &request.CollectionCenterRequest{
			Name:                "Test Center",
			Description:         "Test Description",
			AvailableCurrencies: []string{"COP"},
		}

		_, err := useCase.Create(ctx, *mockCollectionCenter, "1")
		assert.NotNil(t, err)
		assert.EqualError(t, err, "error creating collection center")
		mockRepository.AssertExpectations(t)
	})

	t.Run("error create collection center currency_code invalid", func(t *testing.T) {
		mockRepository := new(mockRepository.CollectionCenterRepositoryIF)
		useCase := NewCollectionCenterUsecase(mockRepository)

		mockCollectionCenter := &request.CollectionCenterRequest{
			Name:                "Test Center",
			Description:         "Test Description",
			AvailableCurrencies: []string{"OPP"},
		}

		_, err := useCase.Create(ctx, *mockCollectionCenter, "1")
		assert.NotNil(t, err)
		assert.EqualError(t, err, "currency OPP is not valid")
		mockRepository.AssertExpectations(t)
	})

	t.Run("error create collection center", func(t *testing.T) {
		mockRepository := new(mockRepository.CollectionCenterRepositoryIF)
		useCase := NewCollectionCenterUsecase(mockRepository)

		mockRepository.On("Create", ctx, mock.Anything).Return(nil)

		mockCollectionCenter := &request.CollectionCenterRequest{
			Name:                "Test Center",
			Description:         "Test Description",
			AvailableCurrencies: []string{"COP"},
		}

		_, err := useCase.Create(ctx, *mockCollectionCenter, "1")
		assert.Nil(t, err)
		mockRepository.AssertExpectations(t)
	})
}

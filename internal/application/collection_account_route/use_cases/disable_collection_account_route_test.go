package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/collection_account_route/dto/response"
	mockRepository "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repositories"
)

func TestDisableCollectionAccountRoute(t *testing.T) {
	ctx := context.TODO()

	mockRepo := mockRepository.NewCollectionCenterAccountRouteRepositoryIF(t)
	usecase := &CollectionAccountRouteUsecase{
		repository: mockRepo,
	}

	t.Run("should disable successfully", func(t *testing.T) {
		mockRepo.On("Disable", ctx, "test-id", "test-enterprise-id").Return(nil)

		resp, err := usecase.Disable(ctx, "test-id", "test-enterprise-id")

		assert.NoError(t, err)
		assert.Equal(t, "test-id", resp.ID)
		assert.Equal(t, "disabled", resp.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return a repository error", func(t *testing.T) {
		mockRepo.On("Disable", mock.Anything, "test-id", "test-enterprise-id").Return(errors.New("repository error"))

		resp, err := usecase.Disable(context.Background(), "test-id", "test-enterprise-id")

		assert.Error(t, err)
		assert.Equal(t, response.CollectionAccountRouteDisableResponse{}, resp)

		mockRepo.AssertExpectations(t)
	})
}

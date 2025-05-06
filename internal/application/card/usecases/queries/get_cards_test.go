package queries

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
	userEntities "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/user/entities"
	"gorm.io/gorm"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

var (
	ctx = context.TODO()
)

func TestGetCardsByUserID(t *testing.T) {
	t.Run("should return an error getting cards", func(t *testing.T) {
		mockRepo := repository.NewCardReadRepositoryIF(t)
		mockRepoUser := repository.NewUserReadRepositoryIF(t)
		usecase := NewGetCardUsecase(mockRepo, mockRepoUser)

		userId := "test-user-id"
		enterpriseId := "test-enterprise-id"

		mockRepoUser.On("GetUserByID", ctx, userId, enterpriseId).Return(userEntities.UserEntity{}, nil)

		mockRepo.On("GetCardsByUserID", ctx, userId, enterpriseId).Return(entities.CardEntities{}, errors.New("syntax error ,/9"))

		result, err := usecase.GetCardsByUserID(ctx, userId, enterpriseId)

		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "CARD_ERROR_GET_DB"))
		assert.Equal(t, 0, len(result))

		mockRepoUser.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return a list of cards", func(t *testing.T) {
		mockRepo := repository.NewCardReadRepositoryIF(t)
		mockRepoUser := repository.NewUserReadRepositoryIF(t)
		usecase := NewGetCardUsecase(mockRepo, mockRepoUser)

		userId := "test-user-id"
		enterpriseId := "test-enterprise-id"

		id, _ := uid.NewUniqueID(uid.WithID("userId"))

		cards := entities.CardEntities{
			{
				ID:             id,
				UserID:         userId,
				ExternalCardID: "card-1",
			},
		}

		mockRepoUser.On("GetUserByID", ctx, userId, enterpriseId).Return(userEntities.UserEntity{}, nil)

		mockRepo.On("GetCardsByUserID", ctx, userId, enterpriseId).Return(cards, nil)

		result, err := usecase.GetCardsByUserID(ctx, userId, enterpriseId)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(result))
	})

	t.Run("should return an error getting user", func(t *testing.T) {
		mockRepo := repository.NewCardReadRepositoryIF(t)
		mockRepoUser := repository.NewUserReadRepositoryIF(t)
		usecase := NewGetCardUsecase(mockRepo, mockRepoUser)

		userId := "test-user-id"
		enterpriseId := "test-enterprise-id"

		mockRepoUser.On("GetUserByID", ctx, userId, enterpriseId).
			Return(userEntities.UserEntity{}, errors.New("USER_NOT_FOUND"))

		result, err := usecase.GetCardsByUserID(ctx, userId, enterpriseId)

		assert.Error(t, err)
		assert.EqualError(t, err, "USER_NOT_FOUND")
		assert.Equal(t, 0, len(result))
	})

	t.Run("should return an error record not found", func(t *testing.T) {
		mockRepo := repository.NewCardReadRepositoryIF(t)
		mockRepoUser := repository.NewUserReadRepositoryIF(t)
		usecase := NewGetCardUsecase(mockRepo, mockRepoUser)

		userId := "test-user-id"
		enterpriseId := "test-enterprise-id"

		mockRepoUser.On("GetUserByID", ctx, userId, enterpriseId).Return(userEntities.UserEntity{}, gorm.ErrRecordNotFound)

		result, err := usecase.GetCardsByUserID(ctx, userId, enterpriseId)

		assert.Error(t, err)
		assert.EqualError(t, err, "Business Error code: USER_NOT_FOUND, message: user not found with memberId: test-user-id")
		assert.Equal(t, 0, len(result))
	})
}

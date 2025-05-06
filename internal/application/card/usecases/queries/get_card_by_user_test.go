package queries

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

var (
	userID       = "user-id"
	cardID       = "card-id"
	enterpriseID = "enterprise-id"
	id, _        = uid.NewUniqueID(uid.WithID("id"))
)

func TestQueryCardByUser(t *testing.T) {
	t.Run("error getting card", func(t *testing.T) {
		mockReadRepo := repository.NewCardReadRepositoryIF(t)
		mockReadRepo.On("GetCardByUserID", mock.Anything, userID, cardID, enterpriseID).Return(entities.CardEntity{}, errors.New("database error"))
		usecase := NewGetCardByUserUsecase(mockReadRepo)

		// Execute the use case
		result, err := usecase.GetCardByIDAndUserID(ctx, userID, cardID, enterpriseID)

		// Assert the results
		assert.EqualError(t, err, "Business Error code: CARD_IS_NOT_FROM_MEMBER, message: Card is not from member user-id ")
		assert.Zero(t, result)

		mockReadRepo.AssertExpectations(t)
	})

	t.Run("get card successful", func(t *testing.T) {
		mockReadRepo := repository.NewCardReadRepositoryIF(t)

		mockReadRepo.On("GetCardByUserID", mock.Anything, userID, cardID, enterpriseID).Return(entities.CardEntity{
			ID:             id,
			ExternalCardID: "external-card-id",
			LastFour:       "1234",
			Brand:          "visa",
			CardHolder:     "card-holder",
		}, nil)
		usecase := NewGetCardByUserUsecase(mockReadRepo)

		// Execute the use case
		result, err := usecase.GetCardByIDAndUserID(ctx, userID, cardID, enterpriseID)

		// Assert the results
		assert.Nil(t, err)
		assert.NotZero(t, result)

		mockReadRepo.AssertExpectations(t)
	})

}

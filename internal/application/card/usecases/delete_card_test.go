package usecases

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/gommon/uid"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/card/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/services"
	"gorm.io/gorm"
)

var (
	enterpriseID = "enterpriseID"
	cardID       = "cardID"
	id, _        = uid.NewUniqueID(uid.WithID("id"))
)

func TestDeleteCard(t *testing.T) {
	t.Run("should return an error getting card and email by user, record not found", func(t *testing.T) {
		cardAdapterMock := adapters.NewCardAdapter(t)
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		readRepository := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.DeleteCardRequest{
			CardID: cardID,
			UserID: "userID",
		}

		readRepository.On("GetCardAndUserEmailByUserID", ctx, requestMock.UserID, requestMock.CardID, enterpriseID).
			Return(nil, gorm.ErrRecordNotFound)

		usecase := NewDeleteCardUseCase(
			cardAdapterMock,
			repositoryMock,
			readRepository,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := usecase.DeleteCard(ctx, requestMock, enterpriseID, "en")

		assert.EqualError(t, err, "Business Error code: CARD_NOT_FOUND, message: Card not found [ cardID ]")
	})

	t.Run("should return an error getting card by user, error db", func(t *testing.T) {
		cardAdapterMock := adapters.NewCardAdapter(t)
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		readRepository := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.DeleteCardRequest{
			CardID: cardID,
			UserID: "userID",
		}

		readRepository.On("GetCardAndUserEmailByUserID", ctx, requestMock.UserID, requestMock.CardID, enterpriseID).
			Return(nil, errors.New("database error"))

		usecase := NewDeleteCardUseCase(
			cardAdapterMock,
			repositoryMock,
			readRepository,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := usecase.DeleteCard(ctx, requestMock, enterpriseID, "en")

		assert.EqualError(t, err, "database error")
	})

	t.Run("should return an error deleting an card", func(t *testing.T) {
		cardAdapterMock := adapters.NewCardAdapter(t)
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		readRepository := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.DeleteCardRequest{
			CardID: cardID,
			UserID: "userID",
		}

		readRepository.On("GetCardAndUserEmailByUserID", ctx, requestMock.UserID, requestMock.CardID, enterpriseID).
			Return(&projections.CardUserEmailProjection{
				ID:             "id",
				UserID:         "userID",
				ExternalCardID: "external-card-id",
				Email:          "email@gmail.com",
			}, nil)

		cardAdapterMock.On("DeleteCard", ctx, id.String(), "external-card-id", "userID", enterpriseID).
			Return(response.DeleteCardResponseDTO{
				Code:    "code",
				Message: "message",
			}, assert.AnError)

		usecase := NewDeleteCardUseCase(
			cardAdapterMock,
			repositoryMock,
			readRepository,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := usecase.DeleteCard(ctx, requestMock, enterpriseID, "en")

		assert.EqualError(t, err, "Business Error code: DELETE_CARD_ERROR, message: error deleting card DEUNA service")
	})

	t.Run("should successfully delete a card", func(t *testing.T) {
		cardAdapterMock := adapters.NewCardAdapter(t)
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		readRepository := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.DeleteCardRequest{
			CardID: cardID,
			UserID: "userID",
		}

		readRepository.On("GetCardAndUserEmailByUserID", ctx, requestMock.UserID, requestMock.CardID, enterpriseID).
			Return(&projections.CardUserEmailProjection{
				ID:             "id",
				UserID:         "userID",
				ExternalCardID: "external-card-id",
				Email:          "email@gmail.com",
			}, nil)

		cardAdapterMock.On("DeleteCard", ctx, id.String(), "external-card-id", "userID", enterpriseID).
			Return(response.DeleteCardResponseDTO{
				Code:    "200",
				Message: "",
			}, nil)

		repositoryMock.On("DeleteCard", ctx, requestMock.CardID).
			Return(nil)

		notificationServiceMock.On(
			"NotifyCardDeletion",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil)

		usecase := NewDeleteCardUseCase(
			cardAdapterMock,
			repositoryMock,
			readRepository,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		resp, err := usecase.DeleteCard(ctx, requestMock, enterpriseID, "en")

		assert.NoError(t, err)
		assert.Equal(t, "Success", resp.Status)
		assert.Equal(t, "Card successfully deleted.", resp.Message)
	})

	t.Run("should return an error deleting card from database", func(t *testing.T) {
		cardAdapterMock := adapters.NewCardAdapter(t)
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		readRepository := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.DeleteCardRequest{
			CardID: cardID,
			UserID: "userID",
		}

		readRepository.On("GetCardAndUserEmailByUserID", ctx, requestMock.UserID, requestMock.CardID, enterpriseID).
			Return(&projections.CardUserEmailProjection{
				ID:             "id",
				UserID:         "userID",
				ExternalCardID: "external-card-id",
				Email:          "email@gmail.com",
			}, nil)

		cardAdapterMock.On("DeleteCard", ctx, id.String(), "external-card-id", "userID", enterpriseID).
			Return(response.DeleteCardResponseDTO{
				Code:    "200",
				Message: "",
			}, nil)

		repositoryMock.On("DeleteCard", ctx, requestMock.CardID).
			Return(assert.AnError)

		usecase := NewDeleteCardUseCase(
			cardAdapterMock,
			repositoryMock,
			readRepository,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := usecase.DeleteCard(ctx, requestMock, enterpriseID, "en")

		assert.EqualError(t, err, "Business Error code: DELETE_CARD_DB_ERROR, message: error deleting card db")
	})
}

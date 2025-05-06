package usecases

import (
	"context"
	"testing"
	"time"

	"gitlab.com/clubhub.ai1/go-libraries/observability/errorx/pkg/domain"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/api/handler/card/dto/request"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/infra/repositories/card/entities"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
)

var (
	ctx          = context.Background()
	userID       = "userId"
	enterpriseId = "enterpriseId"
	cardId       = "1111"
	alias        = "alias"
	status       = "actived"
	cardType     = "credit"
	cardBrand    = "visa"
	lastFour     = "1111"
	firstSix     = "111111"
	expiration   = "01/06"
)

func TestCreateCard(t *testing.T) {
	t.Run("should fail when create card adapter fails", func(t *testing.T) {
		// Arrange
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		cardAdapter := adapters.NewCardAdapter(t)
		cardReadRepositoryMock := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.CardRequest{
			UserID:         userID,
			CardId:         cardId,
			Alias:          alias,
			Status:         status,
			CardType:       cardType,
			CardBrand:      cardBrand,
			LastFour:       lastFour,
			FirstSix:       firstSix,
			ExpirationDate: expiration,
			IsDefault:      false,
			IsRecurrent:    false,
		}

		// Act
		repositoryMock.On(
			"CreateCard", mock.Anything, mock.MatchedBy(func(card *entities.CardEntity) bool {
				return card.UserID == requestMock.UserID &&
					card.ExternalCardID == requestMock.CardId &&
					card.Alias == requestMock.Alias &&
					card.Status == requestMock.Status &&
					card.CardType == requestMock.CardType &&
					card.Brand == requestMock.CardBrand &&
					card.LastFour == requestMock.LastFour
			})).
			Return(assert.AnError)

		useCase := NewCardUseCase(
			cardAdapter,
			repositoryMock,
			cardReadRepositoryMock,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := useCase.CreateCard(ctx, requestMock, enterpriseId, "en")

		// Assert
		assert.Error(t, err)
		assert.True(t, domain.IsBusinessErrorCode(err, "CREATE_CARD_DB_ERROR"))
	})

	t.Run("should create a card successful", func(t *testing.T) {
		// Arrange
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		cardAdapter := adapters.NewCardAdapter(t)
		cardReadRepositoryMock := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.CardRequest{
			UserID:         userID,
			CardId:         cardId,
			Alias:          alias,
			Status:         status,
			CardType:       cardType,
			CardBrand:      cardBrand,
			LastFour:       lastFour,
			FirstSix:       firstSix,
			ExpirationDate: expiration,
			IsDefault:      false,
			IsRecurrent:    false,
		}

		// Act
		repositoryMock.On(
			"CreateCard", mock.Anything, mock.MatchedBy(func(card *entities.CardEntity) bool {
				return card.UserID == requestMock.UserID &&
					card.ExternalCardID == requestMock.CardId &&
					card.Alias == requestMock.Alias &&
					card.Status == requestMock.Status
			})).
			Return(nil)

		userReadRepositoryMock.On(
			"GetEmailByUserID", mock.Anything, mock.Anything, mock.Anything).Return("test@test.com", nil)

		notificationServiceMock.On(
			"NotifyCardAddition",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil)

		useCase := NewCardUseCase(
			cardAdapter,
			repositoryMock,
			cardReadRepositoryMock,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := useCase.CreateCard(ctx, requestMock, enterpriseId, "en")

		// Assert
		assert.Nil(t, err)
	})

	t.Run("should create a card successful when get email by user id fails", func(t *testing.T) {
		// Arrange
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		cardAdapter := adapters.NewCardAdapter(t)
		cardReadRepositoryMock := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.CardRequest{
			UserID:         userID,
			CardId:         cardId,
			Alias:          alias,
			Status:         status,
			CardType:       cardType,
			CardBrand:      cardBrand,
			LastFour:       lastFour,
			FirstSix:       firstSix,
			ExpirationDate: expiration,
			IsDefault:      false,
			IsRecurrent:    false,
		}

		repositoryMock.On(
			"CreateCard", mock.Anything, mock.MatchedBy(func(card *entities.CardEntity) bool {
				return card.UserID == requestMock.UserID &&
					card.ExternalCardID == requestMock.CardId &&
					card.Alias == requestMock.Alias &&
					card.Status == requestMock.Status
			})).
			Return(nil)

		userReadRepositoryMock.On(
			"GetEmailByUserID", mock.Anything, mock.Anything, mock.Anything).Return("", assert.AnError)

		logMock.On(
			"Warn", mock.Anything, mock.Anything).Return(nil)

		notificationServiceMock.On(
			"NotifyCardAddition",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil)

		useCase := NewCardUseCase(
			cardAdapter,
			repositoryMock,
			cardReadRepositoryMock,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := useCase.CreateCard(ctx, requestMock, enterpriseId, "en")

		// Assert
		assert.Nil(t, err)
	})
}

func TestTriggerCardExpiringSoonNotifications(t *testing.T) {

	notificationCardExpiringSoonProjections := []projections.NotificationCardExpiringSoonProjection{
		{
			UserID:         userID,
			LastFour:       lastFour,
			Email:          "test@test.com",
			ExpirationDate: time.Now().AddDate(0, 1, 0),
			EnterpriseID:   enterpriseId,
		},
	}

	t.Run("should fail when get cards expiring soon fails", func(t *testing.T) {
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		cardAdapter := adapters.NewCardAdapter(t)
		cardReadRepositoryMock := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.NotificationCardExpiringSoonRequestDTO{
			NotificationChannels: []constants.NotificationChannel{constants.EmailChannel},
		}

		cardReadRepositoryMock.On(
			"GetCardsExpiringSoon", mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)

		useCase := NewCardUseCase(
			cardAdapter,
			repositoryMock,
			cardReadRepositoryMock,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		_, err := useCase.TriggerCardExpiringSoonNotifications(ctx, requestMock)

		// Assert
		assert.Error(t, err)
	})

	t.Run("should not trigger card expiring soon notifications when there are no cards expiring soon", func(t *testing.T) {
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		cardAdapter := adapters.NewCardAdapter(t)
		cardReadRepositoryMock := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.NotificationCardExpiringSoonRequestDTO{
			NotificationChannels: []constants.NotificationChannel{constants.EmailChannel},
		}

		cardReadRepositoryMock.On(
			"GetCardsExpiringSoon", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		useCase := NewCardUseCase(
			cardAdapter,
			repositoryMock,
			cardReadRepositoryMock,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		response, err := useCase.TriggerCardExpiringSoonNotifications(ctx, requestMock)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, "no cards found to trigger expiring soon notifications", response.Message)
	})

	t.Run("should trigger card expiring soon notifications successfully", func(t *testing.T) {
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		cardAdapter := adapters.NewCardAdapter(t)
		cardReadRepositoryMock := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.NotificationCardExpiringSoonRequestDTO{
			NotificationChannels: []constants.NotificationChannel{constants.EmailChannel},
		}

		cardReadRepositoryMock.On(
			"GetCardsExpiringSoon", mock.Anything, mock.Anything, mock.Anything).Return(notificationCardExpiringSoonProjections, nil)

		notificationServiceMock.On(
			"NotifyCardExpiringSoon",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil)

		logMock.On(
			"Info", mock.Anything, mock.Anything).Return(nil)

		useCase := NewCardUseCase(
			cardAdapter,
			repositoryMock,
			cardReadRepositoryMock,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)
		response, err := useCase.TriggerCardExpiringSoonNotifications(ctx, requestMock)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, "card expiring soon notification has been triggered", response.Message)
	})

	t.Run("should fail when notify card expiring soon fails", func(t *testing.T) {
		repositoryMock := repository.NewCardWriteRepositoryIF(t)
		cardAdapter := adapters.NewCardAdapter(t)
		cardReadRepositoryMock := repository.NewCardReadRepositoryIF(t)
		userReadRepositoryMock := repository.NewUserReadRepositoryIF(t)
		notificationServiceMock := services.NewNotificationService(t)
		logMock := adapters.NewLogger(t)

		requestMock := request.NotificationCardExpiringSoonRequestDTO{
			NotificationChannels: []constants.NotificationChannel{constants.EmailChannel},
		}

		cardReadRepositoryMock.On(
			"GetCardsExpiringSoon", mock.Anything, mock.Anything, mock.Anything).Return(notificationCardExpiringSoonProjections, nil)

		notificationServiceMock.On(
			"NotifyCardExpiringSoon",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(assert.AnError)

		logMock.On(
			"Warn", mock.Anything, mock.Anything).Return(nil)

		useCase := NewCardUseCase(
			cardAdapter,
			repositoryMock,
			cardReadRepositoryMock,
			userReadRepositoryMock,
			notificationServiceMock,
			logMock,
		)

		response, err := useCase.TriggerCardExpiringSoonNotifications(ctx, requestMock)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, "card expiring soon notification has been triggered", response.Message)
	})

}

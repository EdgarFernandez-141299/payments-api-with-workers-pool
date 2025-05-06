package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/user/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/adapters"
)

func TestNotifyCardAddition(t *testing.T) {
	mailAdapterMock := adapters.NewMailAdapterIF(t)
	memberAdapterMock := adapters.NewMemberAdapterIF(t)

	userID := "123"
	email := "test@test.com"
	actionDate := "2021-01-01"
	lastFour := "1234"
	userLanguage := "en"

	t.Run("should notify card addition", func(t *testing.T) {
		mailAdapterMock := adapters.NewMailAdapterIF(t)
		memberAdapterMock := adapters.NewMemberAdapterIF(t)

		notificationService := NewNotificationService(mailAdapterMock, memberAdapterMock)

		mailAdapterMock.On("Send", mock.Anything, mock.Anything).Return(nil)

		err := notificationService.NotifyCardAddition(
			context.Background(),
			[]constants.NotificationChannel{constants.EmailChannel},
			&userID,
			&email,
			&actionDate,
			&lastFour,
			userLanguage,
		)

		assert.NoError(t, err)
		mailAdapterMock.ExpectedCalls = nil
		memberAdapterMock.ExpectedCalls = nil
	})

	t.Run("should return error", func(t *testing.T) {

		notificationService := NewNotificationService(mailAdapterMock, memberAdapterMock)

		mailAdapterMock.On("Send", mock.Anything, mock.Anything).Return(errors.New("error"))

		err := notificationService.NotifyCardAddition(
			context.Background(),
			[]constants.NotificationChannel{constants.EmailChannel},
			&userID,
			&email,
			&actionDate,
			&lastFour,
			userLanguage,
		)

		assert.Error(t, err)
		mailAdapterMock.ExpectedCalls = nil
		memberAdapterMock.ExpectedCalls = nil
	})
}

func TestNotifyCardDeletion(t *testing.T) {
	mailAdapterMock := adapters.NewMailAdapterIF(t)
	memberAdapterMock := adapters.NewMemberAdapterIF(t)

	userID := "123"
	email := "test@test.com"
	actionDate := "2021-01-01"
	lastFour := "1234"
	userLanguage := "fr"

	t.Run("should notify card deletion", func(t *testing.T) {
		notificationService := NewNotificationService(mailAdapterMock, memberAdapterMock)

		mailAdapterMock.On("Send", mock.Anything, mock.Anything).Return(nil)

		err := notificationService.NotifyCardDeletion(
			context.Background(),
			[]constants.NotificationChannel{constants.EmailChannel},
			&userID,
			&email,
			&actionDate,
			&lastFour,
			userLanguage,
		)

		assert.NoError(t, err)
		mailAdapterMock.ExpectedCalls = nil
		memberAdapterMock.ExpectedCalls = nil
	})

	t.Run("should return error", func(t *testing.T) {
		notificationService := NewNotificationService(mailAdapterMock, memberAdapterMock)

		mailAdapterMock.On("Send", mock.Anything, mock.Anything).Return(errors.New("error"))

		err := notificationService.NotifyCardDeletion(
			context.Background(),
			[]constants.NotificationChannel{constants.EmailChannel},
			&userID,
			&email,
			&actionDate,
			&lastFour,
			userLanguage,
		)

		assert.Error(t, err)
		mailAdapterMock.ExpectedCalls = nil
		memberAdapterMock.ExpectedCalls = nil
	})
}

func TestNotifyCardExpiringSoon(t *testing.T) {
	mailAdapterMock := adapters.NewMailAdapterIF(t)
	memberAdapterMock := adapters.NewMemberAdapterIF(t)

	t.Run("should notify card expiring soon", func(t *testing.T) {
		notificationService := NewNotificationService(mailAdapterMock, memberAdapterMock)

		memberAdapterMock.On("GetUserProfileInfo", mock.Anything, mock.Anything, mock.Anything).Return(
			&response.UserProfileInfoDTO{
				PreferenceLanguage: "",
			}, nil)

		mailAdapterMock.On("Send", mock.Anything, mock.Anything).Return(nil)

		err := notificationService.NotifyCardExpiringSoon(
			context.Background(),
			[]constants.NotificationChannel{constants.EmailChannel},
			projections.NotificationCardExpiringSoonProjection{},
		)

		assert.NoError(t, err)
		mailAdapterMock.ExpectedCalls = nil
		memberAdapterMock.ExpectedCalls = nil
	})

	t.Run("should return error when getting user profile info fails", func(t *testing.T) {
		notificationService := NewNotificationService(mailAdapterMock, memberAdapterMock)

		memberAdapterMock.On("GetUserProfileInfo", mock.Anything, mock.Anything, mock.Anything).Return(
			nil, errors.New("error"))

		err := notificationService.NotifyCardExpiringSoon(
			context.Background(),
			[]constants.NotificationChannel{constants.EmailChannel},
			projections.NotificationCardExpiringSoonProjection{},
		)

		assert.Error(t, err)
		mailAdapterMock.ExpectedCalls = nil
		memberAdapterMock.ExpectedCalls = nil
	})

	t.Run("should return error when sending email fails", func(t *testing.T) {
		notificationService := NewNotificationService(mailAdapterMock, memberAdapterMock)

		memberAdapterMock.On("GetUserProfileInfo", mock.Anything, mock.Anything, mock.Anything).Return(
			&response.UserProfileInfoDTO{
				PreferenceLanguage: "",
			}, nil)

		mailAdapterMock.On("Send", mock.Anything, mock.Anything).Return(errors.New("error"))

		err := notificationService.NotifyCardExpiringSoon(
			context.Background(),
			[]constants.NotificationChannel{constants.EmailChannel},
			projections.NotificationCardExpiringSoonProjection{},
		)

		assert.Error(t, err)
		mailAdapterMock.ExpectedCalls = nil
		memberAdapterMock.ExpectedCalls = nil
	})
}

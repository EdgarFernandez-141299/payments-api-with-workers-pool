package services

import (
	"context"
	"strings"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/config"
	commonAdapters "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/common/adapters"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/dto"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/card/projections"
)

type NotificationService interface {
	NotifyCardAddition(
		ctx context.Context,
		notificationChannels []constants.NotificationChannel,
		userID *string,
		email *string,
		actionDate *string,
		lastFour *string,
		userLanguage string,
	) error

	NotifyCardDeletion(
		ctx context.Context,
		notificationChannels []constants.NotificationChannel,
		userID *string,
		email *string,
		actionDate *string,
		lastFour *string,
		userLanguage string,
	) error

	NotifyCardExpiringSoon(
		ctx context.Context,
		notificationChannels []constants.NotificationChannel,
		projection projections.NotificationCardExpiringSoonProjection,
	) error
}

type notificationServiceImpl struct {
	mailAdapter   commonAdapters.MailAdapterIF
	memberAdapter commonAdapters.MemberAdapterIF
}

func NewNotificationService(
	mailAdapter commonAdapters.MailAdapterIF,
	memberAdapter commonAdapters.MemberAdapterIF,
) NotificationService {
	return &notificationServiceImpl{
		mailAdapter:   mailAdapter,
		memberAdapter: memberAdapter,
	}
}

func (nsi *notificationServiceImpl) NotifyCardAddition(
	ctx context.Context,
	channels []constants.NotificationChannel,
	userID, email, actionDate, lastFour *string,
	lang string,
) error {
	return nsi.sendCardNotification(
		ctx,
		&dto.CardNotificationDTO{
			NotificationType: constants.NotificationTypeCardAddition,
			Template:         "card-addition",
			UserID:           *userID,
			Email:            *email,
			ActionDate:       *actionDate,
			LastFour:         *lastFour,
			Language:         lang,
			Channels:         channels,
		},
	)
}

func (nsi *notificationServiceImpl) NotifyCardDeletion(
	ctx context.Context,
	channels []constants.NotificationChannel,
	userID, email, actionDate, lastFour *string,
	lang string,
) error {
	return nsi.sendCardNotification(
		ctx,
		&dto.CardNotificationDTO{
			NotificationType: constants.NotificationTypeCardDeletion,
			Template:         "card-remotion",
			UserID:           *userID,
			Email:            *email,
			ActionDate:       *actionDate,
			LastFour:         *lastFour,
			Language:         lang,
			Channels:         channels,
		},
	)
}

func (nsi *notificationServiceImpl) NotifyCardExpiringSoon(
	ctx context.Context,
	channels []constants.NotificationChannel,
	notificationCardExpiringSoonProjection projections.NotificationCardExpiringSoonProjection,
) error {
	formattedDate := notificationCardExpiringSoonProjection.ExpirationDate.Format("01/06")

	userInfo, err := nsi.memberAdapter.GetUserProfileInfo(
		ctx,
		notificationCardExpiringSoonProjection.UserID,
		notificationCardExpiringSoonProjection.EnterpriseID,
	)
	if err != nil {
		return err
	}

	title, lang := nsi.getTitleAndLanguage(
		userInfo.PreferenceLanguage,
		constants.NotificationTypeCardExpiringSoon,
	)

	metadata := commonAdapters.Metadata{}
	metadata.Add("language", lang)
	metadata.Add("template", "card-expiration")
	metadata.Add("to", notificationCardExpiringSoonProjection.Email)
	metadata.Add("card_last_digits", notificationCardExpiringSoonProjection.LastFour)
	metadata.Add("card_expiration_date", formattedDate)
	metadata.Add("update_card_button_cta", config.Config().NotificationLinkEmail.Wallet)

	return nsi.mailAdapter.Send(ctx, commonAdapters.MailRequest{
		Recipient: notificationCardExpiringSoonProjection.UserID,
		Content:   "Your card is about to expire",
		Title:     title,
		Channels:  channels,
		Metadata:  metadata,
	})
}

func (nsi *notificationServiceImpl) sendCardNotification(
	ctx context.Context,
	cardNotificationDTO *dto.CardNotificationDTO,
) error {
	title, language := nsi.getTitleAndLanguage(cardNotificationDTO.Language, cardNotificationDTO.NotificationType)

	metadata := commonAdapters.Metadata{}
	metadata.Add("language", language)
	metadata.Add("template", cardNotificationDTO.Template)
	metadata.Add("to", cardNotificationDTO.Email)
	metadata.Add("card_last_digits", cardNotificationDTO.LastFour)
	metadata.Add("action_date", cardNotificationDTO.ActionDate)
	metadata.Add("manage_cards_button_cta", config.Config().NotificationLinkEmail.Wallet)
	metadata.Add("contact_support_button_cta", config.Config().NotificationLinkEmail.ContactSupport)

	return nsi.mailAdapter.Send(ctx, commonAdapters.MailRequest{
		Recipient: cardNotificationDTO.UserID,
		Content:   "Card notification",
		Title:     title,
		Channels:  cardNotificationDTO.Channels,
		Metadata:  metadata,
	})
}

func (nsi *notificationServiceImpl) getTitleAndLanguage(
	userLanguage string,
	notificationType constants.NotificationType,
) (string, string) {
	language := strings.ToLower(userLanguage)
	if language == "" {
		language = "en"
	}

	titleMap, ok := constants.TitlesMap[notificationType]
	if !ok {
		titleMap = constants.TitlesMap[constants.NotificationTypeCardExpiringSoon]
	}

	title, ok := titleMap[language]
	if !ok {
		title = titleMap["en"]
	}

	return title, language
}

package dto

import "gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/notification/constants"

type CardNotificationDTO struct {
	NotificationType constants.NotificationType
	Template         string
	UserID           string
	Email            string
	ActionDate       string
	LastFour         string
	Language         string
	Channels         []constants.NotificationChannel
}

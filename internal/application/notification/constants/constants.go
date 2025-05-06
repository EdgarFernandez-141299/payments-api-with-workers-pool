package constants

type NotificationChannel string
type NotificationType string

const (
	EmailChannel NotificationChannel = "EMAIL"
)

const (
	CardAdditionTitleEnglish   = "Your Card Has Been Added"
	CardDeletionTitleEnglish   = "Your Card Has Been Deleted"
	CardExpirationTitleEnglish = "Your Saved Card is Expiring Soon - Update Now"
	CardAdditionTitleSpanish   = "Tu Tarjeta Ha Sido Agregada"
	CardDeletionTitleSpanish   = "Tu Tarjeta Ha Sido Eliminada"
	CardExpirationTitleSpanish = "Su Tarjeta Guardada Expirará Pronto - Actualícela Ahora"
)

const (
	NotificationTypeCardAddition     NotificationType = "cardAddition"
	NotificationTypeCardDeletion     NotificationType = "cardDeletion"
	NotificationTypeCardExpiringSoon NotificationType = "cardExpiringSoon"
)

var AllowedNotificationChannels = map[NotificationChannel]bool{
	EmailChannel: true,
}

var TitlesMap = map[NotificationType]map[string]string{
	NotificationTypeCardAddition: {
		"en": CardAdditionTitleEnglish,
		"es": CardAdditionTitleSpanish,
	},
	NotificationTypeCardDeletion: {
		"en": CardDeletionTitleEnglish,
		"es": CardDeletionTitleSpanish,
	},
	NotificationTypeCardExpiringSoon: {
		"en": CardExpirationTitleEnglish,
		"es": CardExpirationTitleSpanish,
	},
}

package resources

import (
	"context"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type WebhookNotificationResourceIF interface {
	SendNotification(ctx context.Context, webhookUrl value_objects.WebhookUrl, payload interface{}) error
}

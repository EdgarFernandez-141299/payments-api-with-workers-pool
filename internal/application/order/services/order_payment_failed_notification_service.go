package services

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/integration_events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type OrderFailedNotificationServiceIF interface {
	Notify(ctx context.Context, webhookUrl value_objects.WebhookUrl, event integration_events.OrderFailedPaidIntegrationEvent) error
}

type OrderFailedNotificationService struct {
	resource resources.WebhookNotificationResourceIF
}

func NewOrderFailedNotificationService(
	resource resources.WebhookNotificationResourceIF,
) OrderFailedNotificationServiceIF {
	return &OrderFailedNotificationService{
		resource: resource,
	}
}

func (s *OrderFailedNotificationService) Notify(ctx context.Context, webhookUrl value_objects.WebhookUrl, event integration_events.OrderFailedPaidIntegrationEvent) error {
	return s.resource.SendNotification(ctx, webhookUrl, event)
}

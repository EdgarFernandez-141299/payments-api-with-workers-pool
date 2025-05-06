package services

import (
	"context"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/order/adapters/resources"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/integration_events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
)

type OrderProcessedNotificationServiceIF interface {
	Notify(ctx context.Context, webhookUrl value_objects.WebhookUrl, event integration_events.OrderPaymentProcessedIntegrationEvent) error
}

type OrderProcessedNotificationService struct {
	resource resources.WebhookNotificationResourceIF
}

func NewOrderProcessedNotificationService(
	resource resources.WebhookNotificationResourceIF,
) OrderProcessedNotificationServiceIF {
	return &OrderProcessedNotificationService{
		resource: resource,
	}
}

func (s *OrderProcessedNotificationService) Notify(ctx context.Context, webhookUrl value_objects.WebhookUrl, event integration_events.OrderPaymentProcessedIntegrationEvent) error {
	return s.resource.SendNotification(ctx, webhookUrl, event)
}

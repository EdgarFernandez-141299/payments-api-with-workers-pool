package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/integration_events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	mockResources "gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

func TestOrderProcessedNotificationService_Notify(t *testing.T) {
	tests := []struct {
		name          string
		webhookUrl    value_objects.WebhookUrl
		event         integration_events.OrderPaymentProcessedIntegrationEvent
		expectedError error
	}{
		{
			name:       "should send notification successfully",
			webhookUrl: value_objects.NewWebhookUrl("https://example.com/webhook"),
			event: integration_events.NewOrderPaymentProcessedIntegrationEvent(
				integration_events.IntegrationEventsParams{
					ReferenceOrderID:   "123",
					ReferencePaymentID: "payment123",
					AssociatedPayment:  "assoc123",
					TotalOrderAmount:   100.0,
					Currency:           "USD",
					UserID:             "user123",
					UserType:           "regular",
					EnterpriseID:       "enterprise123",
					TotalOrderPaid:     100.0,
					TotalPaymentAmount: 100.0,
					CardData: integration_events.CardData{
						CardNumber: "****1234",
						CardType:   "visa",
					},
					Metadata: map[string]interface{}{"key": "value"},
				},
				"AUTH123",
				"completed",
			),
			expectedError: nil,
		},
		{
			name:       "should return error when notification fails",
			webhookUrl: value_objects.NewWebhookUrl("https://example.com/webhook"),
			event: integration_events.NewOrderPaymentProcessedIntegrationEvent(
				integration_events.IntegrationEventsParams{
					ReferenceOrderID:   "123",
					ReferencePaymentID: "payment123",
					AssociatedPayment:  "assoc123",
					TotalOrderAmount:   100.0,
					Currency:           "USD",
					UserID:             "user123",
					UserType:           "regular",
					EnterpriseID:       "enterprise123",
					TotalOrderPaid:     100.0,
					TotalPaymentAmount: 100.0,
					CardData: integration_events.CardData{
						CardNumber: "****1234",
						CardType:   "visa",
					},
					Metadata: map[string]interface{}{"key": "value"},
				},
				"AUTH123",
				"completed",
			),
			expectedError: errors.New("failed to send notification"),
		},
		{
			name:       "should handle empty metadata",
			webhookUrl: value_objects.NewWebhookUrl("https://example.com/webhook"),
			event: integration_events.NewOrderPaymentProcessedIntegrationEvent(
				integration_events.IntegrationEventsParams{
					ReferenceOrderID:   "123",
					ReferencePaymentID: "payment123",
					AssociatedPayment:  "assoc123",
					TotalOrderAmount:   100.0,
					Currency:           "USD",
					UserID:             "user123",
					UserType:           "regular",
					EnterpriseID:       "enterprise123",
					TotalOrderPaid:     100.0,
					TotalPaymentAmount: 100.0,
					CardData: integration_events.CardData{
						CardNumber: "****1234",
						CardType:   "visa",
					},
					Metadata: nil,
				},
				"AUTH123",
				"completed",
			),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockResource := new(mockResources.WebhookNotificationResourceIF)
			service := NewOrderProcessedNotificationService(mockResource)

			mockResource.On("SendNotification", mock.Anything, tt.webhookUrl, tt.event).Return(tt.expectedError)

			// Act
			err := service.Notify(context.Background(), tt.webhookUrl, tt.event)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
			mockResource.AssertExpectations(t)
		})
	}
}

package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/enums"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/order/integration_events"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/domain/value_objects"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/event_store"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository/fixture"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/services"
)

func TestNewOrderNotificationOrchestrator_NotifyChange(t *testing.T) {
	billsApiWebhookUrl := value_objects.NewWebhookUrl("http://bill_api.co/test")
	orderWebhookUrl := value_objects.NewWebhookUrl("http://order.api/test")

	tests := []struct {
		name          string
		orderID       string
		paymentID     string
		paymentStatus enums.PaymentStatus
		setupMocks    func(*event_store.OrderEventRepository, *services.OrderFailedNotificationServiceIF, *services.OrderProcessedNotificationServiceIF)
		expectedError error
	}{
		{
			name:          "successful payment notification",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentStatus: enums.PaymentProcessed,
			setupMocks: func(
				orderEventMock *event_store.OrderEventRepository,
				orderFailedNotificationMock *services.OrderFailedNotificationServiceIF,
				orderProcessedNotificationServiceMock *services.OrderProcessedNotificationServiceIF,
			) {
				fixture.NewFromOrderMock(orderEventMock).
					WithPaymentOrderID("payment123").
					WithReferenceOrderID("order123").
					WithPaymentStatus(enums.PaymentProcessed).
					WithWebhookURL(orderWebhookUrl).
					Build()

				orderProcessedNotificationServiceMock.On(
					"Notify",
					mock.Anything,
					mock.Anything,
					mock.IsType(integration_events.OrderPaymentProcessedIntegrationEvent{}),
				).
					Times(2).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:          "when fail successful payment notification",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentStatus: enums.PaymentProcessed,
			setupMocks: func(
				orderEventMock *event_store.OrderEventRepository,
				orderFailedNotificationMock *services.OrderFailedNotificationServiceIF,
				orderProcessedNotificationServiceMock *services.OrderProcessedNotificationServiceIF,
			) {
				fixture.NewFromOrderMock(orderEventMock).
					WithPaymentOrderID("payment123").
					WithReferenceOrderID("order123").
					WithPaymentStatus(enums.PaymentProcessed).
					WithWebhookURL(orderWebhookUrl).
					Build()

				orderProcessedNotificationServiceMock.On(
					"Notify",
					mock.Anything,
					mock.Anything,
					mock.IsType(integration_events.OrderPaymentProcessedIntegrationEvent{}),
				).
					Times(2).
					Return(assert.AnError)
			},
			expectedError: assert.AnError,
		},
		{
			name:          "when success payment failure notification",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentStatus: enums.PaymentFailed,
			setupMocks: func(
				orderEventMock *event_store.OrderEventRepository,
				orderFailedNotificationMock *services.OrderFailedNotificationServiceIF,
				orderProcessedNotificationServiceMock *services.OrderProcessedNotificationServiceIF,
			) {
				fixture.NewFromOrderMock(orderEventMock).
					WithPaymentOrderID("payment123").
					WithReferenceOrderID("order123").
					WithPaymentStatus(enums.PaymentFailed).
					WithWebhookURL(orderWebhookUrl).
					Build()

				orderFailedNotificationMock.On(
					"Notify",
					mock.Anything,
					mock.Anything,
					mock.IsType(integration_events.OrderFailedPaidIntegrationEvent{}),
				).
					Times(2).
					Return(nil)
			},
			expectedError: fmt.Errorf("billsError: %v", nil),
		},
		{
			name:          "when fail failed payment notification",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentStatus: enums.PaymentFailed,
			setupMocks: func(
				orderEventMock *event_store.OrderEventRepository,
				orderFailedNotificationMock *services.OrderFailedNotificationServiceIF,
				orderProcessedNotificationServiceMock *services.OrderProcessedNotificationServiceIF,
			) {
				fixture.NewFromOrderMock(orderEventMock).
					WithPaymentOrderID("payment123").
					WithReferenceOrderID("order123").
					WithPaymentStatus(enums.PaymentFailed).
					WithWebhookURL(orderWebhookUrl).
					Build()

				orderFailedNotificationMock.On(
					"Notify",
					mock.Anything,
					mock.Anything,
					mock.IsType(integration_events.OrderFailedPaidIntegrationEvent{}),
				).
					Times(2).
					Return(assert.AnError)
			},
			expectedError: fmt.Errorf("billsError: %v", assert.AnError),
		},
		{
			name:          "when order not found",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentStatus: enums.PaymentProcessed,
			setupMocks: func(
				orderEventMock *event_store.OrderEventRepository,
				orderFailedNotificationMock *services.OrderFailedNotificationServiceIF,
				orderProcessedNotificationServiceMock *services.OrderProcessedNotificationServiceIF,
			) {
				orderEventMock.On("Get", mock.Anything, "order123", mock.Anything).
					Return(errors.NewOrderNotFoundError("order123"))
			},
			expectedError: errors.NewOrderNotFoundError("order123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(event_store.OrderEventRepository)
			mockFailedNotif := new(services.OrderFailedNotificationServiceIF)
			mockSuccessNotif := new(services.OrderProcessedNotificationServiceIF)

			tt.setupMocks(mockRepo, mockFailedNotif, mockSuccessNotif)

			strategy := NewOrderNotificationOrchestrator(mockRepo, mockFailedNotif, mockSuccessNotif, BillsApiWebhookUrl(billsApiWebhookUrl))

			err := strategy.NotifyChange(context.Background(), tt.orderID, tt.paymentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.name == "when order not found" {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				} else {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockFailedNotif.AssertExpectations(t)
			mockSuccessNotif.AssertExpectations(t)
		})
	}
}

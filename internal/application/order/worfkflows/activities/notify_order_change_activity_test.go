package activities

import (
	"context"
	"errors"
	errors2 "gitlab.com/clubhub.ai1/go-libraries/saga/errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotifyOrderChangeActivity_NotifyOrderChange(t *testing.T) {
	tests := []struct {
		name          string
		params        NotifyOrderChangeParams
		mockSetup     func(*services.OrderNotificationStrategyIF)
		expectedError error
	}{
		{
			name: "notification success",
			params: NotifyOrderChangeParams{
				OrderID:   "order123",
				PaymentID: "payment123",
			},
			mockSetup: func(m *services.OrderNotificationStrategyIF) {
				m.On("NotifyChange", mock.Anything, "order123", "payment123").Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "notification error",
			params: NotifyOrderChangeParams{
				OrderID:   "order123",
				PaymentID: "payment123",
			},
			mockSetup: func(m *services.OrderNotificationStrategyIF) {
				m.On("NotifyChange", mock.Anything, "order123", "payment123").Return(errors.New("notification error"))
			},
			expectedError: errors2.WrapActivityError(errors.New("notification error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStrategy := new(services.OrderNotificationStrategyIF)
			tt.mockSetup(mockStrategy)

			activity := NewNotifyOrderChangeActivity(mockStrategy)
			err := activity.NotifyOrderChange(context.Background(), tt.params)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockStrategy.AssertExpectations(t)
		})
	}
}

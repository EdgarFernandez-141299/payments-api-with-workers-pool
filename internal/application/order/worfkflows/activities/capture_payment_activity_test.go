package activities

import (
	"context"
	"errors"
	"testing"

	sagaErrors "gitlab.com/clubhub.ai1/go-libraries/saga/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCapturePaymentUseCase struct {
	mock.Mock
}

func (m *mockCapturePaymentUseCase) CapturePayment(ctx context.Context, orderID, paymentID string) error {
	args := m.Called(ctx, orderID, paymentID)
	return args.Error(0)
}

func TestCapturePaymentActivity_CapturePayment(t *testing.T) {
	tests := []struct {
		name          string
		orderID       string
		paymentID     string
		expectedError error
		setupMock     func(*mockCapturePaymentUseCase)
	}{
		{
			name:      "Captura exitosa",
			orderID:   "order-123",
			paymentID: "payment-456",
			setupMock: func(m *mockCapturePaymentUseCase) {
				m.On("CapturePayment", mock.Anything, "order-123", "payment-456").Return(nil)
			},
		},
		{
			name:          "Error en la captura",
			orderID:       "order-123",
			paymentID:     "payment-456",
			expectedError: sagaErrors.WrapActivityError(errors.New("error de captura")),
			setupMock: func(m *mockCapturePaymentUseCase) {
				m.On("CapturePayment", mock.Anything, "order-123", "payment-456").Return(errors.New("error de captura"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := new(mockCapturePaymentUseCase)
			tt.setupMock(mockUseCase)

			activity := NewCapturePaymentActivity(mockUseCase)
			err := activity.CapturePayment(context.Background(), tt.orderID, tt.paymentID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockUseCase.AssertExpectations(t)
		})
	}
}

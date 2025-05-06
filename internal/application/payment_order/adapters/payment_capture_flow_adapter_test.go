package adapter

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/repository"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

func TestPaymentCaptureFlowAdapter_CapturePayment(t *testing.T) {
	tests := []struct {
		name          string
		orderID       string
		paymentID     string
		paymentTotal  decimal.Decimal
		token         string
		captureResult bool
		expectedErr   error
	}{
		{
			name:          "Captura exitosa",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentTotal:  decimal.NewFromFloat(100.50),
			token:         "token123",
			captureResult: true,
			expectedErr:   nil,
		},
		{
			name:          "Error al obtener token",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentTotal:  decimal.NewFromFloat(100.50),
			token:         "",
			captureResult: false,
			expectedErr:   errors.New("error al obtener token"),
		},
		{
			name:          "Error en captura",
			orderID:       "order123",
			paymentID:     "payment123",
			paymentTotal:  decimal.NewFromFloat(100.50),
			token:         "token123",
			captureResult: false,
			expectedErr:   errors.New("capture failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockResource := resources.NewDeunaCaptureFlowResourceIF(t)
			mockRepo := repository.NewDeunaOrderRepository(t)

			adapter := NewPaymentCaptureFlowAdapter(mockResource, mockRepo)

			mockRepo.On("GetTokenByOrderAndPaymentID", mock.Anything, tt.orderID, tt.paymentID).
				Return(tt.token, tt.expectedErr)

			if tt.expectedErr == nil {
				mockResource.On("Capture", mock.Anything, tt.token, int64(10050)).
					Return(tt.captureResult, nil)
			}

			err := adapter.CapturePayment(ctx, tt.orderID, tt.paymentID, tt.paymentTotal)

			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockResource.AssertExpectations(t)
		})
	}
}

func TestPaymentCaptureFlowAdapter_ReleasePayment(t *testing.T) {
	tests := []struct {
		name          string
		orderID       string
		paymentID     string
		reason        string
		token         string
		releaseResult bool
		expectedErr   error
	}{
		{
			name:          "Liberación exitosa",
			orderID:       "order123",
			paymentID:     "payment123",
			reason:        "cancelled",
			token:         "token123",
			releaseResult: true,
			expectedErr:   nil,
		},
		{
			name:          "Error al obtener token",
			orderID:       "order123",
			paymentID:     "payment123",
			reason:        "cancelled",
			token:         "",
			releaseResult: false,
			expectedErr:   errors.New("error al obtener token"),
		},
		{
			name:          "Error en liberación",
			orderID:       "order123",
			paymentID:     "payment123",
			reason:        "cancelled",
			token:         "token123",
			releaseResult: false,
			expectedErr:   errors.New("release failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockResource := resources.NewDeunaCaptureFlowResourceIF(t)
			mockRepo := repository.NewDeunaOrderRepository(t)

			adapter := NewPaymentCaptureFlowAdapter(mockResource, mockRepo)

			mockRepo.On("GetTokenByOrderAndPaymentID", mock.Anything, tt.orderID, tt.paymentID).
				Return(tt.token, tt.expectedErr)

			if tt.expectedErr == nil {
				mockResource.On("Release", mock.Anything, tt.token, tt.reason).
					Return(tt.releaseResult, nil)
			}

			err := adapter.ReleasePayment(ctx, tt.orderID, tt.paymentID, tt.reason)

			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockResource.AssertExpectations(t)
		})
	}
}

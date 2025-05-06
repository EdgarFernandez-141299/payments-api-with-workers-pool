package http

import (
	"context"
	"errors"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/internal/application/payment_order/adapters/resources/dto/response"
	"gitlab.com/clubhub.ai1/organization/backend/payments-api/mocks/resources"
)

func TestDeUnaRefundHTTPClient_MakeTotalRefund(t *testing.T) {
	tests := []struct {
		name          string
		body          utils.DeunaTotalRefundRequest
		orderToken    string
		mockResponse  response.DeunaRefundPaymentResponse
		mockError     error
		expectedError error
	}{
		{
			name: "successful refund",
			body: utils.DeunaTotalRefundRequest{
				Reason: "customer request",
			},
			orderToken:   "order-token-123",
			mockResponse: response.DeunaRefundPaymentResponse{
				// Using fields that actually exist in the response struct
				// These would need to be updated based on the actual struct definition
			},
			mockError:     nil,
			expectedError: nil,
		},
		{
			name: "payment client returns error",
			body: utils.DeunaTotalRefundRequest{
				Reason: "customer request",
			},
			orderToken:    "order-token-123",
			mockResponse:  response.DeunaRefundPaymentResponse{},
			mockError:     errors.New("payment service unavailable"),
			expectedError: errors.New("payment service unavailable"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			mockPaymentClient := resources.NewDeunaPaymentResourceIF(t)

			// Configure mock
			mockPaymentClient.On("MakeTotalRefund", ctx, tt.body, tt.orderToken).
				Return(tt.mockResponse, tt.mockError)

			// Create client with mock
			client := NewDeUnaRefundHTTPClient(mockPaymentClient)

			// Act
			response, err := client.MakeTotalRefund(ctx, tt.body, tt.orderToken)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse, response)
			}

			// Verify mock expectations
			mockPaymentClient.AssertExpectations(t)
		})
	}
}
